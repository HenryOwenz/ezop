package cloud

import (
	"errors"
	"sync"
)

// Common errors
var (
	ErrProviderNotFound = errors.New("provider not found")
)

// ProviderRegistry is a registry for cloud providers
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
		return nil, ErrProviderNotFound
	}
	return provider, nil
}

// List returns all registered providers
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
