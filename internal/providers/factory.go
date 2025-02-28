package providers

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud/aws"
)

// InitializeProviders registers all available providers with the registry.
func InitializeProviders(registry *ProviderRegistry) {
	// Register AWS provider
	registry.Register(NewCloudProviderAdapter(aws.New()))

	// TODO: Register other providers as they become available
	// registry.Register(NewCloudProviderAdapter(azure.New()))
	// registry.Register(NewCloudProviderAdapter(gcp.New()))
}

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
