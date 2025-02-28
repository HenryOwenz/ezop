package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
)

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
					// Create a new approval action
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
