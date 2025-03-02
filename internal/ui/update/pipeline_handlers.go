package update

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectApproval handles the selection of a pipeline approval
func SelectApproval(m *model.Model) (tea.Model, tea.Cmd) {
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
					newModel.CurrentView = constants.ViewExecutingAction
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

// HandlePipelineApprovals handles the pipeline approvals operation
func HandlePipelineApprovals(m *model.Model) (tea.Model, tea.Cmd) {
	newModel := m.Clone()
	newModel.IsLoading = true
	newModel.LoadingMsg = constants.MsgLoadingApprovals

	return WrapModel(newModel), func() tea.Msg {
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the CodePipelineManualApprovalOperation from the provider
		approvalOperation, err := provider.GetCodePipelineManualApprovalOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get approvals using the operation
		ctx := context.Background()
		approvals, err := approvalOperation.GetPendingApprovals(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		return model.ApprovalsMsg{
			Approvals: approvals,
			Provider:  provider,
		}
	}
}

// HandlePipelineStatus handles the pipeline status operation
func HandlePipelineStatus(m *model.Model) (tea.Model, tea.Cmd) {
	newModel := m.Clone()
	newModel.IsLoading = true
	newModel.LoadingMsg = constants.MsgLoadingPipelines

	return WrapModel(newModel), func() tea.Msg {
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the PipelineStatusOperation from the provider
		statusOperation, err := provider.GetPipelineStatusOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get pipeline status using the operation
		ctx := context.Background()
		pipelines, err := statusOperation.GetPipelineStatus(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		return model.PipelineStatusMsg{
			Pipelines: pipelines,
			Provider:  provider,
		}
	}
}
