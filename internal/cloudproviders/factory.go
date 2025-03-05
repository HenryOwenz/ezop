package cloudproviders

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws"
)

// InitializeProviders registers all available providers with the registry
func InitializeProviders(registry *cloud.ProviderRegistry) {
	// Create and register AWS provider
	awsProvider := aws.New()
	wrapper := NewAWSProviderWrapper(awsProvider)
	registry.Register(wrapper)
}

// CreateProvider creates a provider with the given name and configuration
func CreateProvider(registry *cloud.ProviderRegistry, name, profile, region string) (cloud.Provider, error) {
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
