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

// Service represents the Lambda service.
type Service struct {
	profile    string
	region     string
	categories []cloud.Category
}

// NewService creates a new Lambda service.
func NewService(profile, region string) *Service {
	service := &Service{
		profile:    profile,
		region:     region,
		categories: make([]cloud.Category, 0),
	}

	// Register categories
	service.categories = append(service.categories, NewWorkflowsCategory(profile, region))

	return service
}

// Name returns the service's name.
func (s *Service) Name() string {
	return "Lambda"
}

// Description returns the service's description.
func (s *Service) Description() string {
	return "Serverless Compute Service"
}

// Categories returns all available categories for this service.
func (s *Service) Categories() []cloud.Category {
	return s.categories
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

// FunctionStatus represents the status of a Lambda function
type FunctionStatus struct {
	Name         string
	Runtime      string
	Memory       int32
	Timeout      int32
	LastUpdate   string
	Role         string
	Handler      string
	Description  string
	FunctionArn  string
	CodeSize     int64
	Version      string
	PackageType  string
	Architecture string
	LogGroup     string
}

// GetFunctionStatus returns the status of all Lambda functions
func GetFunctionStatus(ctx context.Context, profile, region string) ([]FunctionStatus, error) {
	client, err := getClient(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	functions, err := listFunctions(ctx, client)
	if err != nil {
		return nil, err
	}

	functionStatuses := make([]FunctionStatus, len(functions))
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

		functionStatuses[i] = FunctionStatus{
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
