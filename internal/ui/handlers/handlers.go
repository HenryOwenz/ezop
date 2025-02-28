package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// ModelWrapper wraps a core.Model to implement the tea.Model interface
type ModelWrapper struct {
	Model *core.Model
}

// Update implements the tea.Model interface
func (m ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// This is just a placeholder - the actual update logic will be in the UI package
	return m, nil
}

// View implements the tea.Model interface
func (m ModelWrapper) View() string {
	// This is just a placeholder - the actual view logic will be in the UI package
	return ""
}

// Init implements the tea.Model interface
func (m ModelWrapper) Init() tea.Cmd {
	// This is just a placeholder - the actual init logic will be in the UI package
	return nil
}

// WrapModel wraps a core.Model in a ModelWrapper
func WrapModel(m *core.Model) ModelWrapper {
	return ModelWrapper{Model: m}
}

// HandleEnter processes the enter key press based on the current view
func HandleEnter(m *core.Model) (tea.Model, tea.Cmd) {
	// Special handling for manual input in AWS config view
	if m.CurrentView == constants.ViewAWSConfig && m.ManualInput {
		newModel := m.Clone()

		// Get the entered value
		value := strings.TrimSpace(m.TextInput.Value())
		if value == "" {
			// If empty, just exit manual input mode
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}

		// Set the appropriate value based on context
		if m.AwsProfile == "" {
			// Setting profile
			newModel.AwsProfile = value
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
		} else {
			// Setting region and moving to next view
			newModel.AwsRegion = value
			newModel.ManualInput = false
			newModel.ResetTextInput()

			// Create the provider with the selected profile and region
			_, err := providers.CreateProvider(newModel.Registry, "AWS", newModel.AwsProfile, newModel.AwsRegion)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return core.ErrMsg{Err: err}
				}
			}

			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(newModel)
		}

		return WrapModel(newModel), nil
	}

	// Regular view handling
	switch m.CurrentView {
	case constants.ViewProviders:
		return HandleProviderSelection(m)
	case constants.ViewAWSConfig:
		return HandleAWSConfigSelection(m)
	case constants.ViewSelectService:
		return HandleServiceSelection(m)
	case constants.ViewSelectCategory:
		return HandleCategorySelection(m)
	case constants.ViewSelectOperation:
		return HandleOperationSelection(m)
	case constants.ViewApprovals:
		return HandleApprovalSelection(m)
	case constants.ViewConfirmation:
		return HandleConfirmationSelection(m)
	case constants.ViewSummary:
		if !m.ManualInput {
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				if selected := m.Table.SelectedRow(); len(selected) > 0 {
					newModel := m.Clone()
					switch selected[0] {
					case "Latest Commit":
						newModel.CurrentView = constants.ViewExecutingAction
						newModel.Summary = "" // Empty string means use latest commit
						view.UpdateTableForView(newModel)
						return WrapModel(newModel), nil
					case "Manual Input":
						newModel.ManualInput = true
						newModel.TextInput.Focus()
						newModel.TextInput.Placeholder = constants.MsgEnterCommitID
						return WrapModel(newModel), nil
					}
				}
			}
		}
		return HandleSummaryConfirmation(m)
	case constants.ViewExecutingAction:
		return HandleExecutionSelection(m)
	case constants.ViewPipelineStatus:
		if selected := m.Table.SelectedRow(); len(selected) > 0 {
			newModel := m.Clone()
			for _, pipeline := range m.Pipelines {
				if pipeline.Name == selected[0] {
					if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
						newModel.CurrentView = constants.ViewExecutingAction
						newModel.SelectedPipeline = &pipeline
						view.UpdateTableForView(newModel)
						return WrapModel(newModel), nil
					}
					newModel.CurrentView = constants.ViewPipelineStages
					newModel.SelectedPipeline = &pipeline
					view.UpdateTableForView(newModel)
					return WrapModel(newModel), nil
				}
			}
		}
	case constants.ViewPipelineStages:
		// Just view only, no action
	}
	return WrapModel(m), nil
}

