package update

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloudproviders"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

// This file contains updated action handlers that use the cloud package directly
// instead of the providers package.

// InitializeProviders initializes the providers in the registry
func InitializeProviders(m *model.Model) {
	cloudproviders.InitializeProviders(m.Registry)
}

// CreateProvider creates a provider with the given name and configuration
func CreateProvider(m *model.Model, name, profile, region string) (cloud.Provider, error) {
	return cloudproviders.CreateProvider(m.Registry, name, profile, region)
}
