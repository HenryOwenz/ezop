package providers

import (
	"context"
	"fmt"
	"sync"
)

// Global registry for AWS provider replacement
var (
	awsProviderMu    sync.RWMutex
	realAWSProvider  Provider
	awsProviderReady bool
)

// RegisterAWSProvider registers the real AWS provider implementation
// This is called from the AWS provider package's init function
func RegisterAWSProvider(provider Provider) {
	awsProviderMu.Lock()
	defer awsProviderMu.Unlock()
	realAWSProvider = provider
	awsProviderReady = true
}

// InitializeProviders registers all available providers with the registry.
func InitializeProviders(registry *ProviderRegistry) {
	// Register AWS provider
	awsProviderMu.RLock()
	if awsProviderReady && realAWSProvider != nil {
		// Use the real AWS provider if it's been registered
		registry.Register(realAWSProvider)
	} else {
		// Fall back to the placeholder if the real provider hasn't been registered yet
		registry.Register(CreateAWSProvider())
	}
	awsProviderMu.RUnlock()
}

// CreateAWSProvider creates a new AWS provider.
func CreateAWSProvider() Provider {
	// Create a new AWS provider
	provider := &awsProvider{}

	return provider
}

// awsProvider is a placeholder for the AWS provider.
// This will be replaced by the actual AWS provider in the aws package.
type awsProvider struct{}

func (p *awsProvider) Name() string {
	return "AWS"
}

func (p *awsProvider) Description() string {
	return "Amazon Web Services"
}

func (p *awsProvider) Services() []Service {
	return []Service{}
}

func (p *awsProvider) GetProfiles() ([]string, error) {
	return []string{}, nil
}

func (p *awsProvider) LoadConfig(profile, region string) error {
	return nil
}

func (p *awsProvider) GetAuthenticationMethods() []string {
	return []string{}
}

func (p *awsProvider) GetAuthConfigKeys(method string) []string {
	return []string{}
}

func (p *awsProvider) Authenticate(method string, authConfig map[string]string) error {
	return nil
}

func (p *awsProvider) IsAuthenticated() bool {
	return false
}

func (p *awsProvider) GetConfigKeys() []string {
	return []string{}
}

func (p *awsProvider) GetConfigOptions(key string) ([]string, error) {
	return []string{}, nil
}

func (p *awsProvider) Configure(config map[string]string) error {
	return nil
}

func (p *awsProvider) GetApprovals(ctx context.Context) ([]ApprovalAction, error) {
	return []ApprovalAction{}, nil
}

func (p *awsProvider) ApproveAction(ctx context.Context, action ApprovalAction, approved bool, comment string) error {
	return nil
}

func (p *awsProvider) GetStatus(ctx context.Context) ([]PipelineStatus, error) {
	return []PipelineStatus{}, nil
}

func (p *awsProvider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	return nil
}

func (p *awsProvider) GetCodePipelineManualApprovalOperation() (CodePipelineManualApprovalOperation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *awsProvider) GetPipelineStatusOperation() (PipelineStatusOperation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *awsProvider) GetStartPipelineOperation() (StartPipelineOperation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *awsProvider) GetFunctionStatusOperation() (FunctionStatusOperation, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateProvider creates a provider with the given name and configuration.
func CreateProvider(registry *ProviderRegistry, name, profile, region string) (Provider, error) {
	provider, err := registry.Get(name)
	if err != nil {
		return nil, err
	}

	err = provider.LoadConfig(profile, region)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
