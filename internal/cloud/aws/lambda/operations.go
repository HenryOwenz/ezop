package lambda

import (
	"context"
	"errors"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// Common errors.
var (
	ErrLoadConfig    = errors.New("failed to load AWS config")
	ErrListFunctions = errors.New("failed to list functions")
	ErrGetFunction   = errors.New("failed to get function details")
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

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *FunctionStatusOperation) IsUIVisible() bool {
	return true
}

// GetFunctionStatus returns the status of all Lambda functions.
func (o *FunctionStatusOperation) GetFunctionStatus(ctx context.Context) ([]cloud.FunctionStatus, error) {
	// Create a new AWS SDK client
	client, err := getClient(ctx, o.profile, o.region)
	if err != nil {
		return nil, err
	}

	// List all functions
	functions, err := listFunctions(ctx, client)
	if err != nil {
		return nil, err
	}

	// Convert to cloud.FunctionStatus
	functionStatuses := make([]cloud.FunctionStatus, len(functions))
	for i, function := range functions {
		memory := int32(0)
		if function.MemorySize != nil {
			memory = *function.MemorySize
		}

		timeout := int32(0)
		if function.Timeout != nil {
			timeout = *function.Timeout
		}

		// CodeSize is not a pointer in the AWS Lambda API
		codeSize := function.CodeSize

		// Get architecture (default to x86_64 if not specified)
		architecture := "x86_64"
		if len(function.Architectures) > 0 {
			architecture = string(function.Architectures[0])
		}

		// Get log group if available
		logGroup := ""
		if function.LoggingConfig != nil && function.LoggingConfig.LogGroup != nil {
			logGroup = *function.LoggingConfig.LogGroup
		}

		functionStatuses[i] = cloud.FunctionStatus{
			Name:         aws.ToString(function.FunctionName),
			Runtime:      string(function.Runtime),
			Memory:       memory,
			Timeout:      timeout,
			LastUpdate:   aws.ToString(function.LastModified),
			Role:         aws.ToString(function.Role),
			Handler:      aws.ToString(function.Handler),
			Description:  aws.ToString(function.Description),
			FunctionArn:  aws.ToString(function.FunctionArn),
			CodeSize:     codeSize,
			Version:      aws.ToString(function.Version),
			PackageType:  string(function.PackageType),
			Architecture: architecture,
			LogGroup:     logGroup,
		}
	}

	return functionStatuses, nil
}

// Execute executes the operation with the given parameters.
func (o *FunctionStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetFunctionStatus(ctx)
}

// getClient creates a new Lambda client.
func getClient(ctx context.Context, profile, region string) (*lambda.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLoadConfig, err)
	}

	return lambda.NewFromConfig(cfg), nil
}

// listFunctions returns a list of all Lambda functions.
func listFunctions(ctx context.Context, client *lambda.Client) ([]types.FunctionConfiguration, error) {
	var functions []types.FunctionConfiguration
	var marker *string

	for {
		output, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{
			Marker: marker,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrListFunctions, err)
		}

		functions = append(functions, output.Functions...)

		if output.NextMarker == nil {
			break
		}
		marker = output.NextMarker
	}

	return functions, nil
}
