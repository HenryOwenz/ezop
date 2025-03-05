package model

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/providers"
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

// ApprovalAction is an alias for cloud.ApprovalAction
type ApprovalAction = cloud.ApprovalAction

// StageStatus is an alias for cloud.StageStatus
type StageStatus = cloud.StageStatus

// PipelineStatus is an alias for cloud.PipelineStatus
type PipelineStatus = cloud.PipelineStatus

// FunctionStatus is an alias for cloud.FunctionStatus
type FunctionStatus = cloud.FunctionStatus

// ApprovalsMsg represents a message containing approvals
type ApprovalsMsg struct {
	Approvals []ApprovalAction
	Provider  providers.Provider
}

// ApprovalResultMsg represents the result of an approval action
type ApprovalResultMsg struct {
	Err error
}

// PipelineStatusMsg represents a message containing pipeline status
type PipelineStatusMsg struct {
	Pipelines []PipelineStatus
	Provider  providers.Provider
}

// PipelineExecutionMsg represents the result of a pipeline execution
type PipelineExecutionMsg struct {
	Err error
}

// FunctionStatusMsg represents a message containing function status
type FunctionStatusMsg struct {
	Functions []FunctionStatus
	Provider  providers.Provider
}
