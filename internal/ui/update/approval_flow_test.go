package update

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	"github.com/charmbracelet/bubbles/table"
)

// TestCompletePipelineApprovalFlow tests the complete flow for pipeline approvals,
// from selecting an approval to entering a comment and then executing the action.
// This test will fail if someone changes the expected behavior of the flow.
func TestCompletePipelineApprovalFlow(t *testing.T) {
	// Step 1: Create initial model with operation selected
	m := model.New()
	m.CurrentView = constants.ViewSelectOperation
	m.SelectedService = &model.Service{Name: "CodePipeline"}
	m.SelectedCategory = &model.Category{Name: "Operations"}
	m.SelectedOperation = &model.Operation{Name: "Pipeline Approvals"}

	// Step 2: Set up approvals view (simulating HandlePipelineApprovals)
	m.CurrentView = constants.ViewApprovals
	m.Approvals = []cloud.ApprovalAction{
		{
			PipelineName: "TestPipeline",
			StageName:    "TestStage",
			ActionName:   "TestAction",
			Token:        "TestToken",
		},
	}
	view.UpdateTableForView(m)

	// Step 3: Select an approval (simulating SelectApproval)
	// This would normally be done by selecting a row in the table
	m.SelectedApproval = &m.Approvals[0]
	result, cmd := SelectApproval(m)

	// Verify we get a model wrapper and no command
	wrapper, ok := result.(ModelWrapper)
	if !ok {
		t.Fatalf("Expected SelectApproval to return a ModelWrapper, got %T", result)
	}
	if cmd != nil {
		t.Errorf("Expected SelectApproval to return nil command, got %T", cmd)
	}

	// Verify we're now at the confirmation view
	if wrapper.Model.CurrentView != constants.ViewConfirmation {
		t.Errorf("Expected to be at ViewConfirmation, got %v", wrapper.Model.CurrentView)
	}

	// Step 4: Choose to approve (simulating HandleConfirmationSelection)
	result, cmd = HandleConfirmationSelection(wrapper.Model)

	// Verify we get a model wrapper and no command
	wrapper, ok = result.(ModelWrapper)
	if !ok {
		t.Fatalf("Expected HandleConfirmationSelection to return a ModelWrapper, got %T", result)
	}
	if cmd != nil {
		t.Errorf("Expected HandleConfirmationSelection to return nil command, got %T", cmd)
	}

	// Verify we're now at the summary view with manual input enabled
	if wrapper.Model.CurrentView != constants.ViewSummary {
		t.Errorf("Expected to be at ViewSummary, got %v", wrapper.Model.CurrentView)
	}
	if !wrapper.Model.ManualInput {
		t.Errorf("Expected ManualInput to be true")
	}
	if !wrapper.Model.ApproveAction {
		t.Errorf("Expected ApproveAction to be true")
	}

	// Step 5: Enter a comment (simulating text input)
	wrapper.Model.TextInput.SetValue("Test approval comment")

	// Step 6: Submit the comment (simulating HandleTextInputSubmission)
	result, cmd = HandleTextInputSubmission(wrapper.Model)

	// Verify we get a model wrapper and no command
	wrapper, ok = result.(ModelWrapper)
	if !ok {
		t.Fatalf("Expected HandleTextInputSubmission to return a ModelWrapper, got %T", result)
	}
	if cmd != nil {
		t.Errorf("Expected HandleTextInputSubmission to return nil command, got %T", cmd)
	}

	// Verify we're now at the executing action view
	if wrapper.Model.CurrentView != constants.ViewExecutingAction {
		t.Errorf("Expected to be at ViewExecutingAction, got %v", wrapper.Model.CurrentView)
	}

	// Verify the comment was stored
	if wrapper.Model.ApprovalComment != "Test approval comment" {
		t.Errorf("Expected ApprovalComment to be 'Test approval comment', got '%s'", wrapper.Model.ApprovalComment)
	}

	// Verify manual input is reset
	if wrapper.Model.ManualInput {
		t.Errorf("Expected ManualInput to be false")
	}

	// Step 7: Execute the approval (simulating HandleExecutionSelection)
	// Set up the table with execution options
	columns := []table.Column{
		{Title: "Action", Width: 10},
		{Title: "Description", Width: 30},
	}
	rows := []table.Row{
		{"Execute", "Execute approval action"},
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

	// Step 8: Verify that navigating back from executing action view goes to summary view
	backResult := NavigateBack(wrapper.Model)

	// Verify we're back at the summary view
	if backResult.CurrentView != constants.ViewSummary {
		t.Errorf("Expected to navigate back to ViewSummary, got %v", backResult.CurrentView)
	}

	// Verify the text input is focused and has the previous comment
	if !backResult.TextInput.Focused() {
		t.Errorf("Expected TextInput to be focused")
	}
	if backResult.TextInput.Value() != wrapper.Model.Summary {
		t.Errorf("Expected TextInput value to be '%s', got '%s'", wrapper.Model.Summary, backResult.TextInput.Value())
	}

	// Step 5: Handle the approval result
	// This would normally be done by HandleApprovalResult
	HandleApprovalResult(m, nil)
}

// TestHandleTextInputSubmissionTransitionsToExecutionView tests that the HandleTextInputSubmission function
// does not execute any actions immediately after entering text, but instead transitions to the execution view.
// This test will fail if someone changes the behavior to execute actions immediately.
func TestHandleTextInputSubmissionTransitionsToExecutionView(t *testing.T) {
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
				m.SelectedApproval = &cloud.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
					Token:        "TestToken",
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
				m.SelectedPipeline = &cloud.PipelineStatus{Name: "TestPipeline"}
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

			// Verify that the model has a table with execution options
			if len(resultModel.Model.Table.Rows()) == 0 {
				t.Errorf("Expected table to have execution options")
			}

			// Verify that one of the options is "Execute"
			hasExecuteOption := false
			for _, row := range resultModel.Model.Table.Rows() {
				if row[0] == "Execute" {
					hasExecuteOption = true
					break
				}
			}
			if !hasExecuteOption {
				t.Errorf("Expected table to have an 'Execute' option")
			}
		})
	}
}
