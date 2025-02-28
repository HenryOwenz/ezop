package providers

import (
	"fmt"
)

// InitializeProviders registers all available providers with the registry.
func InitializeProviders(registry *ProviderRegistry) {
	// Register AWS provider
	// We'll use a direct approach to avoid import cycles
	awsProvider := CreateAWSProvider()
	if awsProvider != nil {
		registry.Register(awsProvider)
	} else {
		panic("Failed to create AWS provider")
	}

	// TODO: Register other providers as they become available
	// RegisterAzureProvider(registry)
	// RegisterGCPProvider(registry)
}

// CreateAWSProvider creates a new AWS provider.
// This is a placeholder that will be replaced by the actual implementation.
var CreateAWSProvider func() Provider

// CreateProvider creates a provider with the given name, profile, and region.
func CreateProvider(registry *ProviderRegistry, name, profile, region string) (Provider, error) {
	// Get the provider from the registry
	provider, err := registry.Get(name)
	if err != nil {
		return nil, err
	}

	// Load the provider configuration
	err = provider.LoadConfig(profile, region)
	if err != nil {
		return nil, fmt.Errorf("failed to load provider configuration: %w", err)
	}

	return provider, nil
}
