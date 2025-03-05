package aws

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/providers"
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
func (a *codePipelineManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
	cloudApprovals, err := a.operation.GetPendingApprovals(ctx)
	if err != nil {
		return nil, err
	}

	// Convert cloud.ApprovalAction to providers.ApprovalAction
	providerApprovals := make([]providers.ApprovalAction, len(cloudApprovals))
	for i, approval := range cloudApprovals {
		providerApprovals[i] = providers.ApprovalAction{
			PipelineName: approval.PipelineName,
			StageName:    approval.StageName,
			ActionName:   approval.ActionName,
			Token:        approval.Token,
		}
	}

	return providerApprovals, nil
}

// ApproveAction approves or rejects an approval action
func (a *codePipelineManualApprovalOperation) ApproveAction(ctx context.Context, action providers.ApprovalAction, approved bool, comment string) error {
	// Convert providers.ApprovalAction to cloud.ApprovalAction
	cloudAction := cloud.ApprovalAction{
		PipelineName: action.PipelineName,
		StageName:    action.StageName,
		ActionName:   action.ActionName,
		Token:        action.Token,
	}

	return a.operation.ApproveAction(ctx, cloudAction, approved, comment)
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
func (a *pipelineStatusOperation) GetPipelineStatus(ctx context.Context) ([]providers.PipelineStatus, error) {
	cloudStatuses, err := a.operation.GetPipelineStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Convert cloud.PipelineStatus to providers.PipelineStatus
	providerStatuses := make([]providers.PipelineStatus, len(cloudStatuses))
	for i, status := range cloudStatuses {
		providerStatus := providers.PipelineStatus{
			Name:   status.Name,
			Stages: make([]providers.StageStatus, len(status.Stages)),
		}

		for j, stage := range status.Stages {
			providerStatus.Stages[j] = providers.StageStatus{
				Name:        stage.Name,
				Status:      stage.Status,
				LastUpdated: stage.LastUpdated,
			}
		}

		providerStatuses[i] = providerStatus
	}

	return providerStatuses, nil
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
func (a *functionStatusOperation) GetFunctionStatus(ctx context.Context) ([]providers.FunctionStatus, error) {
	cloudFunctions, err := a.operation.GetFunctionStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Convert cloud.FunctionStatus to providers.FunctionStatus
	providerFunctions := make([]providers.FunctionStatus, len(cloudFunctions))
	for i, function := range cloudFunctions {
		providerFunctions[i] = providers.FunctionStatus{
			Name:         function.Name,
			Runtime:      function.Runtime,
			Memory:       function.Memory,
			Timeout:      function.Timeout,
			LastUpdate:   function.LastUpdate,
			Role:         function.Role,
			Handler:      function.Handler,
			Description:  function.Description,
			FunctionArn:  function.FunctionArn,
			CodeSize:     function.CodeSize,
			Version:      function.Version,
			PackageType:  function.PackageType,
			Architecture: function.Architecture,
			LogGroup:     function.LogGroup,
		}
	}

	return providerFunctions, nil
}
