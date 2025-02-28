# UI Operation Interfaces

## Overview

This document describes the UI operation interfaces implemented in CloudGate to provide a clearer separation between UI operations and their provider-specific implementations. This architecture addresses terminology issues and creates a more flexible structure for supporting multiple cloud providers.

## Problem Statement

The previous architecture had several issues:

1. **Terminology Confusion**: The term "approval" was too generic and didn't clearly indicate its AWS CodePipeline context.
2. **Nested Operations**: The structure had operations within operations, leading to linguistic confusion.
3. **Direct Provider Coupling**: UI operations were directly tied to specific technical implementations.
4. **Limited Flexibility**: Adding new providers required significant changes to the UI code.

## Solution

We implemented a new architecture with dedicated interfaces for UI operations:

1. **Base UIOperation Interface**: Defines common methods for all UI operations.
2. **Provider-Specific Operation Interfaces**: Extend the base interface with provider-specific methods.
3. **Provider Interface Extensions**: Added methods to get specific operation implementations.
4. **Operation Wrappers**: Created adapter types that implement the operation interfaces.

## Interface Hierarchy

```
UIOperation
├── CodePipelineManualApprovalOperation
├── PipelineStatusOperation
└── StartPipelineOperation
```

### UIOperation Interface

The base interface for all UI operations:

```go
type UIOperation interface {
    // Name returns the operation's name
    Name() string

    // Description returns the operation's description
    Description() string

    // IsUIVisible returns whether this operation should be visible in the UI
    IsUIVisible() bool
}
```

### Provider-Specific Operation Interfaces

Specialized interfaces for specific operations:

```go
type CodePipelineManualApprovalOperation interface {
    UIOperation
    
    // GetPendingApprovals returns all pending manual approval actions
    GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error)
    
    // ApproveAction approves or rejects an approval action
    ApproveAction(ctx context.Context, action ApprovalAction, approved bool, comment string) error
}

type PipelineStatusOperation interface {
    UIOperation
    
    // GetPipelineStatus returns the status of all pipelines
    GetPipelineStatus(ctx context.Context) ([]PipelineStatus, error)
}

type StartPipelineOperation interface {
    UIOperation
    
    // StartPipelineExecution starts a pipeline execution
    StartPipelineExecution(ctx context.Context, pipelineName string, commitID string) error
}
```

### Provider Interface Extensions

New methods added to the Provider interface:

```go
type Provider interface {
    // ... existing methods ...
    
    // GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
    GetCodePipelineManualApprovalOperation() (CodePipelineManualApprovalOperation, error)
    
    // GetPipelineStatusOperation returns the pipeline status operation
    GetPipelineStatusOperation() (PipelineStatusOperation, error)
    
    // GetStartPipelineOperation returns the start pipeline operation
    GetStartPipelineOperation() (StartPipelineOperation, error)
}
```

## Implementation Details

### Operation Wrappers

Each provider implements wrapper types that adapt the provider's functionality to the operation interfaces:

```go
type codePipelineManualApprovalOperation struct {
    provider *Provider
}

func (o *codePipelineManualApprovalOperation) Name() string {
    return "Pipeline Approvals"
}

func (o *codePipelineManualApprovalOperation) Description() string {
    return "Manage Pipeline Approvals"
}

func (o *codePipelineManualApprovalOperation) IsUIVisible() bool {
    return true
}

func (o *codePipelineManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
    return o.provider.GetApprovals(ctx)
}

func (o *codePipelineManualApprovalOperation) ApproveAction(ctx context.Context, action providers.ApprovalAction, approved bool, comment string) error {
    return o.provider.ApproveAction(ctx, action, approved, comment)
}
```

### UI Code Updates

The UI code has been updated to use the new operation interfaces:

```go
// Get the CodePipelineManualApprovalOperation from the provider
approvalOperation, err := provider.GetCodePipelineManualApprovalOperation()
if err != nil {
    return model.ErrMsg{Err: err}
}

// Get approvals using the operation
ctx := context.Background()
approvals, err := approvalOperation.GetPendingApprovals(ctx)
if err != nil {
    return model.ErrMsg{Err: err}
}
```

## Benefits

1. **Clear Terminology**: The new interfaces use more specific names like `CodePipelineManualApprovalOperation`.
2. **Separation of Concerns**: UI operations are clearly separated from their implementations.
3. **Provider-Specific Implementations**: Each provider can implement operations in its own way.
4. **Flexibility**: Adding new providers or operations is easier with the new architecture.
5. **Maintainability**: The code is more maintainable with clearer interfaces and responsibilities.

## Future Enhancements

1. **Additional Providers**: Implement similar interfaces for Azure and GCP.
2. **Operation Discovery**: Add mechanisms for dynamically discovering available operations.
3. **Operation Parameters**: Standardize parameter passing for operations.
4. **Operation Results**: Create standardized result types for operations.

## Conclusion

The new UI operation interfaces provide a more flexible and maintainable architecture for CloudGate. By clearly separating UI operations from their implementations, we've addressed terminology issues and created a foundation for supporting multiple cloud providers. 