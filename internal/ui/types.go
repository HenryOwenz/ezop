package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/aws"
)

// Service represents an AWS service
type Service struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Category represents a group of operations
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

// Message types for internal communication
type (
	errMsg       struct{ err error }
	approvalsMsg struct {
		provider  *aws.Provider
		approvals []aws.ApprovalAction
	}
	pipelineStatusMsg struct {
		provider  *aws.Provider
		pipelines []aws.PipelineStatus
	}
	approvalResultMsg    struct{ err error }
	pipelineExecutionMsg struct{ err error }
) 