// HandleProviderSelection handles the selection of a cloud provider
func HandleProviderSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		providerName := selected[0]

		// Initialize providers if not already done
		if len(m.Registry.GetProviderNames()) == 0 {
			providers.InitializeProviders(m.Registry)
		}

		// Check if the provider exists in the registry
		_, exists := m.Registry.GetProvider(providerName)
		if exists {
			newModel := m.Clone()
			newModel.CurrentView = constants.ViewAWSConfig
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleAWSConfigSelection handles the selection of AWS profile or region
func HandleAWSConfigSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()

		// Handle "Manual Entry" option
		if selected[0] == "Manual Entry" {
			newModel.ManualInput = true
			newModel.TextInput.Focus()

			// Set appropriate placeholder based on context
			if m.AwsProfile == "" {
				newModel.TextInput.Placeholder = constants.MsgEnterProfile
			} else {
				newModel.TextInput.Placeholder = constants.MsgEnterRegion
			}

			return WrapModel(newModel), nil
		}

		// Handle regular selection
		if m.AwsProfile == "" {
			newModel.AwsProfile = selected[0]
			view.UpdateTableForView(newModel)
		} else {
			newModel.AwsRegion = selected[0]

			// Create the provider with the selected profile and region
			_, err := providers.CreateProvider(newModel.Registry, "AWS", newModel.AwsProfile, newModel.AwsRegion)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return core.ErrMsg{Err: err}
				}
			}

			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(newModel)
		}
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleServiceSelection handles the selection of an AWS service
func HandleServiceSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		serviceName := selected[0]

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return core.ErrMsg{Err: err}
			}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == serviceName {
				selectedService = service
				break
			}
		}

		if selectedService != nil {
			newModel := m.Clone()
			newModel.SelectedService = &core.Service{
				Name:        selectedService.Name(),
				Description: selectedService.Description(),
			}
			newModel.CurrentView = constants.ViewSelectCategory
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleCategorySelection handles the selection of a service category
func HandleCategorySelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		categoryName := selected[0]

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return core.ErrMsg{Err: err}
			}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return WrapModel(m), nil
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == categoryName {
				selectedCategory = category
				break
			}
		}

		if selectedCategory != nil {
			newModel := m.Clone()
			newModel.SelectedCategory = &core.Category{
				Name:        selectedCategory.Name(),
				Description: selectedCategory.Description(),
			}
			newModel.CurrentView = constants.ViewSelectOperation
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleOperationSelection handles the selection of a service operation
func HandleOperationSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		operationName := selected[0]

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return core.ErrMsg{Err: err}
			}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return WrapModel(m), nil
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == m.SelectedCategory.Name {
				selectedCategory = category
				break
			}
		}

		if selectedCategory == nil {
			return WrapModel(m), nil
		}

		// Find the selected operation
		var selectedOperation providers.Operation
		for _, operation := range selectedCategory.Operations() {
			if operation.Name() == operationName {
				selectedOperation = operation
				break
			}
		}

		if selectedOperation != nil {
			newModel := m.Clone()
			newModel.SelectedOperation = &core.Operation{
				Name:        selectedOperation.Name(),
				Description: selectedOperation.Description(),
			}

			// Handle different operations
			switch operationName {
			case "Pipeline Approvals":
				newModel.IsLoading = true
				newModel.LoadingMsg = "Fetching approvals..."
				return WrapModel(newModel), FetchApprovals(newModel)
			case "Pipeline Status":
				newModel.IsLoading = true
				newModel.LoadingMsg = "Fetching pipeline status..."
				return WrapModel(newModel), FetchPipelineStatus(newModel)
			case "Start Pipeline":
				newModel.CurrentView = constants.ViewPipelineStatus
				newModel.IsLoading = true
				newModel.LoadingMsg = "Fetching pipelines..."
				return WrapModel(newModel), FetchPipelineStatus(newModel)
			default:
				newModel.CurrentView = constants.ViewConfirmation
				view.UpdateTableForView(newModel)
				return WrapModel(newModel), nil
			}
		}
	}
	return WrapModel(m), nil
}

// HandleApprovalSelection handles the selection of a pipeline approval
func HandleApprovalSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		for _, approval := range m.Approvals {
			if approval.PipelineName == selected[0] &&
				approval.StageName == selected[1] &&
				approval.ActionName == selected[2] {
				newModel.SelectedApproval = &approval
				newModel.CurrentView = constants.ViewConfirmation
				view.UpdateTableForView(newModel)
				return WrapModel(newModel), nil
			}
		}
	}
	return WrapModel(m), nil
}

