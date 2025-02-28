package constants

// View represents different screens in the application
type View int

const (
	ViewProviders View = iota
	ViewAWSConfig
	ViewSelectService
	ViewSelectCategory
	ViewSelectOperation
	ViewApprovals
	ViewPipelineStatus
	ViewPipelineStages
	ViewConfirmation
	ViewSummary
	ViewExecutingAction
	ViewError
	ViewSuccess
	ViewHelp

	// New view states for the updated model
	ViewAuthMethodSelect
	ViewAuthConfig
	ViewProviderConfig
	ViewAuthError
)
