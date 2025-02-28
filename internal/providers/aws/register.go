package aws

import (
	"github.com/HenryOwenz/cloudgate/internal/providers"
)

func init() {
	// Set the CreateAWSProvider function in the providers package
	providers.CreateAWSProvider = func() providers.Provider {
		return New()
	}
}
