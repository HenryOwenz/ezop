package integration

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// CreateMockAWSProvider creates a mock AWS provider for testing
func CreateMockAWSProvider() cloud.Provider {
	return &MockAWSProvider{}
}
