package core

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
	ErrMsg       struct{ Err error }
	ApprovalsMsg struct {
		Provider  *aws.Provider
		Approvals []aws.ApprovalAction
	}
	PipelineStatusMsg struct {
		Provider  *aws.Provider
		Pipelines []aws.PipelineStatus
	}
	ApprovalResultMsg    struct{ Err error }
	PipelineExecutionMsg struct{ Err error }
)
