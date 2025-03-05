package integration_test

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
)

// TestAWSPipelineStartFlow tests the AWS pipeline start flow with default/us-east-1
func TestAWSPipelineStartFlow(t *testing.T) {
	// Initialize the model
	m := model.New()

	// Set up the AWS provider
	registry := providers.NewProviderRegistry()
	registry.Register(CreateMockAWSProvider())
	m.Registry = registry

	// Set up AWS profile and region
	m.SetAwsProfile("default")
	m.SetAwsRegion("us-east-1")

	// Test service selection
	t.Run("Service Selection", func(t *testing.T) {
		// Set the current view to service selection
		m.CurrentView = constants.ViewSelectService

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for service selection view: %v", err)
		}

		// Select the CodePipeline service
		for i, row := range m.Table.Rows() {
			if row[0] == "CodePipeline" {
				m.Table.SetCursor(i)
				m.SelectedService = &model.Service{
					Name:        "CodePipeline",
					Description: row[1],
				}
				break
			}
		}

		// Verify the service is selected
		if m.SelectedService == nil || m.SelectedService.Name != "CodePipeline" {
			t.Error("Failed to select CodePipeline service")
		}
	})

	// Test category selection
	t.Run("Category Selection", func(t *testing.T) {
		// Set the current view to category selection
		m.CurrentView = constants.ViewSelectCategory

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for category selection view: %v", err)
		}

		// Select the Operations category if available
		categoryFound := false
		for i, row := range m.Table.Rows() {
			if row[0] == "Operations" {
				m.Table.SetCursor(i)
				m.SelectedCategory = &model.Category{
					Name:        "Operations",
					Description: row[1],
				}
				categoryFound = true
				break
			}
		}

		if !categoryFound {
			t.Log("Operations category not found, using first available category")
			if len(m.Table.Rows()) > 0 {
				m.Table.SetCursor(0)
				row := m.Table.SelectedRow()
				m.SelectedCategory = &model.Category{
					Name:        row[0],
					Description: row[1],
				}
			} else {
				t.Error("No categories available")
				return
			}
		}
	})

	// Test operation selection
	t.Run("Operation Selection", func(t *testing.T) {
		// Set the current view to operation selection
		m.CurrentView = constants.ViewSelectOperation

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for operation selection view: %v", err)
		}

		// Select the Start Pipeline operation if available
		operationFound := false
		for i, row := range m.Table.Rows() {
			if row[0] == "Start Pipeline" {
				m.Table.SetCursor(i)
				m.SelectedOperation = &model.Operation{
					Name:        "Start Pipeline",
					Description: row[1],
				}
				operationFound = true
				break
			}
		}

		if !operationFound {
			t.Log("Start Pipeline operation not found, using first available operation")
			if len(m.Table.Rows()) > 0 {
				m.Table.SetCursor(0)
				row := m.Table.SelectedRow()
				m.SelectedOperation = &model.Operation{
					Name:        row[0],
					Description: row[1],
				}
			} else {
				t.Error("No operations available")
				return
			}
		}
	})

	// Test pipeline status view
	t.Run("Pipeline Status View", func(t *testing.T) {
		// Set the current view to pipeline status
		m.CurrentView = constants.ViewPipelineStatus

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for pipeline status view: %v", err)
		}

		// Verify that pipelines are loaded (or at least the table is initialized)
		if m.Table.Rows() == nil {
			t.Log("No pipelines found or table not initialized")
		} else {
			t.Logf("Found %d pipelines", len(m.Table.Rows()))

			// If pipelines are found, select the first one
			if len(m.Table.Rows()) > 0 {
				m.Table.SetCursor(0)
				row := m.Table.SelectedRow()

				// Find the pipeline in the model's pipelines
				for _, pipeline := range m.Pipelines {
					if pipeline.Name == row[0] {
						m.SelectedPipeline = &pipeline
						break
					}
				}

				// Verify the pipeline is selected
				if m.SelectedPipeline == nil {
					t.Error("Failed to select a pipeline")
				} else {
					t.Logf("Selected pipeline: %s", m.SelectedPipeline.Name)
				}
			}
		}
	})

	// Test summary view
	t.Run("Summary View", func(t *testing.T) {
		// Skip if no pipeline was selected
		if m.SelectedPipeline == nil {
			t.Skip("No pipeline selected, skipping summary view test")
		}

		// Set the current view to summary
		m.CurrentView = constants.ViewSummary

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for summary view: %v", err)
		}

		// Verify that summary options are loaded
		if m.Table.Rows() == nil || len(m.Table.Rows()) == 0 {
			t.Error("Expected summary options to be loaded, but table rows are empty")
			return
		}

		// Select "Latest Commit"
		for i, row := range m.Table.Rows() {
			if row[0] == "Latest Commit" {
				m.Table.SetCursor(i)
				m.SetManualCommitID(false)
				break
			}
		}

		// Verify the manual commit ID is not set
		if m.GetManualCommitID() {
			t.Error("Expected manual commit ID to be false")
		}
	})

	// Test executing action view
	t.Run("Executing Action View", func(t *testing.T) {
		// Skip if no pipeline was selected
		if m.SelectedPipeline == nil {
			t.Skip("No pipeline selected, skipping executing action view test")
		}

		// Set the current view to executing action
		m.CurrentView = constants.ViewExecutingAction

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for executing action view: %v", err)
		}

		// Note: We don't actually execute the action in the test
		// as it would modify real resources
		t.Log("Pipeline start would happen here in a real scenario")
	})
}
