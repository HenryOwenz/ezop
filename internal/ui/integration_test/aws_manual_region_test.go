package integration

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
)

// TestAWSManualRegionEntry verifies that manually entering a region works correctly,
// especially when using the default profile.
func TestAWSManualRegionEntry(t *testing.T) {
	// Create a mock AWS provider
	provider := CreateMockAWSProvider()

	// Create a new model
	m := model.New()

	// Initialize the provider registry
	m.Registry = update.InitializeTestRegistry(provider)

	// Set the current view to AWS config
	m.CurrentView = constants.ViewAWSConfig

	// Test steps
	t.Run("Select default profile", func(t *testing.T) {
		// Set the default profile
		m.SetAwsProfile("default")

		// Verify the profile was set
		if m.GetAwsProfile() != "default" {
			t.Errorf("Expected profile to be 'default', got '%s'", m.GetAwsProfile())
		}
	})

	t.Run("Manually enter region", func(t *testing.T) {
		// Enable manual input
		m.ManualInput = true

		// Set up text input for region
		m.TextInput.SetValue("us-east-1")

		// Handle text input submission
		result, _ := update.HandleTextInputSubmission(m)
		updatedModel := result.(update.ModelWrapper).Model

		// Verify the region was set
		if updatedModel.GetAwsRegion() != "us-east-1" {
			t.Errorf("Expected region to be 'us-east-1', got '%s'", updatedModel.GetAwsRegion())
		}

		// Verify manual input was disabled
		if updatedModel.ManualInput {
			t.Error("Expected manual input to be disabled")
		}

		// Verify the view was changed to service selection
		if updatedModel.CurrentView != constants.ViewSelectService {
			t.Errorf("Expected view to be ViewSelectService, got %v", updatedModel.CurrentView)
		}

		// Verify the provider was set
		if updatedModel.Provider == nil {
			t.Error("Expected provider to be set")
		}

		// Verify services are available
		if len(updatedModel.Provider.Services()) == 0 {
			t.Error("Expected services to be available")
		}
	})
}
