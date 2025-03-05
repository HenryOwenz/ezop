package providers

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
)

// ProviderAdapter adapts any cloud.Provider to the providers.Provider interface
type ProviderAdapter struct {
	cloudProvider cloud.Provider
	profile       string
	region        string
}

// NewProviderAdapter creates a new adapter for the given cloud provider
func NewProviderAdapter(cloudProvider cloud.Provider) Provider {
	return &ProviderAdapter{
		cloudProvider: cloudProvider,
	}
}

// Name returns the provider's name
func (a *ProviderAdapter) Name() string {
	return a.cloudProvider.Name()
}

// Description returns the provider's description
func (a *ProviderAdapter) Description() string {
	return a.cloudProvider.Description()
}

// Services returns all available services for this provider
func (a *ProviderAdapter) Services() []Service {
	cloudServices := a.cloudProvider.Services()
	services := make([]Service, len(cloudServices))
	for i, service := range cloudServices {
		services[i] = &ServiceAdapter{cloudService: service}
	}
	return services
}

// GetProfiles returns all available profiles for this provider
func (a *ProviderAdapter) GetProfiles() ([]string, error) {
	return a.cloudProvider.GetProfiles()
}

// LoadConfig loads the provider configuration with the given profile and region
func (a *ProviderAdapter) LoadConfig(profile, region string) error {
	a.profile = profile
	a.region = region
	return a.cloudProvider.LoadConfig(profile, region)
}

