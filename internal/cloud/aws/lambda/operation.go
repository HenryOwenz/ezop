package lambda

import (
	"context"
	"fmt"
)

// FunctionStatusOperation represents an operation to view Lambda function status.
type FunctionStatusOperation struct {
	profile string
	region  string
}

// NewFunctionStatusOperation creates a new function status operation.
func NewFunctionStatusOperation(profile, region string) *FunctionStatusOperation {
	return &FunctionStatusOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *FunctionStatusOperation) Name() string {
	return "Function Status"
}

// Description returns the operation's description.
func (o *FunctionStatusOperation) Description() string {
	return "View Lambda Function Status"
}

// Execute executes the operation with the given parameters.
func (o *FunctionStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	functions, err := GetFunctionStatus(ctx, o.profile, o.region)
	if err != nil {
		return nil, fmt.Errorf("failed to get function status: %w", err)
	}

	return functions, nil
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *FunctionStatusOperation) IsUIVisible() bool {
	return true
}
