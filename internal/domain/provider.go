package domain

import "context"

// Provider represents a cloud service provider
type Provider struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Service represents a cloud service within a provider
type Service struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Category represents a group of related operations within a service
type Category struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Operation represents an action that can be performed on a service
type Operation struct {
	ID          string
	Name        string
	Description string
}

// CloudProvider defines the interface that all cloud providers must implement
type CloudProvider interface {
	// GetServices returns the list of available services for this provider
	GetServices() []Service

	// GetOperations returns the list of available operations for a service
	GetOperations(serviceID string) []Operation

	// ExecuteOperation executes an operation on a service
	ExecuteOperation(ctx context.Context, serviceID, operationID string, params map[string]interface{}) error
}

// ProviderRegistry holds the list of available providers
type ProviderRegistry struct {
	Providers []Provider
}

// DefaultProviders returns the list of supported cloud providers
var DefaultProviders = []Provider{
	{
		ID:          "aws",
		Name:        "Amazon Web Services",
		Description: "AWS Cloud Services",
		Available:   true,
	},
	{
		ID:          "azure",
		Name:        "Microsoft Azure",
		Description: "Azure Cloud Platform",
		Available:   false,
	},
	{
		ID:          "gcp",
		Name:        "Google Cloud Platform",
		Description: "Google Cloud Services",
		Available:   false,
	},
}
