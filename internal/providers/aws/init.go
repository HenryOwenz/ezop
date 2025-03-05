package aws

import (
	"github.com/HenryOwenz/cloudgate/internal/providers"
)

// init registers the AWS provider with the global registry
func init() {
	// Replace the placeholder AWS provider in the registry with the real implementation
	providers.RegisterAWSProvider(CreateProvider())
}
