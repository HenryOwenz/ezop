package aws

// This file contains type definitions that are used by the AWS provider.
// The implementation of the provider methods has been moved to use the adapter pattern,
// delegating to the cloud layer through UI operation interfaces.

// ApprovalAction represents a pending approval in a pipeline.
// This type is kept for backward compatibility but should be phased out
// in favor of using providers.ApprovalAction directly.
type ApprovalAction struct {
	PipelineName string
	StageName    string
	ActionName   string
	Token        string
}

// PipelineStatus represents the status of a pipeline and its stages.
// This type is kept for backward compatibility but should be phased out
// in favor of using providers.PipelineStatus directly.
type PipelineStatus struct {
	Name   string
	Stages []StageStatus
}

// StageStatus represents the status of a pipeline stage.
// This type is kept for backward compatibility but should be phased out
// in favor of using providers.StageStatus directly.
type StageStatus struct {
	Name        string
	Status      string
	LastUpdated string
}
