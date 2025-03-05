package cloud

import (
	"context"
)

// Provider represents a cloud provider.
type Provider interface {
	// Name returns the provider's name.
	Name() string

	// Description returns the provider's description.
	Description() string

	// Services returns all available services for this provider.
	Services() []Service

	// GetProfiles returns all available profiles for this provider.
	GetProfiles() ([]string, error)

	// LoadConfig loads the provider configuration with the given profile and region.
	LoadConfig(profile, region string) error

	// GetFunctionStatusOperation returns the function status operation
	GetFunctionStatusOperation() (FunctionStatusOperation, error)

	// GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
	GetCodePipelineManualApprovalOperation() (CodePipelineManualApprovalOperation, error)

	// GetPipelineStatusOperation returns the pipeline status operation
	GetPipelineStatusOperation() (PipelineStatusOperation, error)

	// GetStartPipelineOperation returns the start pipeline operation
	GetStartPipelineOperation() (StartPipelineOperation, error)

	// GetAuthenticationMethods returns available authentication methods
	GetAuthenticationMethods() []string

	// GetAuthConfigKeys returns configuration keys for the specified authentication method
	GetAuthConfigKeys(method string) []string

	// Authenticate authenticates using the provided method and configuration
	Authenticate(method string, authConfig map[string]string) error

	// IsAuthenticated checks if the provider is authenticated
	IsAuthenticated() bool

	// GetConfigKeys returns required configuration keys
	GetConfigKeys() []string

	// GetConfigOptions returns available options for a configuration key
	GetConfigOptions(key string) ([]string, error)

	// Configure configures the provider with the given configuration
	Configure(config map[string]string) error

	// GetApprovals returns pending approvals for the provider
	GetApprovals(ctx context.Context) ([]ApprovalAction, error)

	// ApproveAction approves or rejects an approval action
	ApproveAction(ctx context.Context, action ApprovalAction, approved bool, comment string) error

	// GetStatus returns the status of all pipelines
	GetStatus(ctx context.Context) ([]PipelineStatus, error)

	// StartPipeline starts a pipeline execution
	StartPipeline(ctx context.Context, pipelineName string, commitID string) error
}

// Service represents a cloud service.
type Service interface {
	// Name returns the service's name.
	Name() string

	// Description returns the service's description.
	Description() string

	// Categories returns all available categories for this service.
	Categories() []Category
}

// Category represents a group of operations.
type Category interface {
	// Name returns the category's name.
	Name() string

	// Description returns the category's description.
	Description() string

	// Operations returns all available operations for this category.
	Operations() []Operation

	// IsUIVisible returns whether this category should be visible in the UI.
	IsUIVisible() bool
}

// Operation represents a cloud operation.
type Operation interface {
	// Name returns the operation's name.
	Name() string

	// Description returns the operation's description.
	Description() string

	// Execute executes the operation with the given parameters.
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)

	// IsUIVisible returns whether this operation should be visible in the UI.
	IsUIVisible() bool
}

// UIOperation represents a user-facing operation in the UI
type UIOperation interface {
	// Name returns the operation's name
	Name() string

	// Description returns the operation's description
	Description() string

	// IsUIVisible returns whether this operation should be visible in the UI
	IsUIVisible() bool
}

// ApprovalAction represents a pending approval in a pipeline
type ApprovalAction struct {
	PipelineName string
	StageName    string
	ActionName   string
	Token        string
}

// StageStatus represents the status of a pipeline stage
type StageStatus struct {
	Name        string
	Status      string
	LastUpdated string
}

// PipelineStatus represents the status of a pipeline and its stages
type PipelineStatus struct {
	Name   string
	Stages []StageStatus
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
