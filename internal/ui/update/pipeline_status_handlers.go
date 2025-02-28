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

		// Find the operation for pipeline status
		var statusOperation providers.Operation
		for _, operation := range selectedCategory.Operations() {
			if operation.Name() == "Pipeline Status" {
				statusOperation = operation
				break
			}
		}

		if statusOperation == nil {
			return model.ErrMsg{Err: fmt.Errorf("pipeline status operation not found")}
		}

		// We don't need to execute the operation anymore, as we'll get pipeline status directly from the provider
		// Just keeping this code to maintain the validation of service/category/operation

		// Get pipeline status directly from the provider
		ctx := context.Background()
		pipelines, err := provider.GetStatus(ctx)
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

		// Execute the pipeline
		ctx := context.Background()
		err = provider.StartPipeline(ctx, m.SelectedPipeline.Name, m.CommitID)
		if err != nil {
			return model.PipelineExecutionMsg{Err: err}
		}

		return model.PipelineExecutionMsg{Err: nil}
	}
}