// HandleConfirmationSelection handles the selection of an approval action
func HandleConfirmationSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		if selected[0] == "Approve" {
			newModel.ApproveAction = true
			newModel.CurrentView = constants.ViewSummary
			newModel.ManualInput = true
			newModel.SetTextInputForApproval(true)
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		} else if selected[0] == "Reject" {
			newModel.ApproveAction = false
			newModel.CurrentView = constants.ViewSummary
			newModel.ManualInput = true
			newModel.SetTextInputForApproval(false)
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleSummaryConfirmation handles the confirmation of the summary
func HandleSummaryConfirmation(m *core.Model) (tea.Model, tea.Cmd) {
	if m.ManualInput {
		// For manual input, just store the value and continue
		newModel := m.Clone()

		// Store the comment
		if m.SelectedApproval != nil {
			newModel.ApprovalComment = m.TextInput.Value()
			newModel.Summary = m.TextInput.Value()
		}

		// For pipeline execution with manual commit ID
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			newModel.CommitID = m.TextInput.Value()
			newModel.ManualCommitID = true
		}

		// Move to execution view
		newModel.CurrentView = constants.ViewExecutingAction
		newModel.ManualInput = false
		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}

	// For non-manual input, check if we have a selected row
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		if selected[0] == "Execute" {
			// Start loading and execute the action
			newModel.IsLoading = true
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				return WrapModel(newModel), ExecutePipeline(m)
			}
			return WrapModel(newModel), ExecuteApproval(m)
		} else if selected[0] == "Cancel" {
			// Navigate back to the main menu
			newModel.CurrentView = constants.ViewSelectOperation
			newModel.SelectedApproval = nil
			newModel.SelectedPipeline = nil
			newModel.ApprovalComment = ""
			newModel.CommitID = ""
			newModel.ManualCommitID = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleExecutionSelection handles the selection of an execution action
func HandleExecutionSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		if selected[0] == "Execute" {
			// Start loading and execute the action
			newModel.IsLoading = true
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				return WrapModel(newModel), ExecutePipeline(m)
			}
			return WrapModel(newModel), ExecuteApproval(m)
		} else if selected[0] == "Cancel" {
			// Navigate back to the main menu
			newModel.CurrentView = constants.ViewSelectOperation
			newModel.SelectedApproval = nil
			newModel.SelectedPipeline = nil
			newModel.ApprovalComment = ""
			newModel.CommitID = ""
			newModel.ManualCommitID = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// Async operations

// FetchApprovals fetches pipeline approvals from AWS
func FetchApprovals(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected service not found")}
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == m.SelectedCategory.Name {
				selectedCategory = category
				break
			}
		}

		if selectedCategory == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected category not found")}
		}

		// Find the selected operation
		var selectedOperation providers.Operation
		for _, operation := range selectedCategory.Operations() {
			if operation.Name() == m.SelectedOperation.Name {
				selectedOperation = operation
				break
			}
		}

		if selectedOperation == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected operation not found")}
		}

		// Execute the operation
		ctx := context.Background()
		result, err := selectedOperation.Execute(ctx, nil)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Convert the result to approvals
		// The result could be of different types depending on the implementation
		var approvals []aws.ApprovalAction

		// Try to convert from []interface{} first
		if approvalsInterface, ok := result.([]interface{}); ok {
			approvals = make([]aws.ApprovalAction, len(approvalsInterface))
			for i, a := range approvalsInterface {
				if approval, ok := a.(map[string]interface{}); ok {
					approvals[i] = aws.ApprovalAction{
						PipelineName: fmt.Sprintf("%v", approval["PipelineName"]),
						StageName:    fmt.Sprintf("%v", approval["StageName"]),
						ActionName:   fmt.Sprintf("%v", approval["ActionName"]),
						Token:        fmt.Sprintf("%v", approval["Token"]),
					}
				}
			}
		} else if cloudApprovals, ok := result.([]struct {
			PipelineName string
			StageName    string
			ActionName   string
			Token        string
		}); ok {
			// Try to convert from a slice of structs
			approvals = make([]aws.ApprovalAction, len(cloudApprovals))
			for i, a := range cloudApprovals {
				approvals[i] = aws.ApprovalAction{
					PipelineName: a.PipelineName,
					StageName:    a.StageName,
					ActionName:   a.ActionName,
					Token:        a.Token,
				}
			}
		} else {
			// Try to convert using our helper function
			convertedApprovals, err := convertCodePipelineApprovals(result)
			if err != nil {
				return core.ErrMsg{Err: fmt.Errorf("unexpected result type: %T - %v", result, err)}
			}
			approvals = convertedApprovals
		}

		return core.ApprovalsMsg{
			Approvals: approvals,
		}
	}
}

