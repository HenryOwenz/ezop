package update

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// HandlePipelineExecution handles the result of a pipeline execution
func HandlePipelineExecution(m *model.Model, err error) {
	if err != nil {
		m.Error = fmt.Sprintf(constants.MsgErrorGeneric, err.Error())
		m.CurrentView = constants.ViewError
		return
	}

	m.Success = fmt.Sprintf(constants.MsgPipelineStartSuccess, m.SelectedPipeline.Name)

	// Reset pipeline state
	m.SelectedPipeline = nil
	m.CommitID = ""
	m.ManualCommitID = false

	// Completely reset the text input
	m.ResetTextInput()
	m.TextInput.Placeholder = constants.MsgEnterComment
	m.ManualInput = false

	// Navigate back to the operation selection view
	m.CurrentView = constants.ViewSelectOperation

	// Clear the pipelines list to force a refresh next time
	m.Pipelines = nil

	// Update the table for the current view
	view.UpdateTableForView(m)
}

// FetchPipelineStatus fetches pipeline status from the provider
func FetchPipelineStatus(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
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

// ExecutePipeline executes a pipeline
func ExecutePipeline(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedPipeline == nil {
			return model.ErrMsg{Err: fmt.Errorf("no pipeline selected")}
		}

		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the StartPipelineOperation from the provider
		startOperation, err := provider.GetStartPipelineOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Execute the pipeline using the operation
		ctx := context.Background()
		err = startOperation.StartPipelineExecution(ctx, m.SelectedPipeline.Name, m.CommitID)
		if err != nil {
			return model.PipelineExecutionMsg{Err: err}
		}

		return model.PipelineExecutionMsg{Err: nil}
	}
}