// GetAuthenticationMethods returns the available authentication methods
func (a *ProviderAdapter) GetAuthenticationMethods() []string {
	// AWS only supports profile-based authentication for now
	return []string{"profile"}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (a *ProviderAdapter) GetAuthConfigKeys(method string) []string {
	return []string{"profile", "region"}
}

// Authenticate authenticates with the provider using the given method and configuration
func (a *ProviderAdapter) Authenticate(method string, authConfig map[string]string) error {
	profile, ok := authConfig["profile"]
	if !ok {
		return NewProviderError(a.Name(), fmt.Errorf("profile is required"))
	}

	region, ok := authConfig["region"]
	if !ok {
		return NewProviderError(a.Name(), fmt.Errorf("region is required"))
	}

	return a.LoadConfig(profile, region)
}

// IsAuthenticated returns whether the provider is authenticated
func (a *ProviderAdapter) IsAuthenticated() bool {
	return a.profile != "" && a.region != ""
}

// GetConfigKeys returns the configuration keys required by this provider
func (a *ProviderAdapter) GetConfigKeys() []string {
	return []string{constants.AWSProfileKey, constants.AWSRegionKey}
}

// GetConfigOptions returns the available options for a configuration key
func (a *ProviderAdapter) GetConfigOptions(key string) ([]string, error) {
	switch key {
	case constants.AWSProfileKey:
		return a.GetProfiles()
	case constants.AWSRegionKey:
		return constants.DefaultAWSRegions, nil
	default:
		return nil, NewProviderError(a.Name(), fmt.Errorf("unknown config key: %s", key))
	}
}

// Configure configures the provider with the given configuration
func (a *ProviderAdapter) Configure(config map[string]string) error {
	profile, ok := config[constants.AWSProfileKey]
	if !ok || profile == "" {
		return NewProviderError(a.Name(), ErrInvalidConfig)
	}

	region, ok := config[constants.AWSRegionKey]
	if !ok || region == "" {
		return NewProviderError(a.Name(), ErrInvalidConfig)
	}

	return a.LoadConfig(profile, region)
}

// GetFunctionStatusOperation returns the function status operation
func (a *ProviderAdapter) GetFunctionStatusOperation() (FunctionStatusOperation, error) {
	if !a.IsAuthenticated() {
		return nil, NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the cloud layer operation
	cloudOperation, err := a.cloudProvider.GetFunctionStatusOperation()
	if err != nil {
		return nil, NewProviderError(a.Name(), err)
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return &FunctionStatusOperationAdapter{operation: cloudOperation}, nil
}

// GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
func (a *ProviderAdapter) GetCodePipelineManualApprovalOperation() (CodePipelineManualApprovalOperation, error) {
	if !a.IsAuthenticated() {
		return nil, NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the cloud layer operation
	cloudOperation, err := a.cloudProvider.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return nil, NewProviderError(a.Name(), err)
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return &CodePipelineManualApprovalOperationAdapter{operation: cloudOperation}, nil
}

// GetPipelineStatusOperation returns the pipeline status operation
func (a *ProviderAdapter) GetPipelineStatusOperation() (PipelineStatusOperation, error) {
	if !a.IsAuthenticated() {
		return nil, NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the cloud layer operation
	cloudOperation, err := a.cloudProvider.GetPipelineStatusOperation()
	if err != nil {
		return nil, NewProviderError(a.Name(), err)
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return &PipelineStatusOperationAdapter{operation: cloudOperation}, nil
}

// GetStartPipelineOperation returns the start pipeline operation
func (a *ProviderAdapter) GetStartPipelineOperation() (StartPipelineOperation, error) {
	if !a.IsAuthenticated() {
		return nil, NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the cloud layer operation
	cloudOperation, err := a.cloudProvider.GetStartPipelineOperation()
	if err != nil {
		return nil, NewProviderError(a.Name(), err)
	}

	// Create a wrapper that adapts the cloud operation to the provider interface
	return &StartPipelineOperationAdapter{operation: cloudOperation}, nil
}

// GetApprovals returns all pending approvals for the provider
func (a *ProviderAdapter) GetApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	if !a.IsAuthenticated() {
		return nil, NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the CodePipeline manual approval operation
	operation, err := a.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return nil, err
	}

	// Use the operation to get pending approvals
	return operation.GetPendingApprovals(ctx)
}

// ApproveAction approves or rejects an approval action
func (a *ProviderAdapter) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	if !a.IsAuthenticated() {
		return NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the CodePipeline manual approval operation
	operation, err := a.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return err
	}

	// Use the operation to approve or reject the action
	return operation.ApproveAction(ctx, action, approved, comment)
}

// GetStatus returns the status of all pipelines
func (a *ProviderAdapter) GetStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	if !a.IsAuthenticated() {
		return nil, NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the pipeline status operation
	operation, err := a.GetPipelineStatusOperation()
	if err != nil {
		return nil, err
	}

	// Use the operation to get the pipeline status
	return operation.GetPipelineStatus(ctx)
}

// StartPipeline starts a pipeline execution
func (a *ProviderAdapter) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	if !a.IsAuthenticated() {
		return NewProviderError(a.Name(), ErrNotAuthenticated)
	}

	// Get the start pipeline operation
	operation, err := a.GetStartPipelineOperation()
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

// Name returns the service's name
func (a *ServiceAdapter) Name() string {
	return a.cloudService.Name()
}

// Description returns the service's description
func (a *ServiceAdapter) Description() string {
	return a.cloudService.Description()
}

// Categories returns all available categories for this service
func (a *ServiceAdapter) Categories() []Category {
	cloudCategories := a.cloudService.Categories()
	categories := make([]Category, len(cloudCategories))
	for i, category := range cloudCategories {
		categories[i] = &CategoryAdapter{category: category}
	}
	return categories
}

// CategoryAdapter adapts a cloud.Category to a providers.Category
type CategoryAdapter struct {
	category cloud.Category
}

// Name returns the category's name
func (a *CategoryAdapter) Name() string {
	return a.category.Name()
}

// Description returns the category's description
func (a *CategoryAdapter) Description() string {
	return a.category.Description()
}

// Operations returns all available operations for this category
func (a *CategoryAdapter) Operations() []Operation {
	cloudOperations := a.category.Operations()
	operations := make([]Operation, len(cloudOperations))
	for i, operation := range cloudOperations {
		operations[i] = &OperationAdapter{operation: operation}
	}
	return operations
}

// IsUIVisible returns whether this category should be visible in the UI
func (a *CategoryAdapter) IsUIVisible() bool {
	return a.category.IsUIVisible()
}

// OperationAdapter adapts a cloud.Operation to a providers.Operation
type OperationAdapter struct {
	operation cloud.Operation
}

// Name returns the operation's name
func (a *OperationAdapter) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *OperationAdapter) Description() string {
	return a.operation.Description()
}

// Execute executes the operation with the given parameters
func (a *OperationAdapter) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	result, err := a.operation.Execute(ctx, params)
	if err != nil {
		return nil, NewOperationError(a.Name(), err)
	}
	return result, nil
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *OperationAdapter) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// FunctionStatusOperationAdapter adapts a cloud.FunctionStatusOperation to a providers.FunctionStatusOperation
type FunctionStatusOperationAdapter struct {
	operation cloud.FunctionStatusOperation
}

// Name returns the operation's name
func (a *FunctionStatusOperationAdapter) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *FunctionStatusOperationAdapter) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *FunctionStatusOperationAdapter) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// GetFunctionStatus returns the status of all Lambda functions
func (a *FunctionStatusOperationAdapter) GetFunctionStatus(ctx context.Context) ([]cloud.FunctionStatus, error) {
	functions, err := a.operation.GetFunctionStatus(ctx)
	if err != nil {
		return nil, NewOperationError(a.Name(), err)
	}
	return functions, nil
}

// CodePipelineManualApprovalOperationAdapter adapts a cloud.CodePipelineManualApprovalOperation to a providers.CodePipelineManualApprovalOperation
type CodePipelineManualApprovalOperationAdapter struct {
	operation cloud.CodePipelineManualApprovalOperation
}

// Name returns the operation's name
func (a *CodePipelineManualApprovalOperationAdapter) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *CodePipelineManualApprovalOperationAdapter) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *CodePipelineManualApprovalOperationAdapter) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// GetPendingApprovals returns all pending manual approval actions
func (a *CodePipelineManualApprovalOperationAdapter) GetPendingApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	approvals, err := a.operation.GetPendingApprovals(ctx)
	if err != nil {
		return nil, NewOperationError(a.Name(), err)
	}
	return approvals, nil
}