// convertCodePipelineStatus attempts to convert a slice of codepipeline.PipelineStatus to a slice of aws.PipelineStatus
func convertCodePipelineStatus(result interface{}) ([]aws.PipelineStatus, error) {
	// Get the type name as a string
	typeName := fmt.Sprintf("%T", result)
	if !strings.Contains(typeName, "PipelineStatus") {
		return nil, fmt.Errorf("not a PipelineStatus type: %s", typeName)
	}

	// Use reflection to access the slice
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("not a slice type: %s", typeName)
	}

	// Create the output slice
	length := resultValue.Len()
	pipelines := make([]aws.PipelineStatus, length)

	// Process each pipeline
	for i := 0; i < length; i++ {
		pipelineValue := resultValue.Index(i)

		// Extract the Name field
		nameField := pipelineValue.FieldByName("Name")
		if !nameField.IsValid() {
			return nil, fmt.Errorf("Name field not found in PipelineStatus")
		}

		// Convert the name to string
		var name string
		switch nameField.Kind() {
		case reflect.String:
			name = nameField.String()
		case reflect.Ptr:
			if !nameField.IsNil() {
				name = nameField.Elem().String()
			}
		default:
			return nil, fmt.Errorf("Name field is not a string or string pointer")
		}

		// Create the pipeline status
		pipelines[i] = aws.PipelineStatus{
			Name: name,
		}

		// Extract the Stages field
		stagesField := pipelineValue.FieldByName("Stages")
		if !stagesField.IsValid() || stagesField.Kind() != reflect.Slice {
			continue // Skip stages if not found or not a slice
		}

		// Process each stage
		stagesLen := stagesField.Len()
		stages := make([]aws.StageStatus, stagesLen)

		for j := 0; j < stagesLen; j++ {
			stageValue := stagesField.Index(j)

			// Extract stage fields
			nameField := stageValue.FieldByName("Name")
			statusField := stageValue.FieldByName("Status")
			lastUpdatedField := stageValue.FieldByName("LastUpdated")

			// Convert name to string
			var stageName string
			if nameField.IsValid() {
				switch nameField.Kind() {
				case reflect.String:
					stageName = nameField.String()
				case reflect.Ptr:
					if !nameField.IsNil() {
						stageName = nameField.Elem().String()
					}
				}
			}

			// Convert status to string
			var stageStatus string
			if statusField.IsValid() {
				switch statusField.Kind() {
				case reflect.String:
					stageStatus = statusField.String()
				case reflect.Ptr:
					if !statusField.IsNil() {
						stageStatus = statusField.Elem().String()
					}
				}
			}

			// Convert lastUpdated to string
			var lastUpdated string
			if lastUpdatedField.IsValid() {
				switch lastUpdatedField.Kind() {
				case reflect.String:
					lastUpdated = lastUpdatedField.String()
				case reflect.Ptr:
					if !lastUpdatedField.IsNil() {
						lastUpdated = lastUpdatedField.Elem().String()
					}
				}
			}

			// Create the stage status
			stages[j] = aws.StageStatus{
				Name:        stageName,
				Status:      stageStatus,
				LastUpdated: lastUpdated,
			}
		}

		// Assign stages to the pipeline
		pipelines[i].Stages = stages
	}

	return pipelines, nil
}

