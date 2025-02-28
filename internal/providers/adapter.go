package providers

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// CloudProviderAdapter adapts a cloud.Provider to a providers.Provider.
type CloudProviderAdapter struct {
	provider cloud.Provider
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
	return a.provider.LoadConfig(profile, region)
}

// CloudServiceAdapter adapts a cloud.Service to a providers.Service.
type CloudServiceAdapter struct {
	service cloud.Service
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
