package update

import (
	"errors"
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

func TestHandleApprovalResult(t *testing.T) {
	testCases := []struct {
		name           string
		setupModel     func() *model.Model
		err            error
		expectedView   constants.View
		expectedChecks func(t *testing.T, m *model.Model)
	}{
		{
			name: "With error",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewExecutingAction
				m.SelectedApproval = &providers.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
				}
				return m
			},
			err:          errors.New("test error"),
			expectedView: constants.ViewError,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.Error == "" {
					t.Errorf("Expected Error to be set")
				}
			},
		},
		{
			name: "With nil SelectedApproval",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewExecutingAction
				m.SelectedApproval = nil
				return m
			},
			err:          nil,
			expectedView: constants.ViewError,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.Error != constants.MsgErrorNoApproval {
					t.Errorf("Expected Error to be set to MsgErrorNoApproval, got '%s'", m.Error)
				}
			},
		},
		{
			name: "With approval action",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewExecutingAction
				m.SelectedApproval = &providers.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
				}
				m.ApproveAction = true
				return m
			},
			err:          nil,
			expectedView: constants.ViewSelectOperation,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.Success == "" {
					t.Errorf("Expected Success to be set")
				}
				if m.SelectedApproval != nil {
					t.Errorf("Expected SelectedApproval to be nil")
				}
				if m.ApprovalComment != "" {
					t.Errorf("Expected ApprovalComment to be empty")
				}
				if m.Approvals != nil {
					t.Errorf("Expected Approvals to be nil")
				}
			},
		},
		{
			name: "With rejection action",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewExecutingAction
				m.SelectedApproval = &providers.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
				}
				m.ApproveAction = false
				return m
			},
			err:          nil,
			expectedView: constants.ViewSelectOperation,
			expectedChecks: func(t *testing.T, m *model.Model) {
				if m.Success == "" {
					t.Errorf("Expected Success to be set")
				}
				if m.SelectedApproval != nil {
					t.Errorf("Expected SelectedApproval to be nil")
				}
				if m.ApprovalComment != "" {
					t.Errorf("Expected ApprovalComment to be empty")
				}
				if m.Approvals != nil {
					t.Errorf("Expected Approvals to be nil")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the model according to the test case
			m := tc.setupModel()

			// Call HandleApprovalResult
			HandleApprovalResult(m, tc.err)

			// Check that the view was changed as expected
			if m.CurrentView != tc.expectedView {
				t.Errorf("Expected view to be %v, got %v", tc.expectedView, m.CurrentView)
			}

			// Run any additional checks specific to the test case
			tc.expectedChecks(t, m)
		})
	}
}

// TestApprovalFlow tests the complete approval flow
func TestApprovalFlow(t *testing.T) {
	// Create a model with the initial state
	m := model.New()
	m.CurrentView = constants.ViewSelectOperation
	m.SelectedService = &model.Service{Name: "CodePipeline"}
	m.SelectedCategory = &model.Category{Name: "Operations"}
	m.SelectedOperation = &model.Operation{Name: "Pipeline Approvals"}

	// Step 1: Navigate to approvals view
	// This would normally be done by HandlePipelineApprovals
	m.CurrentView = constants.ViewApprovals
	m.Approvals = []providers.ApprovalAction{
		{
			PipelineName: "TestPipeline",
			StageName:    "TestStage",
			ActionName:   "TestAction",
		},
	}

	// Step 2: Select an approval
	// This would normally be done by SelectApproval
	m.CurrentView = constants.ViewConfirmation
	m.SelectedApproval = &m.Approvals[0]

	// Step 3: Choose to approve
	// This would normally be done by HandleConfirmationSelection
	m.CurrentView = constants.ViewSummary
	m.ApproveAction = true
	m.ManualInput = true
	m.TextInput.SetValue("Test comment")
	m.TextInput.Focus()

	// Step 4: Submit the comment
	// This would normally be done by HandleSummaryConfirmation
	newModel, _ := HandleSummaryConfirmation(m)
	if wrapper, ok := newModel.(ModelWrapper); ok {
		m = wrapper.Model
	}

	// Check that we're at the executing action view
	if m.CurrentView != constants.ViewExecutingAction {
		t.Errorf("Expected to be at ViewExecutingAction, got %v", m.CurrentView)
	}

	// Check that the approval comment is set
	if m.ApprovalComment != "Test comment" {
		t.Errorf("Expected ApprovalComment to be 'Test comment', got '%s'", m.ApprovalComment)
	}

	// Step 5: Handle the approval result
	// This would normally be done by HandleApprovalResult
	HandleApprovalResult(m, nil)

	// Check that we're back at the select operation view
	if m.CurrentView != constants.ViewSelectOperation {
		t.Errorf("Expected to navigate back to ViewSelectOperation, got %v", m.CurrentView)
	}

	// Check that the approval state is reset
	if m.SelectedApproval != nil {
		t.Errorf("Expected SelectedApproval to be nil")
	}
	if m.ApprovalComment != "" {
		t.Errorf("Expected ApprovalComment to be empty")
	}
	if m.Approvals != nil {
		t.Errorf("Expected Approvals to be nil")
	}
}
