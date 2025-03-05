package integration_test

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// TestAWSOperationsFlow tests the AWS operations flow with default/us-east-1
func TestAWSOperationsFlow(t *testing.T) {
	// Initialize the model
	m := model.New()

	// Set up the AWS provider
	registry := providers.NewProviderRegistry()
	registry.Register(CreateMockAWSProvider())
	m.Registry = registry

	// Set up AWS profile and region
	t.Run("Setup AWS Configuration", func(t *testing.T) {
		// Set the AWS profile and region
		m.SetAwsProfile("default")
		m.SetAwsRegion("us-east-1")

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for AWS config view: %v", err)
		}

		// Update the table for the view
		view.UpdateTableForView(m)

		// Verify the profile is set correctly
		if m.GetAwsProfile() != "default" {
			t.Errorf("Expected AWS profile to be 'default', got '%s'", m.GetAwsProfile())
		}

		// Create the provider with the selected profile and region
		provider, err := providers.CreateProvider(m.Registry, "AWS", m.GetAwsProfile(), m.GetAwsRegion())
		if err != nil {
			t.Fatalf("Failed to create AWS provider: %v", err)
		}

		// Set the provider in the model
		m.Provider = provider
	})

	// Test service selection
	t.Run("Service Selection", func(t *testing.T) {
		// Set the current view to service selection
		m.CurrentView = constants.ViewSelectService

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for service selection view: %v", err)
		}

		// Verify that services are loaded
		if m.Table.Rows() == nil || len(m.Table.Rows()) == 0 {
			t.Error("Expected services to be loaded, but table rows are empty")
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

		// Verify that categories are loaded
		if m.Table.Rows() == nil {
			t.Error("Expected categories to be loaded, but table rows are nil")
			return
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

		// Verify the category is selected
		if m.SelectedCategory == nil {
			t.Error("Failed to select a category")
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

		// Verify that operations are loaded
		if m.Table.Rows() == nil {
			t.Error("Expected operations to be loaded, but table rows are nil")
			return
		}

		// Select the Pipeline Status operation if available
		operationFound := false
		for i, row := range m.Table.Rows() {
			if row[0] == "Pipeline Status" {
				m.Table.SetCursor(i)
				m.SelectedOperation = &model.Operation{
					Name:        "Pipeline Status",
					Description: row[1],
				}
				operationFound = true
				break
			}
		}

		if !operationFound {
			t.Log("Pipeline Status operation not found, using first available operation")
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

		// Verify the operation is selected
		if m.SelectedOperation == nil {
			t.Error("Failed to select an operation")
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
}
