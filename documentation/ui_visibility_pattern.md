# UI Visibility Pattern

## Overview

This document describes the UI visibility pattern implemented in CloudGate to control which components (categories and operations) are visible in the user interface. This pattern allows for a clear separation between components that should be exposed to users and internal components that should only be accessible programmatically.

## Problem Statement

As CloudGate scales to support hundreds of cloud services and operations, we need a consistent way to:

1. Hide internal operations that should not be directly accessible by users
2. Maintain programmatic access to these operations for internal use
3. Ensure a clean and focused user interface
4. Support a flexible architecture that can adapt to different visibility requirements

## Solution

We implemented a comprehensive visibility control pattern that works at both the category and operation levels:

1. Added an `IsUIVisible()` method to both category and operation interfaces
2. Created internal categories for operations that should not appear in the UI
3. Implemented filtering in the UI layer to respect visibility settings

## Implementation Details

### Interface Changes

We extended both the cloud and providers interfaces to include the `IsUIVisible()` method:

```go
// In internal/cloud/types.go
type Category interface {
    // ... existing methods ...
    IsUIVisible() bool
}

type Operation interface {
    // ... existing methods ...
    IsUIVisible() bool
}

// In internal/providers/interfaces.go
type Category interface {
    // ... existing methods ...
    IsUIVisible() bool
}

type Operation interface {
    // ... existing methods ...
    IsUIVisible() bool
}
```

### Adapter Implementation

The adapter pattern ensures that visibility settings are passed through from the cloud layer to the providers layer:

```go
// In internal/providers/adapter.go
func (a *CloudCategoryAdapter) IsUIVisible() bool {
    return a.category.IsUIVisible()
}

func (a *CloudOperationAdapter) IsUIVisible() bool {
    return a.operation.IsUIVisible()
}
```

### Internal Categories

For operations that should not be visible in the UI but need to be accessible programmatically, we created internal categories:

```go
// In internal/cloud/aws/codepipeline/workflows.go and internal/providers/aws/codepipeline/workflows.go
type InternalOperationsCategory struct {
    profile    string
    region     string
    operations []cloud.Operation // or []providers.Operation
}

func (c *InternalOperationsCategory) IsUIVisible() bool {
    return false
}
```

### UI Filtering

The UI layer respects the visibility settings by filtering out components that should not be visible:

```go
// Filter out internal categories
var visibleCategories []providers.Category
for _, category := range categories {
    if category.IsUIVisible() {
        visibleCategories = append(visibleCategories, category)
    }
}

// Filter out internal operations
var visibleOperations []providers.Operation
for _, operation := range operations {
    if operation.IsUIVisible() {
        visibleOperations = append(visibleOperations, operation)
    }
}
```

## Usage Examples

### Regular Operations (Visible in UI)

```go
func (o *PipelineApprovalsOperation) IsUIVisible() bool {
    return true
}
```

### Internal Operations (Hidden from UI)

```go
func (o *ApprovalOperation) IsUIVisible() bool {
    return false
}
```

### Registering Internal Operations

```go
// Create an internal category
internalCategory := NewInternalOperationsCategory(profile, region)

// Register internal operations
internalCategory.operations = append(internalCategory.operations, NewApprovalOperation(profile, region))

// Add the internal category to the service
service.categories = append(service.categories, internalCategory)
```

## Best Practices

1. **Default to Visible**: Unless there's a specific reason to hide an operation, make it visible by default.
2. **Document Internal Operations**: Clearly document internal operations and their purpose.
3. **Consistent Naming**: Use consistent naming for internal categories (e.g., "InternalOperations").
4. **Separation of Concerns**: Keep UI visibility logic separate from business logic.
5. **Testing**: Test both visible and invisible components to ensure they work correctly.

## Benefits

1. **Clean UI**: Users only see operations that are relevant to them.
2. **Flexible Architecture**: The same operation can be used both internally and externally.
3. **Maintainable Code**: Clear separation between UI and internal operations.
4. **Scalable**: This pattern will work well as we add more cloud providers, services, and operations.

## Future Enhancements

1. **Role-Based Visibility**: Extend the pattern to support different visibility settings based on user roles.
2. **Conditional Visibility**: Allow operations to be visible only under certain conditions.
3. **Visibility Groups**: Group operations by visibility to simplify management.

## Conclusion

The UI visibility pattern provides a clean and consistent way to control which components are visible in the CloudGate user interface. By implementing this pattern at both the category and operation levels, we ensure that users have a focused experience while maintaining the flexibility to use operations programmatically when needed. 