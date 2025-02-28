package providers

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
)

// CloudProviderAdapter adapts a cloud.Provider to a providers.Provider.
type CloudProviderAdapter struct {
	provider cloud.Provider
	profile  string
	region   string
}

// NewCloudProviderAdapter creates a new adapter for a cloud.Provider.
func NewCloudProviderAdapter(provider cloud.Provider) *CloudProviderAdapter {
	return &CloudProviderAdapter{
		provider: provider,
	}
}

// Name returns the provider's name.
func (a *CloudProviderAdapter) Name() string {
	return a.provider.Name()
}

// Description returns the provider's description.
func (a *CloudProviderAdapter) Description() string {
	return a.provider.Description()
}

// Services returns all available services for this provider.
func (a *CloudProviderAdapter) Services() []Service {
	cloudServices := a.provider.Services()
	services := make([]Service, len(cloudServices))
	for i, service := range cloudServices {
		services[i] = &CloudServiceAdapter{service: service}
	}
	return services
}

// GetProfiles returns all available profiles for this provider.
func (a *CloudProviderAdapter) GetProfiles() ([]string, error) {
	return a.provider.GetProfiles()
}

// LoadConfig loads the provider configuration with the given profile and region.
func (a *CloudProviderAdapter) LoadConfig(profile, region string) error {
	a.profile = profile
	a.region = region
	return a.provider.LoadConfig(profile, region)
}

// GetAuthenticationMethods returns the available authentication methods
func (a *CloudProviderAdapter) GetAuthenticationMethods() []string {
	// AWS doesn't need explicit authentication methods
	return []string{}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (a *CloudProviderAdapter) GetAuthConfigKeys(method string) []string {
	// AWS doesn't need auth config keys
	return []string{}
}

// Authenticate authenticates with the provider using the given method and configuration
func (a *CloudProviderAdapter) Authenticate(method string, authConfig map[string]string) error {
	// AWS doesn't need explicit authentication
	return nil
}

// IsAuthenticated returns whether the provider is authenticated
func (a *CloudProviderAdapter) IsAuthenticated() bool {
	// AWS is always "authenticated" if we have a profile and region
	return a.profile != "" && a.region != ""
}

// GetConfigKeys returns the configuration keys required by this provider
func (a *CloudProviderAdapter) GetConfigKeys() []string {
	return []string{constants.AWSProfileKey, constants.AWSRegionKey}
}

// GetConfigOptions returns the available options for a configuration key
func (a *CloudProviderAdapter) GetConfigOptions(key string) ([]string, error) {
	switch key {
	case constants.AWSProfileKey:
		return a.GetProfiles()
	case constants.AWSRegionKey:
		return constants.DefaultAWSRegions, nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", key)
	}
}

// Configure configures the provider with the given configuration
func (a *CloudProviderAdapter) Configure(config map[string]string) error {
	profile, ok := config[constants.AWSProfileKey]
	if !ok || profile == "" {
		return fmt.Errorf("profile is required")
	}

	region, ok := config[constants.AWSRegionKey]
	if !ok || region == "" {
		return fmt.Errorf("region is required")
	}

	return a.LoadConfig(profile, region)
}

// CloudServiceAdapter adapts a cloud.Service to a providers.Service.
type CloudServiceAdapter struct {
	service cloud.Service
}

// NewCloudServiceAdapter creates a new adapter for a cloud.Service.
func NewCloudServiceAdapter(service cloud.Service) *CloudServiceAdapter {
	return &CloudServiceAdapter{
		service: service,
	}
}

// Name returns the service's name.
func (a *CloudServiceAdapter) Name() string {
	return a.service.Name()
}

// Description returns the service's description.
func (a *CloudServiceAdapter) Description() string {
	return a.service.Description()
}

// Categories returns all available categories for this service.
func (a *CloudServiceAdapter) Categories() []Category {
	cloudCategories := a.service.Categories()
	categories := make([]Category, len(cloudCategories))
	for i, category := range cloudCategories {
		categories[i] = &CloudCategoryAdapter{category: category}
	}
	return categories
}

// CloudCategoryAdapter adapts a cloud.Category to a providers.Category.
type CloudCategoryAdapter struct {
	category cloud.Category
}

// NewCloudCategoryAdapter creates a new adapter for a cloud.Category.
func NewCloudCategoryAdapter(category cloud.Category) *CloudCategoryAdapter {
	return &CloudCategoryAdapter{
		category: category,
	}
}

// Name returns the category's name.
func (a *CloudCategoryAdapter) Name() string {
	return a.category.Name()
}

// Description returns the category's description.
func (a *CloudCategoryAdapter) Description() string {
	return a.category.Description()
}

// Operations returns all available operations for this category.
func (a *CloudCategoryAdapter) Operations() []Operation {
	cloudOperations := a.category.Operations()
	operations := make([]Operation, len(cloudOperations))
	for i, operation := range cloudOperations {
		operations[i] = &CloudOperationAdapter{operation: operation}
	}
	return operations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (a *CloudCategoryAdapter) IsUIVisible() bool {
	return a.category.IsUIVisible()
}

// CloudOperationAdapter adapts a cloud.Operation to a providers.Operation.
type CloudOperationAdapter struct {
	operation cloud.Operation
}

// NewCloudOperationAdapter creates a new adapter for a cloud.Operation.
func NewCloudOperationAdapter(operation cloud.Operation) *CloudOperationAdapter {
	return &CloudOperationAdapter{
		operation: operation,
	}
}

// Name returns the operation's name.
func (a *CloudOperationAdapter) Name() string {
	return a.operation.Name()
}

// Description returns the operation's description.
func (a *CloudOperationAdapter) Description() string {
	return a.operation.Description()
}

// Execute executes the operation with the given parameters.
func (a *CloudOperationAdapter) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return a.operation.Execute(ctx, params)
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (a *CloudOperationAdapter) IsUIVisible() bool {
	return a.operation.IsUIVisible()
}
