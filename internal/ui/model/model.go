package model

import (
	"sort"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	// UI Components
	Table     table.Model
	TextInput textinput.Model
	Spinner   spinner.Model
	Styles    styles.Styles

	// Window dimensions
	Width  int
	Height int

	// View state
	CurrentView constants.View
	ManualInput bool
	Err         error
	Error       string // Error message
	Success     string // Success message

	// Loading state
	IsLoading  bool
	LoadingMsg string

	// Provider Registry
	Registry *providers.ProviderRegistry

	// Provider state
	ProviderState ProviderState

	// Input state
	InputState InputState

	// Legacy fields for backward compatibility
	// These will be gradually migrated to the new structure
	AwsProfile        string
	AwsRegion         string
	Profiles          []string
	Regions           []string
	Provider          providers.Provider
	Approvals         []providers.ApprovalAction
	Pipelines         []providers.PipelineStatus
	Services          []Service
	Categories        []Category
	Operations        []Operation
	SelectedService   *Service
	SelectedCategory  *Category
	SelectedOperation *Operation
	SelectedApproval  *providers.ApprovalAction
	ApproveAction     bool
	Summary           string
	SelectedPipeline  *providers.PipelineStatus
	ManualCommitID    bool
	CommitID          string
	ApprovalComment   string
}

// ProviderState represents the state of the selected provider, service, category, and operation
type ProviderState struct {
	// Selected provider
	ProviderName string

	// Provider configuration
	Config map[string]string // Generic configuration (e.g., "profile", "region")

	// Available configuration options
	ConfigOptions map[string][]string // e.g., "profile" -> ["default", "dev", "prod"]

	// Current configuration key being set
	CurrentConfigKey string

	// Authentication state
	AuthState AuthenticationState

	// Selected service, category, and operation
	SelectedService   *ServiceInfo
	SelectedCategory  *CategoryInfo
	SelectedOperation *OperationInfo

	// Provider-specific state (stored as generic interface{})
	ProviderSpecificState map[string]interface{}
}

// AuthenticationState represents the authentication state for different providers
type AuthenticationState struct {
	// Current authentication method
	Method string

	// Authentication configuration
	AuthConfig map[string]string

	// Available authentication methods
	AvailableMethods []string

	// Current authentication config key being set
	CurrentAuthConfigKey string

	// Authentication status
	IsAuthenticated bool

	// Error message if authentication failed
	AuthError string
}

// ServiceInfo represents information about a service
type ServiceInfo struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// CategoryInfo represents information about a category
type CategoryInfo struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// OperationInfo represents information about an operation
type OperationInfo struct {
	ID          string
	Name        string
	Description string
}

// InputState represents the state of user input
type InputState struct {
	// Generic input fields
	TextValues map[string]string // e.g., "comment" -> "This is a comment"
	BoolValues map[string]bool   // e.g., "approve" -> true

	// Operation-specific state
	OperationState map[string]interface{} // Operation-specific state
}

// New creates and initializes a new Model
func New() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(constants.ColorWarning)).Italic(true)

	ti := textinput.New()
	ti.Placeholder = constants.MsgEnterComment
	ti.CharLimit = constants.TextInputCharLimit
	ti.Width = constants.TextInputWidth

	t := table.New(
		table.WithHeight(constants.TableHeight),
		table.WithFocused(true),
	)
	t.SetStyles(styles.DefaultStyles().Table)

	m := &Model{
		Spinner:     s,
		TextInput:   ti,
		Table:       t,
		CurrentView: constants.ViewProviders,
		Styles:      styles.DefaultStyles(),
		Registry:    providers.NewProviderRegistry(),

		// Initialize new state structures
		ProviderState: ProviderState{
			Config:                make(map[string]string),
			ConfigOptions:         make(map[string][]string),
			ProviderSpecificState: make(map[string]interface{}),
			AuthState: AuthenticationState{
				AuthConfig:       make(map[string]string),
				AvailableMethods: []string{},
			},
		},
		InputState: InputState{
			TextValues:     make(map[string]string),
			BoolValues:     make(map[string]bool),
			OperationState: make(map[string]interface{}),
		},

		// Initialize legacy fields
		Profiles:   []string{},
		Regions:    []string{},
		Approvals:  []providers.ApprovalAction{},
		Pipelines:  []providers.PipelineStatus{},
		Services:   []Service{},
		Categories: []Category{},
		Operations: []Operation{},
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	m.Regions = constants.DefaultAWSRegions

	// Initialize the AWS provider to get profiles
	if providers.CreateAWSProvider != nil {
		awsProvider := providers.CreateAWSProvider()
		if awsProvider != nil {
			profiles, err := awsProvider.GetProfiles()
			if err == nil && len(profiles) > 0 {
				// Sort the profiles alphabetically
				sort.Strings(profiles)
				m.Profiles = profiles
			} else {
				// Fallback to default profile if there's an error
				m.Profiles = []string{"default"}
			}
		}
	}

	return m.Spinner.Tick
}

