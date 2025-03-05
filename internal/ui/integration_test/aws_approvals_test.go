package integration

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloudproviders"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
)

// TestAWSApprovalsFlow tests the AWS approvals flow with default/us-east-1
func TestAWSApprovalsFlow(t *testing.T) {
	// Initialize the model
	m := model.New()

	// Set up the AWS provider
	registry := cloud.NewProviderRegistry()
	registry.Register(CreateMockAWSProvider())
	m.Registry = registry

	// Set up AWS profile and region
	t.Run("Setup AWS Configuration", func(t *testing.T) {
		// Set the AWS profile and region
		m.SetAwsProfile("default")
		m.SetAwsRegion("us-east-1")

		// Create the provider with the selected profile and region
		provider, err := cloudproviders.CreateProvider(m.Registry, "AWS", m.GetAwsProfile(), m.GetAwsRegion())
		if err != nil {
			t.Fatalf("Failed to create AWS provider: %v", err)
		}

		// Set the provider in the model
		m.Provider = provider

		// Verify the configuration is set correctly
		if m.GetAwsProfile() != "default" || m.GetAwsRegion() != "us-east-1" {
			t.Errorf("AWS configuration not set correctly. Profile: %s, Region: %s",
				m.GetAwsProfile(), m.GetAwsRegion())
		}
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

		// Select the Manual Approval operation if available
		operationFound := false
		for i, row := range m.Table.Rows() {
			if row[0] == "Manual Approval" {
				m.Table.SetCursor(i)
				m.SelectedOperation = &model.Operation{
					Name:        "Manual Approval",
					Description: row[1],
				}
				operationFound = true
				break
			}
		}

		if !operationFound {
			t.Log("Manual Approval operation not found, using first available operation")
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

	// Test approvals view
	t.Run("Approvals View", func(t *testing.T) {
		// Set the current view to approvals
		m.CurrentView = constants.ViewApprovals

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for approvals view: %v", err)
		}

		// Verify that approvals are loaded (or at least the table is initialized)
		if m.Table.Rows() == nil {
			t.Log("No approvals found or table not initialized")
		} else {
			t.Logf("Found %d approvals", len(m.Table.Rows()))

			// If approvals are found, select the first one
			if len(m.Table.Rows()) > 0 {
				m.Table.SetCursor(0)
				row := m.Table.SelectedRow()

				// Find the approval in the model's approvals
				for _, approval := range m.Approvals {
					if approval.PipelineName == row[0] &&
						approval.StageName == row[1] &&
						approval.ActionName == row[2] {
						m.SelectedApproval = &approval
						break
					}
				}

				// Verify the approval is selected
				if m.SelectedApproval == nil {
					t.Error("Failed to select an approval")
				} else {
					t.Logf("Selected approval: Pipeline: %s, Stage: %s, Action: %s",
						m.SelectedApproval.PipelineName,
						m.SelectedApproval.StageName,
						m.SelectedApproval.ActionName)
				}
			}
		}
	})

	// Test confirmation view
	t.Run("Confirmation View", func(t *testing.T) {
		// Skip if no approval was selected
		if m.SelectedApproval == nil {
			t.Skip("No approval selected, skipping confirmation view test")
		}

		// Set the current view to confirmation
		m.CurrentView = constants.ViewConfirmation

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for confirmation view: %v", err)
		}

		// Verify that confirmation options are loaded
		if m.Table.Rows() == nil || len(m.Table.Rows()) == 0 {
			t.Error("Expected confirmation options to be loaded, but table rows are empty")
			return
		}

		// Select "Approve"
		for i, row := range m.Table.Rows() {
			if row[0] == "Approve" {
				m.Table.SetCursor(i)
				m.ApproveAction = true
				break
			}
		}

		// Verify the approve action is set
		if !m.ApproveAction {
			t.Error("Failed to set approve action")
		}
	})

	// Test summary view
	t.Run("Summary View", func(t *testing.T) {
		// Skip if no approval was selected
		if m.SelectedApproval == nil {
			t.Skip("No approval selected, skipping summary view test")
		}

		// Set the current view to summary
		m.CurrentView = constants.ViewSummary

		// Update the model for the view
		err := update.UpdateModelForView(m)
		if err != nil {
			t.Fatalf("Failed to update model for summary view: %v", err)
		}

		// Set approval comment
		m.SetApprovalComment("Test approval comment")

		// Verify the comment is set
		if m.GetApprovalComment() != "Test approval comment" {
			t.Errorf("Expected approval comment to be 'Test approval comment', got '%s'",
				m.GetApprovalComment())
		}
	})

	// Test executing action view
	t.Run("Executing Action View", func(t *testing.T) {
		// Skip if no approval was selected
		if m.SelectedApproval == nil {
			t.Skip("No approval selected, skipping executing action view test")
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
		t.Log("Action execution would happen here in a real scenario")
	})
}
