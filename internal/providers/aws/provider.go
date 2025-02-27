package aws

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud/aws"
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
)

// Provider adapts the cloud.Provider to the providers.Provider interface
type Provider struct {
	cloudProvider *aws.Provider
	profile       string
	region        string
	authenticated bool
}

// New creates a new AWS provider adapter
func New() *Provider {
	return &Provider{
		cloudProvider: aws.New(),
	}
}

// Name returns the provider's name
func (p *Provider) Name() string {
	return p.cloudProvider.Name()
}

// Description returns the provider's description
func (p *Provider) Description() string {
	return p.cloudProvider.Description()
}

// Services returns all available services for this provider
func (p *Provider) Services() []providers.Service {
	cloudServices := p.cloudProvider.Services()
	services := make([]providers.Service, len(cloudServices))
	for i, service := range cloudServices {
		services[i] = providers.NewCloudServiceAdapter(service)
	}
	return services
}

// GetProfiles returns all available profiles for this provider
func (p *Provider) GetProfiles() ([]string, error) {
	return p.cloudProvider.GetProfiles()
}

// LoadConfig loads the provider configuration with the given profile and region
func (p *Provider) LoadConfig(profile, region string) error {
	p.profile = profile
	p.region = region
	err := p.cloudProvider.LoadConfig(profile, region)
	if err == nil {
		p.authenticated = true
	}
	return err
}

// GetAuthenticationMethods returns the available authentication methods
func (p *Provider) GetAuthenticationMethods() []string {
	// AWS doesn't need explicit authentication methods
	return []string{}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (p *Provider) GetAuthConfigKeys(method string) []string {
	// AWS doesn't need auth config keys
	return []string{}
}

// Authenticate authenticates with the provider using the given method and configuration
func (p *Provider) Authenticate(method string, authConfig map[string]string) error {
	// AWS doesn't need explicit authentication
	return nil
}

// IsAuthenticated returns whether the provider is authenticated
func (p *Provider) IsAuthenticated() bool {
	// AWS is always "authenticated" if we have a profile and region
	return p.profile != "" && p.region != ""
}

// GetConfigKeys returns the configuration keys required by this provider
func (p *Provider) GetConfigKeys() []string {
	return []string{constants.AWSProfileKey, constants.AWSRegionKey}
}

// GetConfigOptions returns the available options for a configuration key
func (p *Provider) GetConfigOptions(key string) ([]string, error) {
	switch key {
	case constants.AWSProfileKey:
		return p.GetProfiles()
	case constants.AWSRegionKey:
		return constants.DefaultAWSRegions, nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", key)
	}
}

// Configure configures the provider with the given configuration
func (p *Provider) Configure(config map[string]string) error {
	profile, ok := config[constants.AWSProfileKey]
	if !ok || profile == "" {
		return fmt.Errorf("profile is required")
	}

	region, ok := config[constants.AWSRegionKey]
	if !ok || region == "" {
		return fmt.Errorf("region is required")
	}

	return p.LoadConfig(profile, region)
}
