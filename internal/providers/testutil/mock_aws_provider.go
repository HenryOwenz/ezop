package testutil

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/providers"
)

// MockAWSProvider is a mock implementation of the AWS provider for testing
type MockAWSProvider struct {
	profiles []string
	regions  []string
	services []providers.Service
}

// NewMockAWSProvider creates a new mock AWS provider
func NewMockAWSProvider() providers.Provider {
	return &MockAWSProvider{
		profiles: []string{"default", "dev", "prod"},
		regions:  []string{"us-east-1", "us-west-2", "eu-west-1"},
		services: []providers.Service{
			NewMockCodePipelineService(),
		},
	}
}

// NewMockAWSProviderWithProfiles creates a new mock AWS provider with custom profiles
func NewMockAWSProviderWithProfiles(profiles []string) providers.Provider {
	return &MockAWSProvider{
		profiles: profiles,
		regions:  []string{"us-east-1", "us-west-2", "eu-west-1"},
		services: []providers.Service{
			NewMockCodePipelineService(),
		},
	}
}

// Name returns the name of the provider
func (p *MockAWSProvider) Name() string {
	return "AWS"
}

// Description returns the description of the provider
func (p *MockAWSProvider) Description() string {
	return "Amazon Web Services"
}

// GetProfiles returns the available AWS profiles
func (p *MockAWSProvider) GetProfiles() ([]string, error) {
	return p.profiles, nil
}

// LoadConfig loads the provider configuration
func (p *MockAWSProvider) LoadConfig(profile, region string) error {
	return nil
}

// GetAuthenticationMethods returns the available authentication methods
func (p *MockAWSProvider) GetAuthenticationMethods() []string {
	return []string{"profile"}
}

// GetAuthConfigKeys returns the configuration keys for the given authentication method
func (p *MockAWSProvider) GetAuthConfigKeys(method string) []string {
	return []string{}
}

// Authenticate authenticates the provider with the given method and configuration
func (p *MockAWSProvider) Authenticate(method string, config map[string]string) error {
	return nil
}

// IsAuthenticated returns whether the provider is authenticated
func (p *MockAWSProvider) IsAuthenticated() bool {
	return true
}

// GetConfigKeys returns the configuration keys for the provider
func (p *MockAWSProvider) GetConfigKeys() []string {
	return []string{}
}

// GetConfigOptions returns the available options for the given configuration key
func (p *MockAWSProvider) GetConfigOptions(key string) ([]string, error) {
	switch key {
	case "profile":
		return p.profiles, nil
	case "region":
		return p.regions, nil
	default:
		return []string{}, nil
	}
}

// Configure configures the provider with the given configuration
func (p *MockAWSProvider) Configure(config map[string]string) error {
	return nil
}

// Services returns the available services
func (p *MockAWSProvider) Services() []providers.Service {
	return p.services
}

// GetApprovals returns all pending approvals for the provider
func (p *MockAWSProvider) GetApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
	return []providers.ApprovalAction{
		{
			PipelineName: "TestPipeline",
			StageName:    "TestStage",
			ActionName:   "TestAction",
			Token:        "TestToken",
		},
	}, nil
}

// ApproveAction approves or rejects an approval action
func (p *MockAWSProvider) ApproveAction(ctx context.Context, action providers.ApprovalAction, approved bool, comment string) error {
	return nil
}

// GetStatus returns the status of all pipelines
func (p *MockAWSProvider) GetStatus(ctx context.Context) ([]providers.PipelineStatus, error) {
	return []providers.PipelineStatus{
		{
			Name: "TestPipeline",
			Stages: []providers.StageStatus{
				{
					Name:        "TestStage",
					Status:      "Succeeded",
					LastUpdated: "2023-01-01 12:00:00",
				},
			},
		},
	}, nil
}

// StartPipeline starts a pipeline execution
func (p *MockAWSProvider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	return nil
}

// GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
func (p *MockAWSProvider) GetCodePipelineManualApprovalOperation() (providers.CodePipelineManualApprovalOperation, error) {
	return &MockCodePipelineManualApprovalOperation{}, nil
}

// GetPipelineStatusOperation returns the pipeline status operation
func (p *MockAWSProvider) GetPipelineStatusOperation() (providers.PipelineStatusOperation, error) {
	return &MockPipelineStatusOperation{}, nil
}

// GetStartPipelineOperation returns the start pipeline operation
func (p *MockAWSProvider) GetStartPipelineOperation() (providers.StartPipelineOperation, error) {
	return &MockStartPipelineOperation{}, nil
}

// GetFunctionStatusOperation returns the function status operation
func (p *MockAWSProvider) GetFunctionStatusOperation() (providers.FunctionStatusOperation, error) {
	return &MockFunctionStatusOperation{}, nil
}

// MockCodePipelineService is a mock implementation of the CodePipeline service
type MockCodePipelineService struct{}

// NewMockCodePipelineService creates a new mock CodePipeline service
func NewMockCodePipelineService() providers.Service {
	return &MockCodePipelineService{}
}

// Name returns the name of the service
func (s *MockCodePipelineService) Name() string {
	return "CodePipeline"
}

// Description returns the description of the service
func (s *MockCodePipelineService) Description() string {
	return "AWS CodePipeline"
}

