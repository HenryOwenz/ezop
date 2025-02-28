package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
)

// FetchPipelineStatus fetches pipeline status from AWS
func FetchPipelineStatus(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected service not found")}
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == m.SelectedCategory.Name {
				selectedCategory = category
				break
			}
		}

		if selectedCategory == nil {
			return core.ErrMsg{Err: fmt.Errorf("selected category not found")}
		}

		// Find the operation for pipeline status
		var statusOperation providers.Operation
		for _, operation := range selectedCategory.Operations() {
			if operation.Name() == "Pipeline Status" {
				statusOperation = operation
				break
			}
		}

		if statusOperation == nil {
			return core.ErrMsg{Err: fmt.Errorf("pipeline status operation not found")}
		}

		// Execute the operation
		ctx := context.Background()
		result, err := statusOperation.Execute(ctx, nil)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Convert the result to pipeline status
		// The result could be of different types depending on the implementation
		var pipelines []aws.PipelineStatus

		// Try to convert from []interface{} first
		if pipelinesInterface, ok := result.([]interface{}); ok {
			pipelines = make([]aws.PipelineStatus, len(pipelinesInterface))
			for i, p := range pipelinesInterface {
				if pipeline, ok := p.(map[string]interface{}); ok {
					// Create a new pipeline status
					pipelines[i] = aws.PipelineStatus{
						Name: fmt.Sprintf("%v", pipeline["Name"]),
					}

					// Convert stages if they exist
					if stagesInterface, ok := pipeline["Stages"].([]interface{}); ok {
						stages := make([]aws.StageStatus, len(stagesInterface))
						for j, s := range stagesInterface {
							if stage, ok := s.(map[string]interface{}); ok {
								stages[j] = aws.StageStatus{
									Name:        fmt.Sprintf("%v", stage["Name"]),
									Status:      fmt.Sprintf("%v", stage["Status"]),
									LastUpdated: fmt.Sprintf("%v", stage["LastUpdated"]),
								}
							}
						}
						pipelines[i].Stages = stages
					}
				}
			}
		} else if cloudPipelines, ok := result.([]struct {
			Name   string
			Stages []struct {
				Name        string
				Status      string
				LastUpdated string
			}
		}); ok {
			// Try to convert from a slice of structs
			pipelines = make([]aws.PipelineStatus, len(cloudPipelines))
			for i, p := range cloudPipelines {
				pipelines[i] = aws.PipelineStatus{
					Name: p.Name,
				}

				// Convert stages
				stages := make([]aws.StageStatus, len(p.Stages))
				for j, s := range p.Stages {
					stages[j] = aws.StageStatus{
						Name:        s.Name,
						Status:      s.Status,
						LastUpdated: s.LastUpdated,
					}
				}
				pipelines[i].Stages = stages
			}
		} else {
			// Try to convert using our helper function
			convertedPipelines, err := convertCodePipelineStatus(result)
			if err != nil {
				return core.ErrMsg{Err: fmt.Errorf("unexpected result type: %T - %v", result, err)}
			}
			pipelines = convertedPipelines
		}

		return core.PipelineStatusMsg{
			Pipelines: pipelines,
		}
	}
}

// ExecutePipeline executes a pipeline
func ExecutePipeline(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedPipeline == nil {
			return core.ErrMsg{Err: fmt.Errorf("no pipeline selected")}
		}

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		// Find the CodePipeline service
		var codePipelineService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == "CodePipeline" {
				codePipelineService = service
				break
			}
		}

		if codePipelineService == nil {
			return core.ErrMsg{Err: fmt.Errorf("CodePipeline service not found")}
		}

		// Find the Workflows category
		var workflowsCategory providers.Category
		for _, category := range codePipelineService.Categories() {
			if category.Name() == "Workflows" {
				workflowsCategory = category
				break
			}
		}

		if workflowsCategory == nil {
			return core.ErrMsg{Err: fmt.Errorf("Workflows category not found")}
		}

		// Find the start pipeline operation
		var startPipelineOperation providers.Operation
		for _, operation := range workflowsCategory.Operations() {
			if operation.Name() == "Start Pipeline" {
				startPipelineOperation = operation
				break
			}
		}

		if startPipelineOperation == nil {
			return core.ErrMsg{Err: fmt.Errorf("start pipeline operation not found")}
		}

		// Prepare parameters for the operation
		params := map[string]interface{}{
			"pipeline_name": m.SelectedPipeline.Name,
		}

		// Add commit ID if specified
		if m.ManualCommitID && m.CommitID != "" {
			params["commit_id"] = m.CommitID
		}

		// Execute the operation
		ctx := context.Background()
		_, err = startPipelineOperation.Execute(ctx, params)
		if err != nil {
			return core.PipelineExecutionMsg{Err: err}
		}

		return core.PipelineExecutionMsg{Err: nil}
	}
}

