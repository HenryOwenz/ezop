package aws

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/providers"
)

// Provider adapts a cloud.Provider to a providers.Provider
type Provider struct {
	cloudProvider cloud.Provider
	profile       string
	region        string
}

// New creates a new AWS provider.
func New(cloudProvider cloud.Provider) *Provider {
	return &Provider{
		cloudProvider: cloudProvider,
	}
}

// Name returns the provider's name.
func (p *Provider) Name() string {
	return p.cloudProvider.Name()
}

// Description returns the provider's description.
func (p *Provider) Description() string {
	return p.cloudProvider.Description()
}

// Services returns all available services for this provider.
func (p *Provider) Services() []providers.Service {
	// Adapt cloud services to provider services
	cloudServices := p.cloudProvider.Services()
	providerServices := make([]providers.Service, 0, len(cloudServices))

	for _, cloudService := range cloudServices {
		// Create a service adapter for each cloud service
		serviceAdapter := &ServiceAdapter{
			cloudService: cloudService,
		}
		providerServices = append(providerServices, serviceAdapter)
	}

	return providerServices
}

// GetProfiles returns all available profiles for this provider.
func (p *Provider) GetProfiles() ([]string, error) {
	return p.cloudProvider.GetProfiles()
}

// LoadConfig loads the provider configuration with the given profile and region.
func (p *Provider) LoadConfig(profile, region string) error {
	p.profile = profile
	p.region = region
	return p.cloudProvider.LoadConfig(profile, region)
}

// GetAuthenticationMethods returns all available authentication methods for this provider.
func (p *Provider) GetAuthenticationMethods() []string {
	// AWS only supports profile-based authentication for now
	return []string{"profile"}
}

// GetAuthConfigKeys returns the configuration keys required for authentication.
func (p *Provider) GetAuthConfigKeys(method string) []string {
	return []string{"profile", "region"}
}

// Authenticate authenticates the provider with the given credentials.
func (p *Provider) Authenticate(method string, authConfig map[string]string) error {
	profile, ok := authConfig["profile"]
	if !ok {
		return fmt.Errorf("profile is required")
	}

	region, ok := authConfig["region"]
	if !ok {
		return fmt.Errorf("region is required")
	}

	return p.LoadConfig(profile, region)
}

// IsAuthenticated returns whether the provider is authenticated.
func (p *Provider) IsAuthenticated() bool {
	return p.profile != "" && p.region != ""
}

// GetConfigKeys returns the configuration keys for this provider.
func (p *Provider) GetConfigKeys() []string {
	return []string{}
}

// GetConfigOptions returns the configuration options for the given key.
func (p *Provider) GetConfigOptions(key string) ([]string, error) {
	return []string{}, nil
}

// Configure configures the provider with the given configuration.
func (p *Provider) Configure(config map[string]string) error {
	return nil
}

// GetFunctionStatusOperation returns the function status operation
func (p *Provider) GetFunctionStatusOperation() (providers.FunctionStatusOperation, error) {
	if !p.IsAuthenticated() {
		return nil, fmt.Errorf("provider not authenticated")
	}

	// Get the cloud layer operation
	cloudOperation, err := p.cloudProvider.GetFunctionStatusOperation()
	if err != nil {
		return nil, err
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return newFunctionStatusOperation(cloudOperation), nil
}

// GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
func (p *Provider) GetCodePipelineManualApprovalOperation() (providers.CodePipelineManualApprovalOperation, error) {
	if !p.IsAuthenticated() {
		return nil, fmt.Errorf("provider not authenticated")
	}

	// Get the cloud layer operation
	cloudOperation, err := p.cloudProvider.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return nil, err
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return newCodePipelineManualApprovalOperation(cloudOperation), nil
}

// GetPipelineStatusOperation returns the pipeline status operation
func (p *Provider) GetPipelineStatusOperation() (providers.PipelineStatusOperation, error) {
	if !p.IsAuthenticated() {
		return nil, fmt.Errorf("provider not authenticated")
	}

	// Get the cloud layer operation
	cloudOperation, err := p.cloudProvider.GetPipelineStatusOperation()
	if err != nil {
		return nil, err
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return newPipelineStatusOperation(cloudOperation), nil
}

// GetStartPipelineOperation returns the start pipeline operation
func (p *Provider) GetStartPipelineOperation() (providers.StartPipelineOperation, error) {
	if !p.IsAuthenticated() {
		return nil, fmt.Errorf("provider not authenticated")
	}

	// Get the cloud layer operation
	cloudOperation, err := p.cloudProvider.GetStartPipelineOperation()
	if err != nil {
		return nil, err
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return newStartPipelineOperation(cloudOperation), nil
}

// GetApprovals returns all pending approvals for the provider
func (p *Provider) GetApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	if !p.IsAuthenticated() {
		return nil, fmt.Errorf("provider not authenticated")
	}

	// Get the CodePipeline manual approval operation
	operation, err := p.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return nil, err
	}

	// Use the operation to get pending approvals
	return operation.GetPendingApprovals(ctx)
}

// ApproveAction approves or rejects an approval action
func (p *Provider) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	if !p.IsAuthenticated() {
		return fmt.Errorf("provider not authenticated")
	}

	// Get the CodePipeline manual approval operation
	operation, err := p.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return err
	}

	// Use the operation to approve or reject the action
	return operation.ApproveAction(ctx, action, approved, comment)
}

