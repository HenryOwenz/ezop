package integration_test

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
)

// TestAWSCodePipelineOperations tests that all three CodePipeline operations
// (Pipeline Status, Start Pipeline, and Pipeline Approvals) are available in the Workflows category.
func TestAWSCodePipelineOperations(t *testing.T) {
	// Initialize the model
	m := model.New()

	// Set up the AWS provider
	registry := providers.NewProviderRegistry()
	registry.Register(CreateMockAWSProvider())
	m.Registry = registry

	// Set up AWS profile and region
	m.SetAwsProfile("default")
	m.SetAwsRegion("us-east-1")

	// Create the provider with the selected profile and region
	provider, err := providers.CreateProvider(m.Registry, "AWS", m.GetAwsProfile(), m.GetAwsRegion())
	if err != nil {
		t.Fatalf("Failed to create AWS provider: %v", err)
	}

	// Set the provider in the model
	m.Provider = provider

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
		foundCodePipeline := false
		for i, row := range m.Table.Rows() {
			if row[0] == "CodePipeline" {
				m.Table.SetCursor(i)
				m.SelectedService = &model.Service{
					Name:        "CodePipeline",
					Description: row[1],
				}
				foundCodePipeline = true
				break
			}
		}

		if !foundCodePipeline {
			t.Fatal("CodePipeline service not found in the table")
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

		// In the mock provider, the category is called "Operations" not "Workflows"
		foundCategory := false
		for i, row := range m.Table.Rows() {
			if row[0] == "Operations" {
				m.Table.SetCursor(i)
				m.SelectedCategory = &model.Category{
					Name:        "Operations",
					Description: row[1],
				}
				foundCategory = true
				break
			}
		}

		if !foundCategory {
			t.Fatal("Operations category not found in the table")
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

		// Verify that all three operations are available
		operations := map[string]bool{
			"Pipeline Status": false,
			"Start Pipeline":  false,
			"Manual Approval": false,
		}

		for _, row := range m.Table.Rows() {
			if _, exists := operations[row[0]]; exists {
				operations[row[0]] = true
			}
		}

		for op, found := range operations {
			if !found {
				t.Errorf("Operation '%s' not found in the table", op)
			}
		}
	})

	// Test Pipeline Status operation
	t.Run("Pipeline Status Operation", func(t *testing.T) {
		// Select the Pipeline Status operation
		for i, row := range m.Table.Rows() {
			if row[0] == "Pipeline Status" {
				m.Table.SetCursor(i)
				m.SelectedOperation = &model.Operation{
					Name:        "Pipeline Status",
					Description: row[1],
				}
				break
			}
		}

		if m.SelectedOperation == nil || m.SelectedOperation.Name != "Pipeline Status" {
			t.Fatal("Failed to select Pipeline Status operation")
		}

		// Set the current view to pipeline status
		m.CurrentView = constants.ViewPipelineStatus

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for pipeline status view: %v", err)
		}

		// We don't need to verify actual pipelines in this test
		// Just ensure the view is updated without errors
	})

	// Test Start Pipeline operation
	t.Run("Start Pipeline Operation", func(t *testing.T) {
		// Go back to operation selection
		m.CurrentView = constants.ViewSelectOperation

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for operation selection view: %v", err)
		}

		// Select the Start Pipeline operation
		for i, row := range m.Table.Rows() {
			if row[0] == "Start Pipeline" {
				m.Table.SetCursor(i)
				m.SelectedOperation = &model.Operation{
					Name:        "Start Pipeline",
					Description: row[1],
				}
				break
			}
		}

		if m.SelectedOperation == nil || m.SelectedOperation.Name != "Start Pipeline" {
			t.Fatal("Failed to select Start Pipeline operation")
		}

		// Set the current view to pipeline status (for pipeline selection)
		m.CurrentView = constants.ViewPipelineStatus

		// Update the model for the view
		err = update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for pipeline status view: %v", err)
		}

		// We don't need to verify actual pipelines in this test
		// Just ensure the view is updated without errors
	})

	// Test Manual Approval operation
	t.Run("Manual Approval Operation", func(t *testing.T) {
		// Go back to operation selection
		m.CurrentView = constants.ViewSelectOperation

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for operation selection view: %v", err)
		}

		// Select the Manual Approval operation
		for i, row := range m.Table.Rows() {
			if row[0] == "Manual Approval" {
				m.Table.SetCursor(i)
				m.SelectedOperation = &model.Operation{
					Name:        "Manual Approval",
					Description: row[1],
				}
				break
			}
		}

		if m.SelectedOperation == nil || m.SelectedOperation.Name != "Manual Approval" {
			t.Fatal("Failed to select Manual Approval operation")
		}

		// Set the current view to approvals
		m.CurrentView = constants.ViewApprovals

		// Update the model for the view
		err = update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for approvals view: %v", err)
		}

		// We don't need to verify actual approvals in this test
		// Just ensure the view is updated without errors
	})
}
