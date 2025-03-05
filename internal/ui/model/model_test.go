package model

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

func TestGetterSetterMethods(t *testing.T) {
	// Create a new model
	m := New()

	// Test AWS profile getter/setter
	t.Run("AwsProfile", func(t *testing.T) {
		// Test empty initial value
		if m.GetAwsProfile() != "" {
			t.Errorf("Expected empty AWS profile, got '%s'", m.GetAwsProfile())
		}

		// Test setting and getting a value
		m.SetAwsProfile("test-profile")
		if m.GetAwsProfile() != "test-profile" {
			t.Errorf("Expected AWS profile 'test-profile', got '%s'", m.GetAwsProfile())
		}
	})

	// Test AWS region getter/setter
	t.Run("AwsRegion", func(t *testing.T) {
		// Test empty initial value
		if m.GetAwsRegion() != "" {
			t.Errorf("Expected empty AWS region, got '%s'", m.GetAwsRegion())
		}

		// Test setting and getting a value
		m.SetAwsRegion("us-west-2")
		if m.GetAwsRegion() != "us-west-2" {
			t.Errorf("Expected AWS region 'us-west-2', got '%s'", m.GetAwsRegion())
		}
	})

	// Test approval comment getter/setter
	t.Run("ApprovalComment", func(t *testing.T) {
		// Test empty initial value
		if m.GetApprovalComment() != "" {
			t.Errorf("Expected empty approval comment, got '%s'", m.GetApprovalComment())
		}

		// Test setting and getting a value
		m.SetApprovalComment("Test comment")
		if m.GetApprovalComment() != "Test comment" {
			t.Errorf("Expected approval comment 'Test comment', got '%s'", m.GetApprovalComment())
		}
	})

	// Test approve action getter/setter
	t.Run("ApproveAction", func(t *testing.T) {
		// Test default value
		if m.GetApproveAction() != false {
			t.Errorf("Expected default approve action to be false, got %v", m.GetApproveAction())
		}

		// Test setting and getting a value
		m.SetApproveAction(true)
		if !m.GetApproveAction() {
			t.Errorf("Expected approve action to be true, got false")
		}
	})

	// Test selected approval getter/setter
	t.Run("SelectedApproval", func(t *testing.T) {
		// Test nil initial value
		if m.GetSelectedApproval() != nil {
			t.Errorf("Expected nil selected approval, got non-nil")
		}

		// Test setting and getting a value
		approval := &cloud.ApprovalAction{
			PipelineName: "TestPipeline",
			StageName:    "TestStage",
			ActionName:   "TestAction",
			Token:        "TestToken",
		}
		m.SetSelectedApproval(approval)

		result := m.GetSelectedApproval()
		if result == nil {
			t.Errorf("Expected non-nil selected approval, got nil")
		} else if result.PipelineName != "TestPipeline" {
			t.Errorf("Expected pipeline name 'TestPipeline', got '%s'", result.PipelineName)
		}
	})

	// Test selected pipeline getter/setter
	t.Run("SelectedPipeline", func(t *testing.T) {
		// Test nil initial value
		if m.GetSelectedPipeline() != nil {
			t.Errorf("Expected nil selected pipeline, got non-nil")
		}

		// Test setting and getting a value
		pipeline := &cloud.PipelineStatus{
			Name: "TestPipeline",
			Stages: []cloud.StageStatus{
				{
					Name:   "TestStage",
					Status: "Succeeded",
				},
			},
		}
		m.SetSelectedPipeline(pipeline)

		result := m.GetSelectedPipeline()
		if result == nil {
			t.Errorf("Expected non-nil selected pipeline, got nil")
		} else if result.Name != "TestPipeline" {
			t.Errorf("Expected pipeline name 'TestPipeline', got '%s'", result.Name)
		}
	})

	// Test approvals getter/setter
	t.Run("Approvals", func(t *testing.T) {
		// Test empty initial value
		if len(m.GetApprovals()) != 0 {
			t.Errorf("Expected empty approvals, got %d items", len(m.GetApprovals()))
		}

		// Test setting and getting a value
		approvals := []cloud.ApprovalAction{
			{
				PipelineName: "TestPipeline1",
				StageName:    "TestStage1",
				ActionName:   "TestAction1",
			},
			{
				PipelineName: "TestPipeline2",
				StageName:    "TestStage2",
				ActionName:   "TestAction2",
			},
		}
		m.SetApprovals(approvals)

		result := m.GetApprovals()
		if len(result) != 2 {
			t.Errorf("Expected 2 approvals, got %d", len(result))
		} else if result[0].PipelineName != "TestPipeline1" {
			t.Errorf("Expected first pipeline name 'TestPipeline1', got '%s'", result[0].PipelineName)
		}
	})

	// Test pipelines getter/setter
	t.Run("Pipelines", func(t *testing.T) {
		// Test empty initial value
		if len(m.GetPipelines()) != 0 {
			t.Errorf("Expected empty pipelines, got %d items", len(m.GetPipelines()))
		}

		// Test setting and getting a value
		pipelines := []cloud.PipelineStatus{
			{
				Name: "TestPipeline1",
			},
			{
				Name: "TestPipeline2",
			},
		}
		m.SetPipelines(pipelines)

		result := m.GetPipelines()
		if len(result) != 2 {
			t.Errorf("Expected 2 pipelines, got %d", len(result))
		} else if result[0].Name != "TestPipeline1" {
			t.Errorf("Expected first pipeline name 'TestPipeline1', got '%s'", result[0].Name)
		}
	})
}