// ResetApprovalState resets the approval state
func (m *Model) ResetApprovalState() {
	m.Approvals = nil
	m.Provider = nil
	m.SelectedApproval = nil
	m.Summary = ""
}

// ResetTextInput resets the text input
func (m *Model) ResetTextInput() {
	m.TextInput.SetValue("")
	m.TextInput.Blur()
}

// SetTextInputForApproval configures the text input for approval
func (m *Model) SetTextInputForApproval(isApproval bool) {
	m.TextInput.Focus()
	if isApproval {
		m.TextInput.Placeholder = constants.MsgEnterApprovalComment
	} else {
		m.TextInput.Placeholder = constants.MsgEnterRejectionComment
	}
}

// Clone creates a deep copy of the model
func (m *Model) Clone() *Model {
	newModel := *m

	// Deep copy maps in ProviderState
	newModel.ProviderState.Config = make(map[string]string)
	for k, v := range m.ProviderState.Config {
		newModel.ProviderState.Config[k] = v
	}

	newModel.ProviderState.ConfigOptions = make(map[string][]string)
	for k, v := range m.ProviderState.ConfigOptions {
		newOptions := make([]string, len(v))
		copy(newOptions, v)
		newModel.ProviderState.ConfigOptions[k] = newOptions
	}

	newModel.ProviderState.ProviderSpecificState = make(map[string]interface{})
	for k, v := range m.ProviderState.ProviderSpecificState {
		newModel.ProviderState.ProviderSpecificState[k] = v
	}

	newModel.ProviderState.AuthState.AuthConfig = make(map[string]string)
	for k, v := range m.ProviderState.AuthState.AuthConfig {
		newModel.ProviderState.AuthState.AuthConfig[k] = v
	}

	newModel.ProviderState.AuthState.AvailableMethods = make([]string, len(m.ProviderState.AuthState.AvailableMethods))
	copy(newModel.ProviderState.AuthState.AvailableMethods, m.ProviderState.AuthState.AvailableMethods)

	// Deep copy maps in InputState
	newModel.InputState.TextValues = make(map[string]string)
	for k, v := range m.InputState.TextValues {
		newModel.InputState.TextValues[k] = v
	}

	newModel.InputState.BoolValues = make(map[string]bool)
	for k, v := range m.InputState.BoolValues {
		newModel.InputState.BoolValues[k] = v
	}

	newModel.InputState.OperationState = make(map[string]interface{})
	for k, v := range m.InputState.OperationState {
		newModel.InputState.OperationState[k] = v
	}

	return &newModel
}

// Helper methods for working with the model
func (m *Model) SetProviderConfig(key, value string) {
	m.ProviderState.Config[key] = value
}

func (m *Model) GetProviderConfig(key string) string {
	return m.ProviderState.Config[key]
}

func (m *Model) SetAuthConfig(key, value string) {
	m.ProviderState.AuthState.AuthConfig[key] = value
}

func (m *Model) GetAuthConfig(key string) string {
	return m.ProviderState.AuthState.AuthConfig[key]
}

func (m *Model) SetInputText(key, value string) {
	m.InputState.TextValues[key] = value
}

func (m *Model) GetInputText(key string) string {
	return m.InputState.TextValues[key]
}

func (m *Model) SetInputBool(key string, value bool) {
	m.InputState.BoolValues[key] = value
}

func (m *Model) GetInputBool(key string) bool {
	return m.InputState.BoolValues[key]
}

// Backward compatibility methods

// GetAwsProfile returns the AWS profile from the provider config
func (m *Model) GetAwsProfile() string {
	// First check the new structure
	profile := m.GetProviderConfig("profile")
	if profile != "" {
		return profile
	}
	// Fall back to legacy field
	return m.AwsProfile
}

// SetAwsProfile sets the AWS profile in the provider config
func (m *Model) SetAwsProfile(profile string) {
	m.SetProviderConfig("profile", profile)
	// Also set in legacy field for backward compatibility
	m.AwsProfile = profile
}

// GetAwsRegion returns the AWS region from the provider config
func (m *Model) GetAwsRegion() string {
	// First check the new structure
	region := m.GetProviderConfig("region")
	if region != "" {
		return region
	}
	// Fall back to legacy field
	return m.AwsRegion
}

// SetAwsRegion sets the AWS region in the provider config
func (m *Model) SetAwsRegion(region string) {
	m.SetProviderConfig("region", region)
	// Also set in legacy field for backward compatibility
	m.AwsRegion = region
}

// GetApprovalComment returns the approval comment from the input state
func (m *Model) GetApprovalComment() string {
	// First check the new structure
	comment := m.GetInputText("approval-comment")
	if comment != "" {
		return comment
	}
	// Fall back to legacy field
	return m.ApprovalComment
}

