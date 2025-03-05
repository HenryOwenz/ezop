package aws

import (
	cloudAws "github.com/HenryOwenz/cloudgate/internal/cloud/aws"
)

// CreateProvider creates a new AWS provider.
func CreateProvider() *Provider {
	// Create a new cloud provider
	cloudProvider := cloudAws.New()

	// Create a new provider adapter
	return New(cloudProvider)
}
