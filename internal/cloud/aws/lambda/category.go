package lambda

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// WorkflowsCategory represents the Lambda workflows category.
type WorkflowsCategory struct {
	profile    string
	region     string
	operations []cloud.Operation
}

// NewWorkflowsCategory creates a new Lambda workflows category.
func NewWorkflowsCategory(profile, region string) *WorkflowsCategory {
	category := &WorkflowsCategory{
		profile:    profile,
		region:     region,
		operations: make([]cloud.Operation, 0),
	}

	// Register operations
	category.operations = append(category.operations, NewFunctionStatusOperation(profile, region))

	return category
}

// Name returns the category's name.
func (c *WorkflowsCategory) Name() string {
	return "Workflows"
}

// Description returns the category's description.
func (c *WorkflowsCategory) Description() string {
	return "Lambda Function Workflows"
}

// Operations returns all available operations for this category.
func (c *WorkflowsCategory) Operations() []cloud.Operation {
	return c.operations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (c *WorkflowsCategory) IsUIVisible() bool {
	return true
}