// GetStatus returns the status of all pipelines
func (p *Provider) GetStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	if !p.IsAuthenticated() {
		return nil, fmt.Errorf("provider not authenticated")
	}

	// Get the pipeline status operation
	operation, err := p.GetPipelineStatusOperation()
	if err != nil {
		return nil, err
	}

	// Use the operation to get the pipeline status
	return operation.GetPipelineStatus(ctx)
}

// StartPipeline starts a pipeline execution
func (p *Provider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	if !p.IsAuthenticated() {
		return fmt.Errorf("provider not authenticated")
	}

	// Get the start pipeline operation
	operation, err := p.GetStartPipelineOperation()
	if err != nil {
		return err
	}

	// Use the operation to start the pipeline
	return operation.StartPipelineExecution(ctx, pipelineName, commitID)
}

// ServiceAdapter adapts a cloud.Service to a providers.Service
type ServiceAdapter struct {
	cloudService cloud.Service
}

// Name returns the service's name.
func (s *ServiceAdapter) Name() string {
	return s.cloudService.Name()
}

// Description returns the service's description.
func (s *ServiceAdapter) Description() string {
	return s.cloudService.Description()
}

// Categories returns all available categories for this service.
func (s *ServiceAdapter) Categories() []providers.Category {
	cloudCategories := s.cloudService.Categories()
	providerCategories := make([]providers.Category, 0, len(cloudCategories))

	for _, cloudCategory := range cloudCategories {
		// Create a category adapter for each cloud category
		categoryAdapter := &CategoryAdapter{
			cloudCategory: cloudCategory,
		}
		providerCategories = append(providerCategories, categoryAdapter)
	}

	return providerCategories
}

// CategoryAdapter adapts a cloud.Category to a providers.Category
type CategoryAdapter struct {
	cloudCategory cloud.Category
}

// Name returns the category's name.
func (c *CategoryAdapter) Name() string {
	return c.cloudCategory.Name()
}

// Description returns the category's description.
func (c *CategoryAdapter) Description() string {
	return c.cloudCategory.Description()
}

// Operations returns all available operations for this category.
func (c *CategoryAdapter) Operations() []providers.Operation {
	cloudOperations := c.cloudCategory.Operations()
	providerOperations := make([]providers.Operation, 0, len(cloudOperations))

	for _, cloudOperation := range cloudOperations {
		// Create an operation adapter for each cloud operation
		operationAdapter := &OperationAdapter{
			cloudOperation: cloudOperation,
		}
		providerOperations = append(providerOperations, operationAdapter)
	}

	return providerOperations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (c *CategoryAdapter) IsUIVisible() bool {
	return c.cloudCategory.IsUIVisible()
}

// OperationAdapter adapts a cloud.Operation to a providers.Operation
type OperationAdapter struct {
	cloudOperation cloud.Operation
}

// Name returns the operation's name.
func (o *OperationAdapter) Name() string {
	return o.cloudOperation.Name()
}

// Description returns the operation's description.
func (o *OperationAdapter) Description() string {
	return o.cloudOperation.Description()
}

// Execute executes the operation with the given parameters.
func (o *OperationAdapter) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.cloudOperation.Execute(ctx, params)
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *OperationAdapter) IsUIVisible() bool {
	return o.cloudOperation.IsUIVisible()
}
