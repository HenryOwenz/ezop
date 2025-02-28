package codepipeline

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/providers"
)

// WorkflowsCategory represents the Workflows category for CodePipeline.
type WorkflowsCategory struct {
	profile    string
	region     string
	operations []providers.Operation
}

// NewWorkflowsCategory creates a new Workflows category.
func NewWorkflowsCategory(profile, region string) *WorkflowsCategory {
	category := &WorkflowsCategory{
		profile:    profile,
		region:     region,
		operations: make([]providers.Operation, 0),
	}

	// Register operations
	category.operations = append(category.operations, NewPipelineApprovalsOperation(profile, region))
	category.operations = append(category.operations, NewPipelineStatusOperation(profile, region))
	category.operations = append(category.operations, NewStartPipelineOperation(profile, region))
	// Don't register ApprovalOperation here as it's an internal operation

	return category
}

// Name returns the category's name.
func (c *WorkflowsCategory) Name() string {
	return "Workflows"
}

// Description returns the category's description.
func (c *WorkflowsCategory) Description() string {
	return "Pipeline Workflows and Approvals"
}

// Operations returns all available operations for this category.
func (c *WorkflowsCategory) Operations() []providers.Operation {
	return c.operations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (c *WorkflowsCategory) IsUIVisible() bool {
	return true
}

// PipelineApprovalsOperation represents the Pipeline Approvals operation.
type PipelineApprovalsOperation struct {
	profile string
	region  string
}

// NewPipelineApprovalsOperation creates a new Pipeline Approvals operation.
func NewPipelineApprovalsOperation(profile, region string) *PipelineApprovalsOperation {
	return &PipelineApprovalsOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *PipelineApprovalsOperation) Name() string {
	return "Pipeline Approvals"
}

// Description returns the operation's description.
func (o *PipelineApprovalsOperation) Description() string {
	return "Manage Pipeline Approvals"
}

// Execute executes the operation with the given parameters.
func (o *PipelineApprovalsOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// For fetching approvals, we don't need any parameters
	return GetPendingApprovals(ctx, o.profile, o.region)
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *PipelineApprovalsOperation) IsUIVisible() bool {
	return true
}

// PipelineStatusOperation represents the Pipeline Status operation.
type PipelineStatusOperation struct {
	profile string
	region  string
}

// NewPipelineStatusOperation creates a new Pipeline Status operation.
func NewPipelineStatusOperation(profile, region string) *PipelineStatusOperation {
	return &PipelineStatusOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *PipelineStatusOperation) Name() string {
	return "Pipeline Status"
}

// Description returns the operation's description.
func (o *PipelineStatusOperation) Description() string {
	return "View Pipeline Status"
}

// Execute executes the operation with the given parameters.
func (o *PipelineStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// For fetching pipeline status, we don't need any parameters
	return GetPipelineStatus(ctx, o.profile, o.region)
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *PipelineStatusOperation) IsUIVisible() bool {
	return true
}

// StartPipelineOperation represents the Start Pipeline operation.
type StartPipelineOperation struct {
	profile string
	region  string
}

// NewStartPipelineOperation creates a new Start Pipeline operation.
func NewStartPipelineOperation(profile, region string) *StartPipelineOperation {
	return &StartPipelineOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *StartPipelineOperation) Name() string {
	return "Start Pipeline"
}

// Description returns the operation's description.
func (o *StartPipelineOperation) Description() string {
	return "Trigger Pipeline Execution"
}

// Execute executes the operation with the given parameters.
func (o *StartPipelineOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Extract parameters
	pipelineName, ok := params["pipeline_name"].(string)
	if !ok {
		return nil, fmt.Errorf("pipeline_name parameter is required")
	}

	// CommitID is optional
	commitID := ""
	if commitIDParam, ok := params["commit_id"].(string); ok {
		commitID = commitIDParam
	}

	// Start the pipeline execution
	err := StartPipelineExecution(ctx, o.profile, o.region, pipelineName, commitID)
	if err != nil {
		return nil, err
	}

	return "Pipeline execution started successfully", nil
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *StartPipelineOperation) IsUIVisible() bool {
	return true
}

// ApprovalOperation represents the Approval operation.
type ApprovalOperation struct {
	profile string
	region  string
}

// NewApprovalOperation creates a new Approval operation.
func NewApprovalOperation(profile, region string) *ApprovalOperation {
	return &ApprovalOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *ApprovalOperation) Name() string {
	return "Approval"
}

// Description returns the operation's description.
func (o *ApprovalOperation) Description() string {
	return "Approve or Reject Pipeline Actions"
}

// Execute executes the operation with the given parameters.
func (o *ApprovalOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Extract parameters
	pipelineName, ok := params["pipeline_name"].(string)
	if !ok {
		return nil, fmt.Errorf("pipeline_name parameter is required")
	}

	stageName, ok := params["stage_name"].(string)
	if !ok {
		return nil, fmt.Errorf("stage_name parameter is required")
	}

	actionName, ok := params["action_name"].(string)
	if !ok {
		return nil, fmt.Errorf("action_name parameter is required")
	}

	token, ok := params["token"].(string)
	if !ok {
		return nil, fmt.Errorf("token parameter is required")
	}

	approved, ok := params["approved"].(bool)
	if !ok {
		return nil, fmt.Errorf("approved parameter is required")
	}

	comment, ok := params["comment"].(string)
	if !ok {
		comment = ""
	}

	// Create the approval action
	action := ApprovalAction{
		PipelineName: pipelineName,
		StageName:    stageName,
		ActionName:   actionName,
		Token:        token,
	}

	// Put the approval result
	err := PutApprovalResult(ctx, o.profile, o.region, action, approved, comment)
	if err != nil {
		return nil, err
	}

	if approved {
		return "Pipeline approved successfully", nil
	}
	return "Pipeline rejected successfully", nil
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *ApprovalOperation) IsUIVisible() bool {
	return false
}

// InternalOperationsCategory represents a category for internal operations that shouldn't appear in the UI.
type InternalOperationsCategory struct {
	profile    string
	region     string
	operations []providers.Operation
}

// NewInternalOperationsCategory creates a new category for internal operations.
func NewInternalOperationsCategory(profile, region string) *InternalOperationsCategory {
	category := &InternalOperationsCategory{
		profile:    profile,
		region:     region,
		operations: make([]providers.Operation, 0),
	}

	// Register internal operations
	category.operations = append(category.operations, NewApprovalOperation(profile, region))

	return category
}

// Name returns the category's name.
func (c *InternalOperationsCategory) Name() string {
	return "InternalOperations"
}

// Description returns the category's description.
func (c *InternalOperationsCategory) Description() string {
	return "Internal Operations (Not Visible in UI)"
}

// Operations returns all available operations for this category.
func (c *InternalOperationsCategory) Operations() []providers.Operation {
	return c.operations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (c *InternalOperationsCategory) IsUIVisible() bool {
	return false
}
