package codepipeline

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// Service represents the CodePipeline service.
type Service struct {
	profile    string
	region     string
	categories []cloud.Category
}

// NewService creates a new CodePipeline service.
func NewService(profile, region string) *Service {
	service := &Service{
		profile:    profile,
		region:     region,
		categories: make([]cloud.Category, 0),
	}

	// Register categories
	service.categories = append(service.categories, NewWorkflowsCategory(profile, region))
	service.categories = append(service.categories, NewInternalOperationsCategory(profile, region))

	return service
}

// Name returns the service's name.
func (s *Service) Name() string {
	return "CodePipeline"
}

// Description returns the service's description.
func (s *Service) Description() string {
	return "AWS CodePipeline"
}

// Categories returns all available categories for this service.
func (s *Service) Categories() []cloud.Category {
	return s.categories
}
