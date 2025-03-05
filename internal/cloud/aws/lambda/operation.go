package lambda

import (
	"context"
	"fmt"
	"time"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
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
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(o.profile),
		config.WithRegion(o.region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := lambda.NewFromConfig(cfg)

	// List all functions
	var marker *string
	var functions []cloud.FunctionStatus

	for {
		output, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{
			Marker: marker,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list functions: %w", err)
		}

		// Convert function configurations to FunctionStatus
		for _, fn := range output.Functions {
			lastModified := "Unknown"
			if fn.LastModified != nil {
				// Parse the timestamp
				t, err := time.Parse(time.RFC3339, *fn.LastModified)
				if err == nil {
					lastModified = t.UTC().Format("Jan 02 15:04:05") + " UTC"
				}
			}

			status := cloud.FunctionStatus{
				Name:         aws.ToString(fn.FunctionName),
				Runtime:      string(fn.Runtime),
				Memory:       aws.ToInt32(fn.MemorySize),
				Timeout:      aws.ToInt32(fn.Timeout),
				LastUpdate:   lastModified,
				Role:         aws.ToString(fn.Role),
				Handler:      aws.ToString(fn.Handler),
				Description:  aws.ToString(fn.Description),
				FunctionArn:  aws.ToString(fn.FunctionArn),
				CodeSize:     fn.CodeSize,
				Version:      aws.ToString(fn.Version),
				PackageType:  string(fn.PackageType),
				Architecture: string(fn.Architectures[0]),
			}

			// Get the function's log group
			status.LogGroup = fmt.Sprintf("/aws/lambda/%s", status.Name)

			functions = append(functions, status)
		}

		// Check if there are more functions to fetch
		if output.NextMarker == nil {
			break
		}
		marker = output.NextMarker
	}

	return functions, nil
}

// Execute executes the operation with the given parameters.
func (o *FunctionStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetFunctionStatus(ctx)
}