// convertCodePipelineStatus attempts to convert a slice of codepipeline.PipelineStatus to a slice of aws.PipelineStatus
func convertCodePipelineStatus(result interface{}) ([]aws.PipelineStatus, error) {
	// Get the type name as a string
	typeName := fmt.Sprintf("%T", result)
	if !strings.Contains(typeName, "PipelineStatus") {
		return nil, fmt.Errorf("not a PipelineStatus type: %s", typeName)
	}

	// Use reflection to access the slice
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("not a slice type: %s", typeName)
	}

	// Create the output slice
	length := resultValue.Len()
	pipelines := make([]aws.PipelineStatus, length)

	// Process each pipeline
	for i := 0; i < length; i++ {
		pipelineValue := resultValue.Index(i)

		// Extract the Name field
		nameField := pipelineValue.FieldByName("Name")
		if !nameField.IsValid() {
			return nil, fmt.Errorf("Name field not found in PipelineStatus")
		}

		// Convert the name to string
		var name string
		switch nameField.Kind() {
		case reflect.String:
			name = nameField.String()
		case reflect.Ptr:
			if !nameField.IsNil() {
				name = nameField.Elem().String()
			}
		default:
			return nil, fmt.Errorf("Name field is not a string or string pointer")
		}

		// Create the pipeline status
		pipelines[i] = aws.PipelineStatus{
			Name: name,
		}

		// Extract the Stages field
		stagesField := pipelineValue.FieldByName("Stages")
		if !stagesField.IsValid() || stagesField.Kind() != reflect.Slice {
			continue // Skip stages if not found or not a slice
		}

		// Process each stage
		stagesLen := stagesField.Len()
		stages := make([]aws.StageStatus, stagesLen)

		for j := 0; j < stagesLen; j++ {
			stageValue := stagesField.Index(j)

			// Extract stage fields
			nameField := stageValue.FieldByName("Name")
			statusField := stageValue.FieldByName("Status")
			lastUpdatedField := stageValue.FieldByName("LastUpdated")

			// Convert name to string
			var stageName string
			if nameField.IsValid() {
				switch nameField.Kind() {
				case reflect.String:
					stageName = nameField.String()
				case reflect.Ptr:
					if !nameField.IsNil() {
						stageName = nameField.Elem().String()
					}
				}
			}

			// Convert status to string
			var stageStatus string
			if statusField.IsValid() {
				switch statusField.Kind() {
				case reflect.String:
					stageStatus = statusField.String()
				case reflect.Ptr:
					if !statusField.IsNil() {
						stageStatus = statusField.Elem().String()
					}
				}
			}

			// Convert lastUpdated to string
			var lastUpdated string
			if lastUpdatedField.IsValid() {
				switch lastUpdatedField.Kind() {
				case reflect.String:
					lastUpdated = lastUpdatedField.String()
				case reflect.Ptr:
					if !lastUpdatedField.IsNil() {
						lastUpdated = lastUpdatedField.Elem().String()
					}
				}
			}

			// Create the stage status
			stages[j] = aws.StageStatus{
				Name:        stageName,
				Status:      stageStatus,
				LastUpdated: lastUpdated,
			}
		}

		// Assign stages to the pipeline
		pipelines[i].Stages = stages
	}

	return pipelines, nil
}