// convertCodePipelineApprovals attempts to convert a slice of codepipeline.ApprovalAction to a slice of aws.ApprovalAction
func convertCodePipelineApprovals(result interface{}) ([]aws.ApprovalAction, error) {
	// Get the type name as a string
	typeName := fmt.Sprintf("%T", result)
	if !strings.Contains(typeName, "ApprovalAction") {
		return nil, fmt.Errorf("not an ApprovalAction type: %s", typeName)
	}

	// Use reflection to access the slice
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("not a slice type: %s", typeName)
	}

	// Create the output slice
	length := resultValue.Len()
	approvals := make([]aws.ApprovalAction, length)

	// Process each approval
	for i := 0; i < length; i++ {
		approvalValue := resultValue.Index(i)

		// Extract fields
		pipelineNameField := approvalValue.FieldByName("PipelineName")
		stageNameField := approvalValue.FieldByName("StageName")
		actionNameField := approvalValue.FieldByName("ActionName")
		tokenField := approvalValue.FieldByName("Token")

		// Convert fields to strings
		var pipelineName, stageName, actionName, token string

		if pipelineNameField.IsValid() {
			switch pipelineNameField.Kind() {
			case reflect.String:
				pipelineName = pipelineNameField.String()
			case reflect.Ptr:
				if !pipelineNameField.IsNil() {
					pipelineName = pipelineNameField.Elem().String()
				}
			}
		}

		if stageNameField.IsValid() {
			switch stageNameField.Kind() {
			case reflect.String:
				stageName = stageNameField.String()
			case reflect.Ptr:
				if !stageNameField.IsNil() {
					stageName = stageNameField.Elem().String()
				}
			}
		}

		if actionNameField.IsValid() {
			switch actionNameField.Kind() {
			case reflect.String:
				actionName = actionNameField.String()
			case reflect.Ptr:
				if !actionNameField.IsNil() {
					actionName = actionNameField.Elem().String()
				}
			}
		}

		if tokenField.IsValid() {
			switch tokenField.Kind() {
			case reflect.String:
				token = tokenField.String()
			case reflect.Ptr:
				if !tokenField.IsNil() {
					token = tokenField.Elem().String()
				}
			}
		}

		// Create the approval action
		approvals[i] = aws.ApprovalAction{
			PipelineName: pipelineName,
			StageName:    stageName,
			ActionName:   actionName,
			Token:        token,
		}
	}

	return approvals, nil
}

// FetchPipelineStatus fetches pipeline status from AWS
func FetchPipelineStatus(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected service not found")}
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == m.SelectedCategory.Name {
				selectedCategory = category
				break
			}
		}

		if selectedCategory == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected category not found")}
		}

		// Find the operation for pipeline status
		var statusOperation providers.Operation
		for _, operation := range selectedCategory.Operations() {
			if operation.Name() == "Pipeline Status" {
				statusOperation = operation
				break
			}
		}

		if statusOperation == nil {
			return core.ErrMsg{Err: fmt.Errorf("pipeline status operation not found")}
		}

		// Execute the operation
		ctx := context.Background()
		result, err := statusOperation.Execute(ctx, nil)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Convert the result to pipeline status
		// The result could be of different types depending on the implementation
		var pipelines []aws.PipelineStatus

		// Try to convert from []interface{} first
		if pipelinesInterface, ok := result.([]interface{}); ok {
			pipelines = make([]aws.PipelineStatus, len(pipelinesInterface))
			for i, p := range pipelinesInterface {
				if pipeline, ok := p.(map[string]interface{}); ok {
					// Create a new pipeline status
					pipelines[i] = aws.PipelineStatus{
						Name: fmt.Sprintf("%v", pipeline["Name"]),
					}

					// Convert stages if they exist
					if stagesInterface, ok := pipeline["Stages"].([]interface{}); ok {
						stages := make([]aws.StageStatus, len(stagesInterface))
						for j, s := range stagesInterface {
							if stage, ok := s.(map[string]interface{}); ok {
								stages[j] = aws.StageStatus{
									Name:        fmt.Sprintf("%v", stage["Name"]),
									Status:      fmt.Sprintf("%v", stage["Status"]),
									LastUpdated: fmt.Sprintf("%v", stage["LastUpdated"]),
								}
							}
						}
						pipelines[i].Stages = stages
					}
				}
			}
		} else if cloudPipelines, ok := result.([]struct {
			Name   string
			Stages []struct {
				Name        string
				Status      string
				LastUpdated string
			}
		}); ok {
			// Try to convert from a slice of structs
			pipelines = make([]aws.PipelineStatus, len(cloudPipelines))
			for i, p := range cloudPipelines {
				pipelines[i] = aws.PipelineStatus{
					Name: p.Name,
				}

				// Convert stages
				stages := make([]aws.StageStatus, len(p.Stages))
				for j, s := range p.Stages {
					stages[j] = aws.StageStatus{
						Name:        s.Name,
						Status:      s.Status,
						LastUpdated: s.LastUpdated,
					}
				}
				pipelines[i].Stages = stages
			}
		} else {
			// Try to convert using our helper function
			convertedPipelines, err := convertCodePipelineStatus(result)
			if err != nil {
				return core.ErrMsg{Err: fmt.Errorf("unexpected result type: %T - %v", result, err)}
			}
			pipelines = convertedPipelines
		}

		return core.PipelineStatusMsg{
			Pipelines: pipelines,
		}
	}
}

