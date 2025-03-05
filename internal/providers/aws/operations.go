package aws

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// codePipelineManualApprovalOperation adapts a cloud.CodePipelineManualApprovalOperation to a providers.CodePipelineManualApprovalOperation
type codePipelineManualApprovalOperation struct {
	operation cloud.CodePipelineManualApprovalOperation
}

func newCodePipelineManualApprovalOperation(operation cloud.CodePipelineManualApprovalOperation) *codePipelineManualApprovalOperation {
	return &codePipelineManualApprovalOperation{
		operation: operation,
	}
}

// Name returns the operation's name
func (a *codePipelineManualApprovalOperation) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *codePipelineManualApprovalOperation) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *codePipelineManualApprovalOperation) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// GetPendingApprovals returns all pending manual approval actions
func (a *codePipelineManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	return a.operation.GetPendingApprovals(ctx)
}

// ApproveAction approves or rejects an approval action
func (a *codePipelineManualApprovalOperation) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	return a.operation.ApproveAction(ctx, action, approved, comment)
}

// pipelineStatusOperation adapts a cloud.PipelineStatusOperation to a providers.PipelineStatusOperation
type pipelineStatusOperation struct {
	operation cloud.PipelineStatusOperation
}

func newPipelineStatusOperation(operation cloud.PipelineStatusOperation) *pipelineStatusOperation {
	return &pipelineStatusOperation{
		operation: operation,
	}
}

// Name returns the operation's name
func (a *pipelineStatusOperation) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *pipelineStatusOperation) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *pipelineStatusOperation) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// GetPipelineStatus returns the status of all pipelines
func (a *pipelineStatusOperation) GetPipelineStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	return a.operation.GetPipelineStatus(ctx)
}

// startPipelineOperation adapts a cloud.StartPipelineOperation to a providers.StartPipelineOperation
type startPipelineOperation struct {
	operation cloud.StartPipelineOperation
}

func newStartPipelineOperation(operation cloud.StartPipelineOperation) *startPipelineOperation {
	return &startPipelineOperation{
		operation: operation,
	}
}

// Name returns the operation's name
func (a *startPipelineOperation) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *startPipelineOperation) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *startPipelineOperation) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// StartPipelineExecution starts a pipeline execution
func (a *startPipelineOperation) StartPipelineExecution(ctx context.Context, pipelineName string, commitID string) error {
	return a.operation.StartPipelineExecution(ctx, pipelineName, commitID)
}

// functionStatusOperation adapts a cloud.FunctionStatusOperation to a providers.FunctionStatusOperation
type functionStatusOperation struct {
	operation cloud.FunctionStatusOperation
}

func newFunctionStatusOperation(operation cloud.FunctionStatusOperation) *functionStatusOperation {
	return &functionStatusOperation{
		operation: operation,
	}
}

// Name returns the operation's name
func (a *functionStatusOperation) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *functionStatusOperation) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *functionStatusOperation) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// GetFunctionStatus returns the status of all Lambda functions
func (a *functionStatusOperation) GetFunctionStatus(ctx context.Context) ([]cloud.FunctionStatus, error) {
	return a.operation.GetFunctionStatus(ctx)
}
