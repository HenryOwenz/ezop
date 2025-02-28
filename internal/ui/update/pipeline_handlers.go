package update

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleApprovalSelection handles the selection of a pipeline approval
func HandleApprovalSelection(m *model.Model) (tea.Model, tea.Cmd) {
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

// HandlePipelineSelection handles the selection of a pipeline
func HandlePipelineSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()
		for _, pipeline := range m.Pipelines {
			if pipeline.Name == selected[0] {
				newModel.SelectedPipeline = &pipeline
				if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
					newModel.CurrentView = constants.ViewSummary
				} else {
					newModel.CurrentView = constants.ViewPipelineStages
				}
				view.UpdateTableForView(newModel)
				return WrapModel(newModel), nil
			}
		}
	}
	return WrapModel(m), nil
}
