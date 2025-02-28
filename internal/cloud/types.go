package cloud

import (
	"context"
)

// Provider represents a cloud provider.
type Provider interface {
	// Name returns the provider's name.
	Name() string

	// Description returns the provider's description.
	Description() string

	// Services returns all available services for this provider.
	Services() []Service

	// GetProfiles returns all available profiles for this provider.
	GetProfiles() ([]string, error)

	// LoadConfig loads the provider configuration with the given profile and region.
	LoadConfig(profile, region string) error
}

// Service represents a cloud service.
type Service interface {
	// Name returns the service's name.
	Name() string

	// Description returns the service's description.
	Description() string

	// Categories returns all available categories for this service.
	Categories() []Category
}

// Category represents a group of operations.
type Category interface {
	// Name returns the category's name.
	Name() string

	// Description returns the category's description.
	Description() string

	// Operations returns all available operations for this category.
	Operations() []Operation

	// IsUIVisible returns whether this category should be visible in the UI.
	IsUIVisible() bool
}

// Operation represents a cloud operation.
type Operation interface {
	// Name returns the operation's name.
	Name() string

	// Description returns the operation's description.
	Description() string

	// Execute executes the operation with the given parameters.
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)

	// IsUIVisible returns whether this operation should be visible in the UI.
	IsUIVisible() bool
}
