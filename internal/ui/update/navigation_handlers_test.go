package update

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

func TestNavigateBack(t *testing.T) {
	// Test navigation from different views
	testCases := []struct {
		name           string
		setupModel     func() *model.Model
		expectedView   constants.View
		expectedChecks func(t *testing.T, m *model.Model)
	}{
		{
			name: "From ViewSelectService to ViewAWSConfig",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectService
				m.SetAwsProfile("test-profile")
				m.SetAwsRegion("us-west-2")
				m.SelectedService = &model.Service{Name: "TestService"}
				return m
			},
			expectedView: constants.ViewAWSConfig,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.SelectedService != nil {
					t.Errorf("Expected SelectedService to be nil after navigating back")
				}
				// Profile and region should be preserved
				if m.GetAwsProfile() != "test-profile" {
					t.Errorf("Expected AWS profile to be preserved, got '%s'", m.GetAwsProfile())
				}
				if m.GetAwsRegion() != "us-west-2" {
					t.Errorf("Expected AWS region to be preserved, got '%s'", m.GetAwsRegion())
				}
			},
		},
		{
			name: "From ViewSelectCategory to ViewSelectService",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectCategory
				m.SelectedService = &model.Service{Name: "TestService"}
				m.SelectedCategory = &model.Category{Name: "TestCategory"}
				return m
			},
			expectedView: constants.ViewSelectService,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.SelectedCategory != nil {
					t.Errorf("Expected SelectedCategory to be nil after navigating back")
				}
				// Service should be preserved
				if m.SelectedService == nil || m.SelectedService.Name != "TestService" {
					t.Errorf("Expected SelectedService to be preserved")
				}
			},
		},
		{
			name: "From ViewSelectOperation to ViewSelectCategory",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectOperation
				m.SelectedService = &model.Service{Name: "TestService"}
				m.SelectedCategory = &model.Category{Name: "TestCategory"}
				m.SelectedOperation = &model.Operation{Name: "TestOperation"}
				return m
			},
			expectedView: constants.ViewSelectCategory,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.SelectedOperation != nil {
					t.Errorf("Expected SelectedOperation to be nil after navigating back")
				}
				// Service and category should be preserved
				if m.SelectedService == nil || m.SelectedService.Name != "TestService" {
					t.Errorf("Expected SelectedService to be preserved")
				}
				if m.SelectedCategory == nil || m.SelectedCategory.Name != "TestCategory" {
					t.Errorf("Expected SelectedCategory to be preserved")
				}
			},
		},
		{
			name: "From ViewAWSConfig with profile to ViewProviders",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewAWSConfig
				m.SetAwsProfile("")
				return m
			},
			expectedView: constants.ViewProviders,
			expectedChecks: func(t *testing.T, m *model.Model) {
				// No specific checks needed
			},
		},
		{
			name: "From ViewAWSConfig with profile and region to profile selection",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewAWSConfig
				m.SetAwsProfile("test-profile")
				m.SetAwsRegion("us-west-2")
				return m
			},
			expectedView: constants.ViewAWSConfig,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.GetAwsProfile() != "" {
					t.Errorf("Expected AWS profile to be cleared, got '%s'", m.GetAwsProfile())
				}
				if m.GetAwsRegion() != "" {
					t.Errorf("Expected AWS region to be cleared, got '%s'", m.GetAwsRegion())
				}
			},
		},
		{
			name: "From ViewApprovals to ViewSelectOperation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewApprovals
				m.SelectedOperation = &model.Operation{Name: "Manual Approval"}
				return m
			},
			expectedView: constants.ViewSelectOperation,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.GetSelectedApproval() != nil {
					t.Errorf("Expected SelectedApproval to be nil after navigating back")
				}
			},
		},
		{
			name: "From ViewPipelineStatus to ViewSelectOperation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.SelectedOperation = &model.Operation{Name: "Pipeline Status"}
				return m
			},
			expectedView: constants.ViewSelectOperation,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.GetSelectedPipeline() != nil {
					t.Errorf("Expected SelectedPipeline to be nil after navigating back")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the model according to the test case
			m := tc.setupModel()

			// Call NavigateBack
			result := NavigateBack(m)

			// Check that the result is not nil
			if result == nil {
				t.Fatal("NavigateBack returned nil")
			}

			// Check that the view was changed as expected
			if result.CurrentView != tc.expectedView {
				t.Errorf("Expected view to be %v, got %v", tc.expectedView, result.CurrentView)
			}

			// Run any additional checks specific to the test case
			tc.expectedChecks(t, result)
		})
	}
}
