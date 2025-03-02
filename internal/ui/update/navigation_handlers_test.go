package update

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
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
		{
			name: "From ViewExecutingAction to ViewPipelineStatus for Start Pipeline operation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewExecutingAction
				m.SelectedOperation = &model.Operation{Name: "Start Pipeline"}
				m.SelectedPipeline = &providers.PipelineStatus{Name: "TestPipeline"}
				return m
			},
			expectedView: constants.ViewPipelineStatus,
			expectedChecks: func(t *testing.T, m *model.Model) {
				// Pipeline should be preserved
				if m.SelectedPipeline == nil || m.SelectedPipeline.Name != "TestPipeline" {
					t.Errorf("Expected SelectedPipeline to be preserved")
				}
				// Approval-related state should be reset
				if m.SelectedApproval != nil {
					t.Errorf("Expected SelectedApproval to be nil")
				}
				if m.ApproveAction {
					t.Errorf("Expected ApproveAction to be false")
				}
				if m.ApprovalComment != "" {
					t.Errorf("Expected ApprovalComment to be empty, got '%s'", m.ApprovalComment)
				}
			},
		},
		{
			name: "From ViewSummary to ViewPipelineStatus for Start Pipeline operation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSummary
				m.SelectedOperation = &model.Operation{Name: "Start Pipeline"}
				m.SelectedPipeline = &providers.PipelineStatus{Name: "TestPipeline"}
				m.Summary = "Test summary"
				return m
			},
			expectedView: constants.ViewPipelineStatus,
			expectedChecks: func(t *testing.T, m *model.Model) {
				// Pipeline should be preserved
				if m.SelectedPipeline == nil || m.SelectedPipeline.Name != "TestPipeline" {
					t.Errorf("Expected SelectedPipeline to be preserved")
				}
				// Summary should be reset
				if m.Summary != "" {
					t.Errorf("Expected Summary to be empty, got '%s'", m.Summary)
				}
			},
		},
		{
			name: "From ViewExecutingAction to ViewSummary for Approval operation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewExecutingAction
				m.SelectedOperation = &model.Operation{Name: "Pipeline Approvals"}
				m.SelectedApproval = &providers.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
				}
				m.ApproveAction = true
				m.Summary = "Test summary"
				return m
			},
			expectedView: constants.ViewSummary,
			expectedChecks: func(t *testing.T, m *model.Model) {
				// Approval should be preserved
				if m.SelectedApproval == nil {
					t.Errorf("Expected SelectedApproval to be preserved")
				}
				// Summary should be preserved
				if m.Summary != "Test summary" {
					t.Errorf("Expected Summary to be preserved, got '%s'", m.Summary)
				}
				// Text input should be focused
				if !m.TextInput.Focused() {
					t.Errorf("Expected TextInput to be focused")
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

// TestPipelineStartNavigationFlow tests the complete navigation flow for the Start Pipeline operation
func TestPipelineStartNavigationFlow(t *testing.T) {
	// Create a model with the initial state
	m := model.New()
	m.CurrentView = constants.ViewSelectOperation
	m.SelectedService = &model.Service{Name: "CodePipeline"}
	m.SelectedCategory = &model.Category{Name: "Operations"}
	m.SelectedOperation = &model.Operation{Name: "Start Pipeline"}

	// Step 1: Navigate to pipeline status view
	// This would normally be done by HandlePipelineStatus
	m.CurrentView = constants.ViewPipelineStatus
	m.SelectedPipeline = &providers.PipelineStatus{Name: "TestPipeline"}

	// Step 2: Navigate to summary view
	// This would normally be done by HandlePipelineSelection
	m.CurrentView = constants.ViewSummary

	// Step 3: Navigate to executing action view
	// This would normally be done by HandleSummaryConfirmation
	m.CurrentView = constants.ViewExecutingAction

	// Step 4: Navigate back from executing action view
	result := NavigateBack(m)

	// Check that we're back at the pipeline status view
	if result.CurrentView != constants.ViewPipelineStatus {
		t.Errorf("Expected to navigate back to ViewPipelineStatus, got %v", result.CurrentView)
	}

	// Check that the pipeline is still selected
	if result.SelectedPipeline == nil || result.SelectedPipeline.Name != "TestPipeline" {
		t.Errorf("Expected SelectedPipeline to be preserved")
	}

	// Check that the operation is still selected
	if result.SelectedOperation == nil || result.SelectedOperation.Name != "Start Pipeline" {
		t.Errorf("Expected SelectedOperation to be preserved")
	}

	// Step 5: Navigate back from pipeline status view
	result = NavigateBack(result)

	// Check that we're back at the select operation view
	if result.CurrentView != constants.ViewSelectOperation {
		t.Errorf("Expected to navigate back to ViewSelectOperation, got %v", result.CurrentView)
	}
}

// TestHandleTextInputSubmissionDoesNotExecuteActions tests that the HandleTextInputSubmission function
// doesn't execute actions immediately after entering text, but instead just transitions to the execution view.
func TestHandleTextInputSubmissionDoesNotExecuteActions(t *testing.T) {
	testCases := []struct {
		name           string
		setupModel     func() *model.Model
		expectedView   constants.View
		expectedFields map[string]string
	}{
		{
			name: "Approval comment submission",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSummary
				m.SelectedOperation = &model.Operation{Name: "Pipeline Approvals"}
				m.SelectedApproval = &providers.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
				}
				m.ApproveAction = true
				m.ManualInput = true
				m.TextInput.SetValue("Test comment")
				m.TextInput.Focus()
				return m
			},
			expectedView: constants.ViewExecutingAction,
			expectedFields: map[string]string{
				"ApprovalComment": "Test comment",
				"Summary":         "Test comment",
			},
		},
		{
			name: "Pipeline start commit ID submission",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSummary
				m.SelectedOperation = &model.Operation{Name: "Start Pipeline"}
				m.SelectedPipeline = &providers.PipelineStatus{Name: "TestPipeline"}
				m.ManualInput = true
				m.TextInput.SetValue("abc123")
				m.TextInput.Focus()
				return m
			},
			expectedView: constants.ViewExecutingAction,
			expectedFields: map[string]string{
				"CommitID": "abc123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the model according to the test case
			m := tc.setupModel()

			// Call HandleTextInputSubmission
			result, cmd := HandleTextInputSubmission(m)

			// Verify that the command is nil (no action execution)
			if cmd != nil {
				t.Errorf("Expected command to be nil (no action execution), got %T", cmd)
			}

			// Unwrap the model
			resultModel, ok := result.(ModelWrapper)
			if !ok {
				t.Fatalf("Expected result to be a ModelWrapper, got %T", result)
			}

			// Verify that the model has transitioned to the expected view
			if resultModel.Model.CurrentView != tc.expectedView {
				t.Errorf("Expected view to be %v, got %v", tc.expectedView, resultModel.Model.CurrentView)
			}

			// Verify that manual input is reset
			if resultModel.Model.ManualInput {
				t.Errorf("Expected ManualInput to be false")
			}

			// Verify that text input is reset
			if resultModel.Model.TextInput.Value() != "" {
				t.Errorf("Expected TextInput value to be empty, got '%s'", resultModel.Model.TextInput.Value())
			}

			// Verify expected fields
			for field, expectedValue := range tc.expectedFields {
				var actualValue string
				switch field {
				case "ApprovalComment":
					actualValue = resultModel.Model.ApprovalComment
				case "Summary":
					actualValue = resultModel.Model.Summary
				case "CommitID":
					actualValue = resultModel.Model.CommitID
				}

				if actualValue != expectedValue {
					t.Errorf("Expected %s to be '%s', got '%s'", field, expectedValue, actualValue)
				}
			}
		})
	}
}
