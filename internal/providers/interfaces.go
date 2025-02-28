package providers

import (
	"context"
	"fmt"
	"sync"
)

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

// Provider interface defines methods that all cloud providers must implement
type Provider interface {
	// Name returns the provider's name
	Name() string

	// Description returns the provider's description
	Description() string

	// Services returns all available services for this provider
	Services() []Service

	// GetProfiles returns all available profiles for this provider
	GetProfiles() ([]string, error)

	// LoadConfig loads the provider configuration with the given profile and region
	LoadConfig(profile, region string) error

	// GetAuthenticationMethods returns the available authentication methods
	GetAuthenticationMethods() []string

	// GetAuthConfigKeys returns the configuration keys required for an authentication method
	GetAuthConfigKeys(method string) []string

	// Authenticate authenticates with the provider using the given method and configuration
	Authenticate(method string, authConfig map[string]string) error

	// IsAuthenticated returns whether the provider is authenticated
	IsAuthenticated() bool

	// GetConfigKeys returns the configuration keys required by this provider
	GetConfigKeys() []string

	// GetConfigOptions returns the available options for a configuration key
	GetConfigOptions(key string) ([]string, error)

	// Configure configures the provider with the given configuration
	Configure(config map[string]string) error

	// GetApprovals returns all pending approvals for the provider
	GetApprovals(ctx context.Context) ([]ApprovalAction, error)

	// ApproveAction approves or rejects an approval action
	ApproveAction(ctx context.Context, action ApprovalAction, approved bool, comment string) error

	// GetStatus returns the status of all pipelines
	GetStatus(ctx context.Context) ([]PipelineStatus, error)

	// StartPipeline starts a pipeline execution
	StartPipeline(ctx context.Context, pipelineName string, commitID string) error
}

// Service interface defines methods that all cloud services must implement
type Service interface {
	// Name returns the service's name
	Name() string

	// Description returns the service's description
	Description() string

	// Categories returns all available categories for this service
	Categories() []Category
}

// Category interface defines methods that all service categories must implement
type Category interface {
	// Name returns the category's name
	Name() string

	// Description returns the category's description
	Description() string

	// Operations returns all available operations for this category
	Operations() []Operation

	// IsUIVisible returns whether this category should be visible in the UI
	IsUIVisible() bool
}

// Operation interface defines methods that all service operations must implement
type Operation interface {
	// Name returns the operation's name
	Name() string

	// Description returns the operation's description
	Description() string

	// Execute executes the operation with the given parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)

	// IsUIVisible returns whether this operation should be visible in the UI
	IsUIVisible() bool
}

// Registry for all providers
type ProviderRegistry struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

// Register registers a provider with the registry
func (r *ProviderRegistry) Register(provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider.Name()] = provider
}

// Providers returns all registered providers
func (r *ProviderRegistry) Providers() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}

// GetProvider returns a provider by name
func (r *ProviderRegistry) GetProvider(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[name]
	return provider, ok
}

// Get returns a provider by name with an error if not found
func (r *ProviderRegistry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// List returns all registered providers (alias for Providers)
func (r *ProviderRegistry) List() []Provider {
	return r.Providers()
}

// GetProviderNames returns the names of all registered providers
func (r *ProviderRegistry) GetProviderNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
