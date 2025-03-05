package providers

import (
	cloudAws "github.com/HenryOwenz/cloudgate/internal/cloud/aws"
)

// InitializeProviders registers all available providers with the registry.
func InitializeProviders(registry *ProviderRegistry) {
	// Register AWS provider
	awsProvider := CreateAWSProvider()
	registry.Register(awsProvider)
}

// CreateAWSProvider creates a new AWS provider.
func CreateAWSProvider() Provider {
	// Create a new cloud provider
	cloudProvider := cloudAws.New()

	// Create a new provider adapter
	return NewProviderAdapter(cloudProvider)
}

// CreateProvider creates a provider with the given name and configuration.
func CreateProvider(registry *ProviderRegistry, name, profile, region string) (Provider, error) {
	provider, err := registry.Get(name)
	if err != nil {
		return nil, NewProviderError(name, err)
	}

	err = provider.LoadConfig(profile, region)
	if err != nil {
		return nil, NewProviderError(name, err)
	}

	return provider, nil
}
