package integration_test

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
)

// TestAWSServiceVisibility tests the visibility of AWS services, categories, and operations
func TestAWSServiceVisibility(t *testing.T) {
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

		// Verify that CodePipeline service is visible
		foundCodePipeline := false
		for _, row := range m.Table.Rows() {
			if row[0] == "CodePipeline" {
				foundCodePipeline = true
				break
			}
		}

		if !foundCodePipeline {
			t.Error("CodePipeline service not found in the table")
		}
	})

	// Test CodePipeline category visibility
	t.Run("CodePipeline Category Visibility", func(t *testing.T) {
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

		// Set the current view to category selection
		m.CurrentView = constants.ViewSelectCategory

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for category selection view: %v", err)
		}

		// Verify that Operations category is visible
		foundOperations := false
		for _, row := range m.Table.Rows() {
			if row[0] == "Operations" {
				foundOperations = true
				break
			}
		}

		if !foundOperations {
			t.Error("Operations category not found in the table")
		}
	})

	// Test CodePipeline operations visibility
	t.Run("CodePipeline Operations Visibility", func(t *testing.T) {
		// Select the Operations category
		for i, row := range m.Table.Rows() {
			if row[0] == "Operations" {
				m.Table.SetCursor(i)
				m.SelectedCategory = &model.Category{
					Name:        "Operations",
					Description: row[1],
				}
				break
			}
		}

		// Set the current view to operation selection
		m.CurrentView = constants.ViewSelectOperation

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for operation selection view: %v", err)
		}

		// Verify that all operations are visible
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
				t.Errorf("%s operation not found in the table", op)
			}
		}
	})
}