// SetApprovalComment sets the approval comment in the input state
func (m *Model) SetApprovalComment(comment string) {
	m.SetInputText("approval-comment", comment)
	// Also set in legacy field for backward compatibility
	m.ApprovalComment = comment
}

// GetApproveAction returns the approve action from the input state
func (m *Model) GetApproveAction() bool {
	// First check the new structure
	if _, ok := m.InputState.BoolValues["approve-action"]; ok {
		return m.GetInputBool("approve-action")
	}
	// Fall back to legacy field
	return m.ApproveAction
}

// SetApproveAction sets the approve action in the input state
func (m *Model) SetApproveAction(approve bool) {
	m.SetInputBool("approve-action", approve)
	// Also set in legacy field for backward compatibility
	m.ApproveAction = approve
}

// GetCommitID returns the commit ID from the input state
func (m *Model) GetCommitID() string {
	// First check the new structure
	commitID := m.GetInputText("commit-id")
	if commitID != "" {
		return commitID
	}
	// Fall back to legacy field
	return m.CommitID
}

// SetCommitID sets the commit ID in the input state
func (m *Model) SetCommitID(commitID string) {
	m.SetInputText("commit-id", commitID)
	// Also set in legacy field for backward compatibility
	m.CommitID = commitID
}

// GetManualCommitID returns whether to use manual commit ID
func (m *Model) GetManualCommitID() bool {
	// First check the new structure
	if _, ok := m.InputState.BoolValues["manual-commit-id"]; ok {
		return m.GetInputBool("manual-commit-id")
	}
	// Fall back to legacy field
	return m.ManualCommitID
}

// SetManualCommitID sets whether to use manual commit ID
func (m *Model) SetManualCommitID(manual bool) {
	m.SetInputBool("manual-commit-id", manual)
	// Also set in legacy field for backward compatibility
	m.ManualCommitID = manual
}

// GetSelectedApproval returns the selected approval from the provider-specific state
func (m *Model) GetSelectedApproval() *providers.ApprovalAction {
	// First check the new structure
	if approval, ok := m.ProviderState.ProviderSpecificState["selected-approval"]; ok {
		if typedApproval, ok := approval.(*providers.ApprovalAction); ok {
			return typedApproval
		}
	}
	// Fall back to legacy field
	return m.SelectedApproval
}

// SetSelectedApproval sets the selected approval in the provider-specific state
func (m *Model) SetSelectedApproval(approval *providers.ApprovalAction) {
	m.ProviderState.ProviderSpecificState["selected-approval"] = approval
	// Also set in legacy field for backward compatibility
	m.SelectedApproval = approval
}

// GetSelectedPipeline returns the selected pipeline from the provider-specific state
func (m *Model) GetSelectedPipeline() *providers.PipelineStatus {
	// First check the new structure
	if pipeline, ok := m.ProviderState.ProviderSpecificState["selected-pipeline"]; ok {
		if typedPipeline, ok := pipeline.(*providers.PipelineStatus); ok {
			return typedPipeline
		}
	}
	// Fall back to legacy field
	return m.SelectedPipeline
}

// SetSelectedPipeline sets the selected pipeline in the provider-specific state
func (m *Model) SetSelectedPipeline(pipeline *providers.PipelineStatus) {
	m.ProviderState.ProviderSpecificState["selected-pipeline"] = pipeline
	// Also set in legacy field for backward compatibility
	m.SelectedPipeline = pipeline
}

// GetApprovals returns the approvals from the provider-specific state
func (m *Model) GetApprovals() []providers.ApprovalAction {
	// First check the new structure
	if approvals, ok := m.ProviderState.ProviderSpecificState["approvals"]; ok {
		if typedApprovals, ok := approvals.([]providers.ApprovalAction); ok {
			return typedApprovals
		}
	}
	// Fall back to legacy field
	return m.Approvals
}

// SetApprovals sets the approvals in the provider-specific state
func (m *Model) SetApprovals(approvals []providers.ApprovalAction) {
	m.ProviderState.ProviderSpecificState["approvals"] = approvals
	// Also set in legacy field for backward compatibility
	m.Approvals = approvals
}

// GetPipelines returns the pipelines from the provider-specific state
func (m *Model) GetPipelines() []providers.PipelineStatus {
	// First check the new structure
	if pipelines, ok := m.ProviderState.ProviderSpecificState["pipelines"]; ok {
		if typedPipelines, ok := pipelines.([]providers.PipelineStatus); ok {
			return typedPipelines
		}
	}
	// Fall back to legacy field
	return m.Pipelines
}

// SetPipelines sets the pipelines in the provider-specific state
func (m *Model) SetPipelines(pipelines []providers.PipelineStatus) {
	m.ProviderState.ProviderSpecificState["pipelines"] = pipelines
	// Also set in legacy field for backward compatibility
	m.Pipelines = pipelines
}
