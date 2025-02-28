package model

import (
	"github.com/HenryOwenz/cloudgate/internal/aws"
)

// Service represents a cloud service
type Service struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Category represents a service category
type Category struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Operation represents a service operation
type Operation struct {
	ID          string
	Name        string
	Description string
}

// ErrMsg represents an error message
type ErrMsg struct {
	Err error
}

// ApprovalsMsg represents a message containing approvals
type ApprovalsMsg struct {
	Approvals []aws.ApprovalAction
	Provider  *aws.Provider
}

// ApprovalResultMsg represents the result of an approval action
type ApprovalResultMsg struct {
	Err error
}

// PipelineStatusMsg represents a message containing pipeline status
type PipelineStatusMsg struct {
	Pipelines []aws.PipelineStatus
	Provider  *aws.Provider
}

// PipelineExecutionMsg represents the result of a pipeline execution
type PipelineExecutionMsg struct {
	Err error
}