// ExecuteApproval executes an approval action
func ExecuteApproval(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedApproval == nil {
			return core.ErrMsg{Err: fmt.Errorf("no approval selected")}
		}

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Find the CodePipeline service
		var codePipelineService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == "CodePipeline" {
				codePipelineService = service
				break
			}
		}

		if codePipelineService == nil {
			return core.ErrMsg{Err: fmt.Errorf("CodePipeline service not found")}
		}

		// Find the InternalOperations category
		var internalCategory providers.Category
		for _, category := range codePipelineService.Categories() {
			if category.Name() == "InternalOperations" {
				internalCategory = category
				break
			}
		}

		if internalCategory == nil {
			return core.ErrMsg{Err: fmt.Errorf("InternalOperations category not found")}
		}

		// Find the approval operation
		var approvalOperation providers.Operation
		for _, operation := range internalCategory.Operations() {
			if operation.Name() == "Approval" {
				approvalOperation = operation
				break
			}
		}

		if approvalOperation == nil {
			return core.ErrMsg{Err: fmt.Errorf("approval operation not found")}
		}

		// Prepare parameters for the operation
		params := map[string]interface{}{
			"pipeline_name": m.SelectedApproval.PipelineName,
			"stage_name":    m.SelectedApproval.StageName,
			"action_name":   m.SelectedApproval.ActionName,
			"token":         m.SelectedApproval.Token,
			"approved":      m.ApproveAction,
			"comment":       m.ApprovalComment,
		}

		// Execute the operation
		ctx := context.Background()
		_, err = approvalOperation.Execute(ctx, params)
		if err != nil {
			return core.ApprovalResultMsg{Err: err}
		}

		return core.ApprovalResultMsg{Err: nil}
	}
}

// ExecutePipeline executes a pipeline
func ExecutePipeline(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedPipeline == nil {
			return core.ErrMsg{Err: fmt.Errorf("no pipeline selected")}
		}

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Find the CodePipeline service
		var codePipelineService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == "CodePipeline" {
				codePipelineService = service
				break
			}
		}

		if codePipelineService == nil {
			return core.ErrMsg{Err: fmt.Errorf("CodePipeline service not found")}
		}

		// Find the Workflows category
		var workflowsCategory providers.Category
		for _, category := range codePipelineService.Categories() {
			if category.Name() == "Workflows" {
				workflowsCategory = category
				break
			}
		}

		if workflowsCategory == nil {
			return core.ErrMsg{Err: fmt.Errorf("Workflows category not found")}
		}

		// Find the start pipeline operation
		var startPipelineOperation providers.Operation
		for _, operation := range workflowsCategory.Operations() {
			if operation.Name() == "Start Pipeline" {
				startPipelineOperation = operation
				break
			}
		}

		if startPipelineOperation == nil {
			return core.ErrMsg{Err: fmt.Errorf("start pipeline operation not found")}
		}

		// Prepare parameters for the operation
		params := map[string]interface{}{
			"pipeline_name": m.SelectedPipeline.Name,
		}

		// Add commit ID if specified
		if m.ManualCommitID && m.CommitID != "" {
			params["commit_id"] = m.CommitID
		}

		// Execute the operation
		ctx := context.Background()
		_, err = startPipelineOperation.Execute(ctx, params)
		if err != nil {
			return core.PipelineExecutionMsg{Err: err}
		}

		return core.PipelineExecutionMsg{Err: nil}
	}
}
