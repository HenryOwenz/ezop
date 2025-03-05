package integration_test

import (
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/providers/testutil"
)

// CreateMockAWSProvider is a helper function that creates a mock AWS provider.
func CreateMockAWSProvider() providers.Provider {
	return testutil.NewMockAWSProvider()
}
