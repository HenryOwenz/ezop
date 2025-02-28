package update

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// HandleApprovalResult handles the result of an approval action
func HandleApprovalResult(m *model.Model, err error) {
	if err != nil {
		m.Error = fmt.Sprintf(constants.MsgErrorGeneric, err.Error())
		m.CurrentView = constants.ViewError
		return
	}

	// Use the appropriate message constant based on approval action
	if m.ApproveAction {
		m.Success = fmt.Sprintf(constants.MsgApprovalSuccess,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)
	} else {
		m.Success = fmt.Sprintf(constants.MsgRejectionSuccess,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)
	}

	// Reset approval state
	m.SelectedApproval = nil
	m.ApprovalComment = ""

	// Completely reset the text input
	m.ResetTextInput()
	m.TextInput.Placeholder = constants.MsgEnterComment
	m.ManualInput = false

	// Navigate back to the operation selection view
	m.CurrentView = constants.ViewSelectOperation

	// Clear the approvals list to force a refresh next time
	m.Approvals = nil

	// Update the table for the current view
	view.UpdateTableForView(m)
}

// FetchApprovals fetches pipeline approvals from the provider
func FetchApprovals(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return model.ErrMsg{Err: err}
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
			return model.ErrMsg{Err: fmt.Errorf("selected service not found")}
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
			return model.ErrMsg{Err: fmt.Errorf("selected category not found")}
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
			return model.ErrMsg{Err: fmt.Errorf("selected operation not found")}
		}

		// We don't need to execute the operation anymore, as we'll get approvals directly from the provider
		// Just keeping this code to maintain the validation of service/category/operation

		// Get approvals directly from the provider
		ctx := context.Background()
		approvals, err := provider.GetApprovals(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		return model.ApprovalsMsg{
			Approvals: approvals,
			Provider:  provider,
		}
	}
}

// ExecuteApproval executes an approval action
func ExecuteApproval(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedApproval == nil {
			return model.ErrMsg{Err: fmt.Errorf("no approval selected")}
		}

		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Convert the legacy ApprovalAction to providers.ApprovalAction
		providerApproval := providers.ApprovalAction{
			PipelineName: m.SelectedApproval.PipelineName,
			StageName:    m.SelectedApproval.StageName,
			ActionName:   m.SelectedApproval.ActionName,
			Token:        m.SelectedApproval.Token,
		}

		// Execute the approval action
		ctx := context.Background()
		err = provider.ApproveAction(ctx, providerApproval, m.ApproveAction, m.ApprovalComment)
		if err != nil {
			return model.ApprovalResultMsg{Err: err}
		}

		return model.ApprovalResultMsg{Err: nil}
	}
}