// Categories returns the available categories
func (s *MockCodePipelineService) Categories() []providers.Category {
	return []providers.Category{
		&MockOperationsCategory{},
	}
}

// MockOperationsCategory is a mock implementation of the Operations category
type MockOperationsCategory struct{}

// Name returns the name of the category
func (c *MockOperationsCategory) Name() string {
	return "Operations"
}

// Description returns the description of the category
func (c *MockOperationsCategory) Description() string {
	return "Pipeline operations"
}

// Operations returns the available operations
func (c *MockOperationsCategory) Operations() []providers.Operation {
	return []providers.Operation{
		&MockOperation{
			name:        "Manual Approval",
			description: "Approve or reject pipeline stages",
		},
		&MockOperation{
			name:        "Pipeline Status",
			description: "View pipeline status",
		},
		&MockOperation{
			name:        "Start Pipeline",
			description: "Start a pipeline execution",
		},
	}
}

// IsUIVisible returns whether the category is visible in the UI
func (c *MockOperationsCategory) IsUIVisible() bool {
	return true
}

// MockOperation is a mock implementation of the Operation interface
type MockOperation struct {
	name        string
	description string
}

// Name returns the name of the operation
func (o *MockOperation) Name() string {
	return o.name
}

// Description returns the description of the operation
func (o *MockOperation) Description() string {
	return o.description
}

// Execute executes the operation with the given parameters
func (o *MockOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}

// IsUIVisible returns whether the operation is visible in the UI
func (o *MockOperation) IsUIVisible() bool {
	return true
}

// MockCodePipelineManualApprovalOperation is a mock implementation of the CodePipeline manual approval operation
type MockCodePipelineManualApprovalOperation struct{}

// Name returns the name of the operation
func (o *MockCodePipelineManualApprovalOperation) Name() string {
	return "Manual Approval"
}

// Description returns the description of the operation
func (o *MockCodePipelineManualApprovalOperation) Description() string {
	return "Approve or reject pipeline stages"
}

// IsUIVisible returns whether the operation is visible in the UI
func (o *MockCodePipelineManualApprovalOperation) IsUIVisible() bool {
	return true
}

// GetPendingApprovals returns the pending approvals
func (o *MockCodePipelineManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
	return []providers.ApprovalAction{
		{
			PipelineName: "TestPipeline",
			StageName:    "TestStage",
			ActionName:   "TestAction",
			Token:        "TestToken",
		},
	}, nil
}

// ApproveAction approves or rejects an action
func (o *MockCodePipelineManualApprovalOperation) ApproveAction(ctx context.Context, approval providers.ApprovalAction, approve bool, comment string) error {
	return nil
}

// MockPipelineStatusOperation is a mock implementation of the Pipeline Status operation
type MockPipelineStatusOperation struct{}

// Name returns the name of the operation
func (o *MockPipelineStatusOperation) Name() string {
	return "Pipeline Status"
}

// Description returns the description of the operation
func (o *MockPipelineStatusOperation) Description() string {
	return "View pipeline status"
}

// IsUIVisible returns whether the operation is visible in the UI
func (o *MockPipelineStatusOperation) IsUIVisible() bool {
	return true
}

// GetPipelineStatus returns the pipeline status
func (o *MockPipelineStatusOperation) GetPipelineStatus(ctx context.Context) ([]providers.PipelineStatus, error) {
	return []providers.PipelineStatus{
		{
			Name: "TestPipeline",
			Stages: []providers.StageStatus{
				{
					Name:        "TestStage",
					Status:      "Succeeded",
					LastUpdated: "2023-01-01 12:00:00",
				},
			},
		},
	}, nil
}

// MockStartPipelineOperation is a mock implementation of the Start Pipeline operation
type MockStartPipelineOperation struct{}

// Name returns the name of the operation
func (o *MockStartPipelineOperation) Name() string {
	return "Start Pipeline"
}

// Description returns the description of the operation
func (o *MockStartPipelineOperation) Description() string {
	return "Start a pipeline execution"
}

// IsUIVisible returns whether the operation is visible in the UI
func (o *MockStartPipelineOperation) IsUIVisible() bool {
	return true
}

// StartPipelineExecution starts a pipeline execution
func (o *MockStartPipelineOperation) StartPipelineExecution(ctx context.Context, pipelineName, commitID string) error {
	return nil
}

// MockFunctionStatusOperation is a mock implementation of the FunctionStatusOperation interface
type MockFunctionStatusOperation struct{}

// Name returns the name of the operation
func (o *MockFunctionStatusOperation) Name() string {
	return "Function Status"
}

// Description returns the description of the operation
func (o *MockFunctionStatusOperation) Description() string {
	return "View Lambda function status"
}

// IsUIVisible returns whether this operation should be visible in the UI
func (o *MockFunctionStatusOperation) IsUIVisible() bool {
	return true
}

// GetFunctionStatus returns mock function status data
func (o *MockFunctionStatusOperation) GetFunctionStatus(ctx context.Context) ([]providers.FunctionStatus, error) {
	return []providers.FunctionStatus{
		{
			Name:       "test-function",
			Runtime:    "nodejs14.x",
			Memory:     128,
			Timeout:    30,
			LastUpdate: "2023-01-01",
		},
	}, nil
}