// ApproveAction approves or rejects an approval action
func (a *CodePipelineManualApprovalOperationAdapter) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	err := a.operation.ApproveAction(ctx, action, approved, comment)
	if err != nil {
		return NewOperationError(a.Name(), err)
	}
	return nil
}

// PipelineStatusOperationAdapter adapts a cloud.PipelineStatusOperation to a providers.PipelineStatusOperation
type PipelineStatusOperationAdapter struct {
	operation cloud.PipelineStatusOperation
}

// Name returns the operation's name
func (a *PipelineStatusOperationAdapter) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *PipelineStatusOperationAdapter) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *PipelineStatusOperationAdapter) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// GetPipelineStatus returns the status of all pipelines
func (a *PipelineStatusOperationAdapter) GetPipelineStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	pipelines, err := a.operation.GetPipelineStatus(ctx)
	if err != nil {
		return nil, NewOperationError(a.Name(), err)
	}
	return pipelines, nil
}

// StartPipelineOperationAdapter adapts a cloud.StartPipelineOperation to a providers.StartPipelineOperation
type StartPipelineOperationAdapter struct {
	operation cloud.StartPipelineOperation
}

// Name returns the operation's name
func (a *StartPipelineOperationAdapter) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description
func (a *StartPipelineOperationAdapter) Description() string {
	return a.operation.Description()
}

// IsUIVisible returns whether this operation should be visible in the UI
func (a *StartPipelineOperationAdapter) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}

// StartPipelineExecution starts a pipeline execution
func (a *StartPipelineOperationAdapter) StartPipelineExecution(ctx context.Context, pipelineName string, commitID string) error {
	err := a.operation.StartPipelineExecution(ctx, pipelineName, commitID)
	if err != nil {
		return NewOperationError(a.Name(), err)
	}
	return nil
}
