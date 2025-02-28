package aws

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/providers"
)

// codePipelineManualApprovalOperation implements the CodePipelineManualApprovalOperation interface
type codePipelineManualApprovalOperation struct {
	provider *Provider
}

// Name returns the operation's name
func (o *codePipelineManualApprovalOperation) Name() string {
	return "Pipeline Approvals"
}

// Description returns the operation's description
func (o *codePipelineManualApprovalOperation) Description() string {
	return "Manage Pipeline Approvals"
}

// IsUIVisible returns whether this operation should be visible in the UI
func (o *codePipelineManualApprovalOperation) IsUIVisible() bool {
	return true
}

// GetPendingApprovals returns all pending manual approval actions
func (o *codePipelineManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
	return o.provider.GetApprovals(ctx)
}

// ApproveAction approves or rejects an approval action
func (o *codePipelineManualApprovalOperation) ApproveAction(ctx context.Context, action providers.ApprovalAction, approved bool, comment string) error {
	return o.provider.ApproveAction(ctx, action, approved, comment)
}

// pipelineStatusOperation implements the PipelineStatusOperation interface
type pipelineStatusOperation struct {
	provider *Provider
}

// Name returns the operation's name
func (o *pipelineStatusOperation) Name() string {
	return "Pipeline Status"
}

// Description returns the operation's description
func (o *pipelineStatusOperation) Description() string {
	return "View Pipeline Status"
}

// IsUIVisible returns whether this operation should be visible in the UI
func (o *pipelineStatusOperation) IsUIVisible() bool {
	return true
}

// GetPipelineStatus returns the status of all pipelines
func (o *pipelineStatusOperation) GetPipelineStatus(ctx context.Context) ([]providers.PipelineStatus, error) {
	return o.provider.GetStatus(ctx)
}

// startPipelineOperation implements the StartPipelineOperation interface
type startPipelineOperation struct {
	provider *Provider
}

// Name returns the operation's name
func (o *startPipelineOperation) Name() string {
	return "Start Pipeline"
}

// Description returns the operation's description
func (o *startPipelineOperation) Description() string {
	return "Trigger Pipeline Execution"
}

// IsUIVisible returns whether this operation should be visible in the UI
func (o *startPipelineOperation) IsUIVisible() bool {
	return true
}

// StartPipelineExecution starts a pipeline execution
func (o *startPipelineOperation) StartPipelineExecution(ctx context.Context, pipelineName string, commitID string) error {
	return o.provider.StartPipeline(ctx, pipelineName, commitID)
}
