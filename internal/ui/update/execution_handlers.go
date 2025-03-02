package update

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleConfirmationSelection handles the selection of an approval action
func HandleConfirmationSelection(m *model.Model) (tea.Model, tea.Cmd) {
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
func HandleSummaryConfirmation(m *model.Model) (tea.Model, tea.Cmd) {
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
		newModel.ResetTextInput()
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
func HandleExecutionSelection(m *model.Model) (tea.Model, tea.Cmd) {
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
