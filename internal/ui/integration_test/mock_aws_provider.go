package integration

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// MockAWSProvider implements cloud.Provider for testing
type MockAWSProvider struct {
	profile string
	region  string
}

// Name returns the provider name
func (p *MockAWSProvider) Name() string {
	return "AWS"
}

// Description returns the provider description
func (p *MockAWSProvider) Description() string {
	return "Amazon Web Services (Mock)"
}

// Services returns available services
func (p *MockAWSProvider) Services() []cloud.Service {
	return []cloud.Service{
		&MockService{
			name:        "CodePipeline",
			description: "AWS CodePipeline (Mock)",
			categories: []cloud.Category{
				&MockServiceCategory{
					name:        "Operations",
					description: "Pipeline Operations",
					operations: []cloud.Operation{
						&MockPipelineStatusOperation{},
						&MockStartPipelineOperation{},
						&MockCodePipelineManualApprovalOperation{},
					},
				},
			},
		},
		&MockService{
			name:        "Lambda",
			description: "AWS Lambda (Mock)",
			categories: []cloud.Category{
				&MockServiceCategory{
					name:        "Operations",
					description: "Function Operations",
					operations: []cloud.Operation{
						&MockFunctionStatusOperation{},
					},
				},
			},
		},
	}
}

// GetProfiles returns available profiles
func (p *MockAWSProvider) GetProfiles() ([]string, error) {
	return []string{"default", "dev", "prod"}, nil
}

// LoadConfig loads the provider configuration
func (p *MockAWSProvider) LoadConfig(profile, region string) error {
	p.profile = profile
	p.region = region
	return nil
}

// GetFunctionStatusOperation returns an operation for viewing Lambda function status
func (p *MockAWSProvider) GetFunctionStatusOperation() (cloud.FunctionStatusOperation, error) {
	return &MockFunctionStatusOperation{}, nil
}

// GetCodePipelineManualApprovalOperation returns an operation for managing pipeline approvals
func (p *MockAWSProvider) GetCodePipelineManualApprovalOperation() (cloud.CodePipelineManualApprovalOperation, error) {
	return &MockCodePipelineManualApprovalOperation{}, nil
}

// GetPipelineStatusOperation returns an operation for viewing pipeline status
func (p *MockAWSProvider) GetPipelineStatusOperation() (cloud.PipelineStatusOperation, error) {
	return &MockPipelineStatusOperation{}, nil
}

// GetStartPipelineOperation returns an operation for starting pipeline execution
func (p *MockAWSProvider) GetStartPipelineOperation() (cloud.StartPipelineOperation, error) {
	return &MockStartPipelineOperation{}, nil
}

// GetAuthenticationMethods returns available authentication methods
func (p *MockAWSProvider) GetAuthenticationMethods() []string {
	return []string{"profile", "access_key"}
}

// GetAuthConfigKeys returns configuration keys for the specified authentication method
func (p *MockAWSProvider) GetAuthConfigKeys(method string) []string {
	switch method {
	case "profile":
		return []string{"profile_name"}
	case "access_key":
		return []string{"access_key_id", "secret_access_key"}
	default:
		return []string{}
	}
}

// Authenticate authenticates using the provided method and configuration
func (p *MockAWSProvider) Authenticate(method string, authConfig map[string]string) error {
	switch method {
	case "profile":
		if _, ok := authConfig["profile_name"]; !ok {
			return fmt.Errorf("profile_name is required")
		}
	case "access_key":
		if _, ok := authConfig["access_key_id"]; !ok {
			return fmt.Errorf("access_key_id is required")
		}
		if _, ok := authConfig["secret_access_key"]; !ok {
			return fmt.Errorf("secret_access_key is required")
		}
	default:
		return fmt.Errorf("unsupported authentication method: %s", method)
	}
	return nil
}

// IsAuthenticated checks if the provider is authenticated
func (p *MockAWSProvider) IsAuthenticated() bool {
	return true
}

// GetConfigKeys returns required configuration keys
func (p *MockAWSProvider) GetConfigKeys() []string {
	return []string{"region"}
}

// GetConfigOptions returns available options for a configuration key
func (p *MockAWSProvider) GetConfigOptions(key string) ([]string, error) {
	if key == "region" {
		return []string{
			"us-east-1",
			"us-east-2",
			"us-west-1",
			"us-west-2",
			"eu-west-1",
			"eu-central-1",
		}, nil
	}
	return []string{}, fmt.Errorf("unknown config key: %s", key)
}

// Configure configures the provider with the given configuration
func (p *MockAWSProvider) Configure(config map[string]string) error {
	if region, ok := config["region"]; ok {
		p.region = region
	} else {
		return fmt.Errorf("region is required")
	}
	return nil
}

// GetApprovals returns pending approvals for the provider
func (p *MockAWSProvider) GetApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	return []cloud.ApprovalAction{
		{
			PipelineName: "mock-pipeline",
			StageName:    "Approval",
			ActionName:   "ManualApproval",
			Token:        "mock-token",
		},
	}, nil
}

// ApproveAction approves or rejects an approval action
func (p *MockAWSProvider) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	// Mock implementation
	return nil
}

// GetStatus returns the status of all pipelines
func (p *MockAWSProvider) GetStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	return []cloud.PipelineStatus{
		{
			Name: "mock-pipeline",
			Stages: []cloud.StageStatus{
				{
					Name:   "Source",
					Status: "Succeeded",
				},
				{
					Name:   "Build",
					Status: "InProgress",
				},
			},
		},
	}, nil
}

