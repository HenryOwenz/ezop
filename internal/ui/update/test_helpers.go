package update

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// InitializeTestRegistry creates a registry with the given provider for testing
func InitializeTestRegistry(provider cloud.Provider) *cloud.ProviderRegistry {
	registry := cloud.NewProviderRegistry()
	registry.Register(provider)
	return registry
}
