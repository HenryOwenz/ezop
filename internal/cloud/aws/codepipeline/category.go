package codepipeline

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// WorkflowsCategory represents the CodePipeline workflows category.
type WorkflowsCategory struct {
	profile    string
	region     string
	operations []cloud.Operation
}

// NewWorkflowsCategory creates a new CodePipeline workflows category.
func NewWorkflowsCategory(profile, region string) *WorkflowsCategory {
	category := &WorkflowsCategory{
		profile:    profile,
		region:     region,
		operations: make([]cloud.Operation, 0),
	}

	// Register operations
	category.operations = append(category.operations, NewCloudPipelineStatusOperation(profile, region))
	category.operations = append(category.operations, NewCloudStartPipelineOperation(profile, region))
	category.operations = append(category.operations, NewCloudManualApprovalOperation(profile, region))

	return category
}

// Name returns the category's name.
func (c *WorkflowsCategory) Name() string {
	return "Workflows"
}

// Description returns the category's description.
func (c *WorkflowsCategory) Description() string {
	return "CodePipeline Workflows"
}

// Operations returns all available operations for this category.
func (c *WorkflowsCategory) Operations() []cloud.Operation {
	return c.operations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (c *WorkflowsCategory) IsUIVisible() bool {
	return true
}

// InternalOperationsCategory represents the CodePipeline internal operations category.
type InternalOperationsCategory struct {
	profile    string
	region     string
	operations []cloud.Operation
}

// NewInternalOperationsCategory creates a new CodePipeline internal operations category.
func NewInternalOperationsCategory(profile, region string) *InternalOperationsCategory {
	category := &InternalOperationsCategory{
		profile:    profile,
		region:     region,
		operations: make([]cloud.Operation, 0),
	}

	// Register operations
	category.operations = append(category.operations, NewCloudManualApprovalOperation(profile, region))

	return category
}

// Name returns the category's name.
func (c *InternalOperationsCategory) Name() string {
	return "Internal Operations"
}

// Description returns the category's description.
func (c *InternalOperationsCategory) Description() string {
	return "CodePipeline Internal Operations"
}

// Operations returns all available operations for this category.
func (c *InternalOperationsCategory) Operations() []cloud.Operation {
	return c.operations
}

// IsUIVisible returns whether this category should be visible in the UI.
func (c *InternalOperationsCategory) IsUIVisible() bool {
	return false
}
