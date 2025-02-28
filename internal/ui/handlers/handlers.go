package handlers

import (
	"context"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/aws"
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
		if selected[0] == "Amazon Web Services" {
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
		newModel := m.Clone()
		newModel.SelectedService = &core.Service{
			Name:        selected[0],
			Description: selected[1],
		}
		newModel.CurrentView = constants.ViewSelectCategory
		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleCategorySelection handles the selection of a service category
func HandleCategorySelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		newModel.SelectedCategory = &core.Category{
			Name:        selected[0],
			Description: selected[1],
		}
		newModel.CurrentView = constants.ViewSelectOperation
		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleOperationSelection handles the selection of a service operation
func HandleOperationSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		newModel.SelectedOperation = &core.Operation{
			Name:        selected[0],
			Description: selected[1],
		}

		if selected[0] == "Pipeline Approvals" {
			// Start loading approvals
			newModel.IsLoading = true
			newModel.LoadingMsg = constants.MsgLoadingApprovals
			return WrapModel(newModel), FetchApprovals(m)
		} else if selected[0] == "Pipeline Status" || selected[0] == "Start Pipeline" {
			// Start loading pipeline status
			newModel.IsLoading = true
			newModel.LoadingMsg = constants.MsgLoadingPipelines
			return WrapModel(newModel), FetchPipelineStatus(m)
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

// FetchApprovals fetches pipeline approvals
func FetchApprovals(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		approvals, err := provider.GetPendingApprovals(ctx)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		return core.ApprovalsMsg{
			Provider:  provider,
			Approvals: approvals,
		}
	}
}

// FetchPipelineStatus fetches pipeline status
func FetchPipelineStatus(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		pipelines, err := provider.GetPipelineStatus(ctx)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		return core.PipelineStatusMsg{
			Provider:  provider,
			Pipelines: pipelines,
		}
	}
}

// ExecuteApproval executes an approval action
func ExecuteApproval(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.ApprovalResultMsg{Err: err}
		}

		err = provider.PutApprovalResult(ctx, *m.SelectedApproval, m.ApproveAction, m.ApprovalComment)
		return core.ApprovalResultMsg{Err: err}
	}
}

// ExecutePipeline executes a pipeline
func ExecutePipeline(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.PipelineExecutionMsg{Err: err}
		}

		commitID := ""
		if m.ManualCommitID {
			commitID = m.CommitID
		}

		err = provider.StartPipelineExecution(ctx, m.SelectedPipeline.Name, commitID)
		return core.PipelineExecutionMsg{Err: err}
	}
}
