package update

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	"github.com/charmbracelet/bubbles/table"
)

// TestCompletePipelineStartFlow tests the complete flow for starting a pipeline,
// from selecting a pipeline to executing the action.
// This test will fail if someone changes the expected behavior of the flow.
func TestCompletePipelineStartFlow(t *testing.T) {
	// Step 1: Create initial model with operation selected
	m := model.New()
	m.CurrentView = constants.ViewSelectOperation
	m.SelectedService = &model.Service{Name: "CodePipeline"}
	m.SelectedCategory = &model.Category{Name: "Operations"}
	m.SelectedOperation = &model.Operation{Name: "Start Pipeline"}

	// Step 2: Set up pipelines view (simulating HandlePipelineStatus)
	m.CurrentView = constants.ViewPipelineStatus
	m.Pipelines = []providers.PipelineStatus{
		{
			Name:   "TestPipeline",
			Stages: []providers.StageStatus{},
		},
	}
	view.UpdateTableForView(m)

	// Step 3: Select a pipeline (simulating HandlePipelineSelection)
	// This would normally be done by selecting a row in the table
	m.SelectedPipeline = &m.Pipelines[0]

	// Create a table with a single row for the pipeline
	columns := []table.Column{
		{Title: "Pipeline", Width: 20},
		{Title: "Stages", Width: 10},
	}
	rows := []table.Row{
		{m.SelectedPipeline.Name, "3"},
	}
	m.Table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	m.Table.SetCursor(0)

	result, cmd := HandlePipelineSelection(m)

	// Verify we get a model wrapper and no command
	wrapper, ok := result.(ModelWrapper)
	if !ok {
		t.Fatalf("Expected HandlePipelineSelection to return a ModelWrapper, got %T", result)
	}
	if cmd != nil {
		t.Errorf("Expected HandlePipelineSelection to return nil command, got %T", cmd)
	}

	// Verify we're now at the execution view directly (not the summary view)
	if wrapper.Model.CurrentView != constants.ViewExecutingAction {
		t.Errorf("Expected to be at ViewExecutingAction, got %v", wrapper.Model.CurrentView)
	}

	// Step 4: Execute the pipeline start (simulating HandleExecutionSelection)
	// Set up the table with execution options
	columns = []table.Column{
		{Title: "Action", Width: 10},
		{Title: "Description", Width: 30},
	}
	rows = []table.Row{
		{"Execute", "Start pipeline with latest commit"},
		{"Cancel", "Cancel and return to main menu"},
	}

	wrapper.Model.Table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	// Select the "Execute" row (index 0)
	wrapper.Model.Table.SetCursor(0)

	result, cmd = HandleExecutionSelection(wrapper.Model)

	// Verify we get a model wrapper and a command
	wrapper, ok = result.(ModelWrapper)
	if !ok {
		t.Fatalf("Expected HandleExecutionSelection to return a ModelWrapper, got %T", result)
	}
	if cmd == nil {
		t.Errorf("Expected HandleExecutionSelection to return a command")
	}

	// Verify we're now in loading state
	if !wrapper.Model.IsLoading {
		t.Errorf("Expected IsLoading to be true")
	}

	// Step 5: Verify that navigating back from executing action view
	// The NavigateBack function should go to the pipeline status view
	// when the selected operation is "Start Pipeline"
	backResult := NavigateBack(wrapper.Model)

	// Check if we're at the pipeline status view
	if backResult.CurrentView != constants.ViewPipelineStatus {
		t.Errorf("Expected to navigate back to ViewPipelineStatus, got %v", backResult.CurrentView)
	}
}
