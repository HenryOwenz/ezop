package integration_test

import (
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/providers/testutil"
)

func init() {
	// Set up the mock AWS provider
	providers.CreateAWSProvider = testutil.NewMockAWSProvider
}
