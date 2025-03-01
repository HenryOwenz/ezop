package view

import (
	"strings"
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

func TestGetContextText(t *testing.T) {
	// Test context text for different views
	testCases := []struct {
		name           string
		setupModel     func() *model.Model
		expectedText   string
		expectedChecks func(t *testing.T, text string)
	}{
		{
			name: "ViewProviders",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewProviders
				return m
			},
			expectedText: constants.MsgAppDescription,
			expectedChecks: func(t *testing.T, text string) {
				// No additional checks needed
			},
		},
		{
			name: "ViewAWSConfig - Profile Selection",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewAWSConfig
				m.SetAwsProfile("")
				return m
			},
			expectedText: "Amazon Web Services",
			expectedChecks: func(t *testing.T, text string) {
				// No additional checks needed
			},
		},
		{
			name: "ViewAWSConfig - Region Selection",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewAWSConfig
				m.SetAwsProfile("test-profile")
				return m
			},
			expectedText: "Profile: test-profile",
			expectedChecks: func(t *testing.T, text string) {
				// No additional checks needed
			},
		},
		{
			name: "ViewSelectService",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectService
				m.SetAwsProfile("test-profile")
				m.SetAwsRegion("us-west-2")
				return m
			},
			expectedText: "",
			expectedChecks: func(t *testing.T, text string) {
				if !strings.Contains(text, "Profile: test-profile") {
					t.Errorf("Expected context text to contain profile, got '%s'", text)
				}
				if !strings.Contains(text, "Region: us-west-2") {
					t.Errorf("Expected context text to contain region, got '%s'", text)
				}
			},
		},
		{
			name: "ViewSelectCategory",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectCategory
				m.SelectedService = &model.Service{Name: "TestService"}
				return m
			},
			expectedText: "Service: TestService",
			expectedChecks: func(t *testing.T, text string) {
				// No additional checks needed
			},
		},
		{
			name: "ViewSelectOperation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectOperation
				m.SelectedService = &model.Service{Name: "TestService"}
				m.SelectedCategory = &model.Category{Name: "TestCategory"}
				return m
			},
			expectedText: "",
			expectedChecks: func(t *testing.T, text string) {
				if !strings.Contains(text, "Service: TestService") {
					t.Errorf("Expected context text to contain service, got '%s'", text)
				}
				if !strings.Contains(text, "Category: TestCategory") {
					t.Errorf("Expected context text to contain category, got '%s'", text)
				}
			},
		},
		{
			name: "ViewApprovals",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewApprovals
				m.SetAwsProfile("test-profile")
				m.SetAwsRegion("us-west-2")
				return m
			},
			expectedText: "",
			expectedChecks: func(t *testing.T, text string) {
				if !strings.Contains(text, "Profile: test-profile") {
					t.Errorf("Expected context text to contain profile, got '%s'", text)
				}
				if !strings.Contains(text, "Region: us-west-2") {
					t.Errorf("Expected context text to contain region, got '%s'", text)
				}
			},
		},
		{
			name: "ViewConfirmation - Pipeline Start",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewConfirmation
				m.SetAwsProfile("test-profile")
				m.SetAwsRegion("us-west-2")
				m.SelectedOperation = &model.Operation{Name: "Start Pipeline"}
				m.SetSelectedPipeline(&providers.PipelineStatus{Name: "TestPipeline"})
				return m
			},
			expectedText: "",
			expectedChecks: func(t *testing.T, text string) {
				if !strings.Contains(text, "Profile: test-profile") {
					t.Errorf("Expected context text to contain profile, got '%s'", text)
				}
				if !strings.Contains(text, "Region: us-west-2") {
					t.Errorf("Expected context text to contain region, got '%s'", text)
				}
				if !strings.Contains(text, "Pipeline: TestPipeline") {
					t.Errorf("Expected context text to contain pipeline, got '%s'", text)
				}
			},
		},
		{
			name: "ViewConfirmation - Approval",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewConfirmation
				m.SetSelectedApproval(&providers.ApprovalAction{
					PipelineName: "TestPipeline",
					StageName:    "TestStage",
					ActionName:   "TestAction",
				})
				return m
			},
			expectedText: "",
			expectedChecks: func(t *testing.T, text string) {
				if !strings.Contains(text, "Pipeline: TestPipeline") {
					t.Errorf("Expected context text to contain pipeline, got '%s'", text)
				}
				if !strings.Contains(text, "Stage: TestStage") {
					t.Errorf("Expected context text to contain stage, got '%s'", text)
				}
				if !strings.Contains(text, "Action: TestAction") {
					t.Errorf("Expected context text to contain action, got '%s'", text)
				}
			},
		},
		{
			name: "ViewPipelineStatus",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.SetAwsProfile("test-profile")
				m.SetAwsRegion("us-west-2")
				return m
			},
			expectedText: "",
			expectedChecks: func(t *testing.T, text string) {
				if !strings.Contains(text, "Profile: test-profile") {
					t.Errorf("Expected context text to contain profile, got '%s'", text)
				}
				if !strings.Contains(text, "Region: us-west-2") {
					t.Errorf("Expected context text to contain region, got '%s'", text)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the model according to the test case
			m := tc.setupModel()

			// Get the context text
			text := getContextText(m)

			// If an exact match is expected, check it
			if tc.expectedText != "" && text != tc.expectedText {
				t.Errorf("Expected context text '%s', got '%s'", tc.expectedText, text)
			}

			// Run any additional checks specific to the test case
			tc.expectedChecks(t, text)
		})
	}
}

func TestGetTitleText(t *testing.T) {
	// Test title text for different views
	testCases := []struct {
		name         string
		setupModel   func() *model.Model
		expectedText string
	}{
		{
			name: "ViewProviders",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewProviders
				return m
			},
			expectedText: constants.TitleProviders,
		},
		{
			name: "ViewAWSConfig - Profile Selection",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewAWSConfig
				m.SetAwsProfile("")
				return m
			},
			expectedText: constants.TitleSelectProfile,
		},
		{
			name: "ViewAWSConfig - Region Selection",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewAWSConfig
				m.SetAwsProfile("test-profile")
				return m
			},
			expectedText: constants.TitleSelectRegion,
		},
		{
			name: "ViewSelectService",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectService
				return m
			},
			expectedText: constants.TitleSelectService,
		},
		{
			name: "ViewSelectCategory",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectCategory
				return m
			},
			expectedText: constants.TitleSelectCategory,
		},
		{
			name: "ViewSelectOperation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewSelectOperation
				return m
			},
			expectedText: constants.TitleSelectOperation,
		},
		{
			name: "ViewApprovals",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewApprovals
				return m
			},
			expectedText: constants.TitleApprovals,
		},
		{
			name: "ViewConfirmation",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewConfirmation
				return m
			},
			expectedText: constants.TitleConfirmation,
		},
		{
			name: "ViewPipelineStatus",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				return m
			},
			expectedText: constants.TitlePipelineStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the model according to the test case
			m := tc.setupModel()

			// Get the title text
			text := getTitleText(m)

			// Check that the title text is as expected
			if text != tc.expectedText {
				t.Errorf("Expected title text '%s', got '%s'", tc.expectedText, text)
			}
		})
	}
}
