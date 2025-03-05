package cloudproviders

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/codepipeline"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/lambda"
)

// AWSProviderWrapper wraps the AWS provider to ensure it implements the cloud.Provider interface
type AWSProviderWrapper struct {
	provider *aws.Provider
	services []cloud.Service
	profile  string
	region   string
}

// NewAWSProviderWrapper creates a new wrapper for the AWS provider
func NewAWSProviderWrapper(provider *aws.Provider) *AWSProviderWrapper {
	return &AWSProviderWrapper{
		provider: provider,
		services: []cloud.Service{},
	}
}

// Name returns the provider's name
func (w *AWSProviderWrapper) Name() string {
	return w.provider.Name()
}

// Description returns the provider's description
func (w *AWSProviderWrapper) Description() string {
	return w.provider.Description()
}

// Services returns all available services for this provider
func (w *AWSProviderWrapper) Services() []cloud.Service {
	return w.services
}

// GetProfiles returns all available profiles for this provider
func (w *AWSProviderWrapper) GetProfiles() ([]string, error) {
	return w.provider.GetProfiles()
}

// LoadConfig loads the provider configuration with the given profile and region
func (w *AWSProviderWrapper) LoadConfig(profile, region string) error {
	// Load the config in the wrapped provider
	err := w.provider.LoadConfig(profile, region)
	if err != nil {
		return err
	}

	// Store the profile and region
	w.profile = profile
	w.region = region

	// Let the provider load its services
	// Then we'll create our own service wrappers
	lambdaService := &LambdaServiceWrapper{
		service: lambda.NewService(profile, region),
	}

	codePipelineService := &CodePipelineServiceWrapper{
		service: codepipeline.NewService(profile, region),
	}

	// Set the services
	w.services = []cloud.Service{
		lambdaService,
		codePipelineService,
	}

	return nil
}

// LambdaServiceWrapper wraps the Lambda service
type LambdaServiceWrapper struct {
	service *lambda.Service
}

// Name returns the service name
func (s *LambdaServiceWrapper) Name() string {
	return s.service.Name()
}

// Description returns the service description
func (s *LambdaServiceWrapper) Description() string {
	return s.service.Description()
}

// Categories returns the service categories
func (s *LambdaServiceWrapper) Categories() []cloud.Category {
	// Convert the categories
	categories := s.service.Categories()
	result := make([]cloud.Category, len(categories))
	// Using a loop instead of copy due to type compatibility issues
	//nolint:gosimple // Cannot use copy due to type system limitations
	for i := range categories {
		result[i] = categories[i]
	}
	return result
}

// CodePipelineServiceWrapper wraps the CodePipeline service
type CodePipelineServiceWrapper struct {
	service *codepipeline.Service
}

// Name returns the service name
func (s *CodePipelineServiceWrapper) Name() string {
	return s.service.Name()
}

// Description returns the service description
func (s *CodePipelineServiceWrapper) Description() string {
	return s.service.Description()
}

// Categories returns the service categories
func (s *CodePipelineServiceWrapper) Categories() []cloud.Category {
	// Convert the categories
	categories := s.service.Categories()
	result := make([]cloud.Category, len(categories))
	// Using a loop instead of copy due to type compatibility issues
	//nolint:gosimple // Cannot use copy due to type system limitations
	for i := range categories {
		result[i] = categories[i]
	}
	return result
}

// GetFunctionStatusOperation returns the function status operation
func (w *AWSProviderWrapper) GetFunctionStatusOperation() (cloud.FunctionStatusOperation, error) {
	return w.provider.GetFunctionStatusOperation()
}

// GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
func (w *AWSProviderWrapper) GetCodePipelineManualApprovalOperation() (cloud.CodePipelineManualApprovalOperation, error) {
	return w.provider.GetCodePipelineManualApprovalOperation()
}

// GetPipelineStatusOperation returns the pipeline status operation
func (w *AWSProviderWrapper) GetPipelineStatusOperation() (cloud.PipelineStatusOperation, error) {
	return w.provider.GetPipelineStatusOperation()
}

// GetStartPipelineOperation returns the start pipeline operation
func (w *AWSProviderWrapper) GetStartPipelineOperation() (cloud.StartPipelineOperation, error) {
	return w.provider.GetStartPipelineOperation()
}

// GetAuthenticationMethods returns the available authentication methods
func (w *AWSProviderWrapper) GetAuthenticationMethods() []string {
	return w.provider.GetAuthenticationMethods()
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (w *AWSProviderWrapper) GetAuthConfigKeys(method string) []string {
	return w.provider.GetAuthConfigKeys(method)
}

// Authenticate authenticates with the provider using the given method and configuration
func (w *AWSProviderWrapper) Authenticate(method string, authConfig map[string]string) error {
	return w.provider.Authenticate(method, authConfig)
}

// IsAuthenticated returns whether the provider is authenticated
func (w *AWSProviderWrapper) IsAuthenticated() bool {
	return w.provider.IsAuthenticated()
}

// GetConfigKeys returns the configuration keys required by this provider
func (w *AWSProviderWrapper) GetConfigKeys() []string {
	return w.provider.GetConfigKeys()
}

// GetConfigOptions returns the available options for a configuration key
func (w *AWSProviderWrapper) GetConfigOptions(key string) ([]string, error) {
	return w.provider.GetConfigOptions(key)
}

// Configure configures the provider with the given configuration
func (w *AWSProviderWrapper) Configure(config map[string]string) error {
	return w.provider.Configure(config)
}

// GetApprovals returns all pending approvals for the provider
func (w *AWSProviderWrapper) GetApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	// Get the approval operation
	op, err := w.provider.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return nil, err
	}

	// Use the operation to get the approvals
	return op.GetPendingApprovals(ctx)
}

// ApproveAction approves or rejects an approval action
func (w *AWSProviderWrapper) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	return w.provider.ApproveAction(ctx, action, approved, comment)
}

// GetStatus returns the status of all pipelines
func (w *AWSProviderWrapper) GetStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	// Get the pipeline status operation
	op, err := w.provider.GetPipelineStatusOperation()
	if err != nil {
		return nil, err
	}

	// Use the operation to get the pipeline status
	return op.GetPipelineStatus(ctx)
}

// StartPipeline starts a pipeline execution
func (w *AWSProviderWrapper) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	return w.provider.StartPipeline(ctx, pipelineName, commitID)
}
