package providers

import (
	"context"
)

// UIOperation represents a user-facing operation in the UI
type UIOperation interface {
	// Name returns the operation's name
	Name() string

	// Description returns the operation's description
	Description() string

	// IsUIVisible returns whether this operation should be visible in the UI
	IsUIVisible() bool
}

// CodePipelineManualApprovalOperation represents a manual approval operation for AWS CodePipeline
type CodePipelineManualApprovalOperation interface {
	UIOperation

	// GetPendingApprovals returns all pending manual approval actions
	GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error)

	// ApproveAction approves or rejects an approval action
	ApproveAction(ctx context.Context, action ApprovalAction, approved bool, comment string) error
}

// PipelineStatusOperation represents an operation to view pipeline status
type PipelineStatusOperation interface {
	UIOperation

	// GetPipelineStatus returns the status of all pipelines
	GetPipelineStatus(ctx context.Context) ([]PipelineStatus, error)
}

// StartPipelineOperation represents an operation to start a pipeline execution
type StartPipelineOperation interface {
	UIOperation

	// StartPipelineExecution starts a pipeline execution
	StartPipelineExecution(ctx context.Context, pipelineName string, commitID string) error
}

// FunctionStatusOperation represents an operation to view Lambda function status
type FunctionStatusOperation interface {
	UIOperation

	// GetFunctionStatus returns the status of all Lambda functions
	GetFunctionStatus(ctx context.Context) ([]FunctionStatus, error)
}

// FunctionStatus represents the status of a Lambda function
type FunctionStatus struct {
	Name         string
	Runtime      string
	Memory       int32
	Timeout      int32
	LastUpdate   string
	Role         string
	Handler      string
	Description  string
	FunctionArn  string
	CodeSize     int64
	Version      string
	PackageType  string
	Architecture string
	LogGroup     string
}