func TestModelClone(t *testing.T) {
	// Create a new model
	m := New()

	// Set up the original model with various values
	m.SetAwsProfile("original-profile")
	m.SetAwsRegion("us-east-1")
	m.SetApprovalComment("Original comment")
	m.SetApproveAction(true)

	// Set up provider specific state
	if m.ProviderState.ProviderSpecificState == nil {
		m.ProviderState.ProviderSpecificState = make(map[string]interface{})
	}
	m.ProviderState.ProviderSpecificState["testKey"] = "testValue"

	// Set up input state
	if m.InputState.TextValues == nil {
		m.InputState.TextValues = make(map[string]string)
	}
	m.InputState.TextValues["testInput"] = "testValue"

	// Clone the model
	clone := m.Clone()

	// Test that basic values are cloned correctly
	t.Run("BasicValues", func(t *testing.T) {
		if clone.GetAwsProfile() != "original-profile" {
			t.Errorf("Expected cloned profile 'original-profile', got '%s'", clone.GetAwsProfile())
		}

		if clone.GetAwsRegion() != "us-east-1" {
			t.Errorf("Expected cloned region 'us-east-1', got '%s'", clone.GetAwsRegion())
		}

		if clone.GetApprovalComment() != "Original comment" {
			t.Errorf("Expected cloned comment 'Original comment', got '%s'", clone.GetApprovalComment())
		}

		if !clone.GetApproveAction() {
			t.Errorf("Expected cloned approve action to be true, got false")
		}
	})

	// Test that modifying the clone doesn't affect the original
	t.Run("Independence", func(t *testing.T) {
		// Modify the clone
		clone.SetAwsProfile("modified-profile")
		clone.SetAwsRegion("eu-west-1")
		clone.SetApprovalComment("Modified comment")
		clone.SetApproveAction(false)

		// Verify the original is unchanged
		if m.GetAwsProfile() != "original-profile" {
			t.Errorf("Original profile was modified when clone was changed")
		}

		if m.GetAwsRegion() != "us-east-1" {
			t.Errorf("Original region was modified when clone was changed")
		}

		if m.GetApprovalComment() != "Original comment" {
			t.Errorf("Original comment was modified when clone was changed")
		}

		if !m.GetApproveAction() {
			t.Errorf("Original approve action was modified when clone was changed")
		}
	})

	// Test that maps are deep copied
	t.Run("DeepCopyMaps", func(t *testing.T) {
		// Check that provider specific state was copied
		if clone.ProviderState.ProviderSpecificState == nil {
			t.Errorf("Provider specific state map was not cloned")
		} else if val, ok := clone.ProviderState.ProviderSpecificState["testKey"]; !ok || val != "testValue" {
			t.Errorf("Provider specific state values were not cloned correctly")
		}

		// Check that input state was copied
		if clone.InputState.TextValues == nil {
			t.Errorf("Input text values map was not cloned")
		} else if val, ok := clone.InputState.TextValues["testInput"]; !ok || val != "testValue" {
			t.Errorf("Input text values were not cloned correctly")
		}

		// Modify the clone's maps
		clone.ProviderState.ProviderSpecificState["testKey"] = "modifiedValue"
		clone.InputState.TextValues["testInput"] = "modifiedValue"

		// Verify the original maps are unchanged
		if val, ok := m.ProviderState.ProviderSpecificState["testKey"]; !ok || val != "testValue" {
			t.Errorf("Original provider specific state was modified when clone was changed")
		}

		if val, ok := m.InputState.TextValues["testInput"]; !ok || val != "testValue" {
			t.Errorf("Original input text values were modified when clone was changed")
		}
	})
}