// StartPipeline starts a pipeline execution
func (p *MockAWSProvider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	// Mock implementation
	return nil
}

// MockFunctionStatusOperation implements cloud.FunctionStatusOperation for testing
type MockFunctionStatusOperation struct{}

func (o *MockFunctionStatusOperation) Name() string {
	return "View Lambda Function Status"
}

func (o *MockFunctionStatusOperation) Description() string {
	return "View the status of Lambda functions"
}

func (o *MockFunctionStatusOperation) IsUIVisible() bool {
	return true
}

func (o *MockFunctionStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetFunctionStatus(ctx)
}

func (o *MockFunctionStatusOperation) GetFunctionStatus(ctx context.Context) ([]cloud.FunctionStatus, error) {
	return []cloud.FunctionStatus{
		{
			Name:        "mock-function-1",
			Runtime:     "nodejs14.x",
			Memory:      128,
			Timeout:     30,
			LastUpdate:  "2023-01-01T00:00:00Z",
			Role:        "arn:aws:iam::123456789012:role/lambda-role",
			Handler:     "index.handler",
			Description: "Mock function 1",
		},
		{
			Name:        "mock-function-2",
			Runtime:     "python3.9",
			Memory:      256,
			Timeout:     60,
			LastUpdate:  "2023-01-02T00:00:00Z",
			Role:        "arn:aws:iam::123456789012:role/lambda-role",
			Handler:     "app.handler",
			Description: "Mock function 2",
		},
	}, nil
}

// MockCodePipelineManualApprovalOperation implements cloud.CodePipelineManualApprovalOperation for testing
type MockCodePipelineManualApprovalOperation struct{}

func (o *MockCodePipelineManualApprovalOperation) Name() string {
	return "Manual Approval"
}

func (o *MockCodePipelineManualApprovalOperation) Description() string {
	return "Approve or reject pipeline approval actions"
}

func (o *MockCodePipelineManualApprovalOperation) IsUIVisible() bool {
	return true
}

func (o *MockCodePipelineManualApprovalOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetPendingApprovals(ctx)
}

func (o *MockCodePipelineManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	return []cloud.ApprovalAction{
		{
			PipelineName: "mock-pipeline",
			StageName:    "Approval",
			ActionName:   "ManualApproval",
			Token:        "mock-token",
		},
	}, nil
}

func (o *MockCodePipelineManualApprovalOperation) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	return nil
}

// MockPipelineStatusOperation implements cloud.PipelineStatusOperation for testing
type MockPipelineStatusOperation struct{}

func (o *MockPipelineStatusOperation) Name() string {
	return "Pipeline Status"
}

func (o *MockPipelineStatusOperation) Description() string {
	return "View the status of CodePipeline pipelines"
}

func (o *MockPipelineStatusOperation) IsUIVisible() bool {
	return true
}

func (o *MockPipelineStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetPipelineStatus(ctx)
}

func (o *MockPipelineStatusOperation) GetPipelineStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	return []cloud.PipelineStatus{
		{
			Name: "mock-pipeline-1",
			Stages: []cloud.StageStatus{
				{
					Name:   "Source",
					Status: "Succeeded",
				},
				{
					Name:   "Build",
					Status: "InProgress",
				},
			},
		},
		{
			Name: "mock-pipeline-2",
			Stages: []cloud.StageStatus{
				{
					Name:   "Source",
					Status: "Succeeded",
				},
				{
					Name:   "Build",
					Status: "Succeeded",
				},
				{
					Name:   "Deploy",
					Status: "Succeeded",
				},
			},
		},
	}, nil
}

// MockStartPipelineOperation implements cloud.StartPipelineOperation for testing
type MockStartPipelineOperation struct{}

func (o *MockStartPipelineOperation) Name() string {
	return "Start Pipeline"
}

func (o *MockStartPipelineOperation) Description() string {
	return "Start a CodePipeline pipeline execution"
}

func (o *MockStartPipelineOperation) IsUIVisible() bool {
	return true
}

func (o *MockStartPipelineOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	pipelineName, _ := params["pipeline_name"].(string)
	commitID, _ := params["commit_id"].(string)
	return nil, o.StartPipelineExecution(ctx, pipelineName, commitID)
}

func (o *MockStartPipelineOperation) StartPipelineExecution(ctx context.Context, pipelineName, commitID string) error {
	return nil
}

// MockService implements cloud.Service for testing
type MockService struct {
	name        string
	description string
	categories  []cloud.Category
}

// Name returns the service name
func (s *MockService) Name() string {
	return s.name
}

// Description returns the service description
func (s *MockService) Description() string {
	return s.description
}

// Categories returns the service categories
func (s *MockService) Categories() []cloud.Category {
	return s.categories
}

// MockServiceCategory implements cloud.Category for testing
type MockServiceCategory struct {
	name        string
	description string
	operations  []cloud.Operation
}

// Name returns the category name
func (c *MockServiceCategory) Name() string {
	return c.name
}

// Description returns the category description
func (c *MockServiceCategory) Description() string {
	return c.description
}

// Operations returns the category operations
func (c *MockServiceCategory) Operations() []cloud.Operation {
	return c.operations
}

// IsUIVisible returns whether this category should be visible in the UI
func (c *MockServiceCategory) IsUIVisible() bool {
	return true
}
