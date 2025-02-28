package core

import (
	"github.com/HenryOwenz/cloudgate/internal/aws"
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

	// AWS Configuration
	AwsProfile string
	AwsRegion  string
	Profiles   []string
	Regions    []string

	// Loading state
	IsLoading  bool
	LoadingMsg string

	// AWS Resources
	Provider   *aws.Provider
	Approvals  []aws.ApprovalAction
	Pipelines  []aws.PipelineStatus
	Services   []Service
	Categories []Category
	Operations []Operation

	// Selection state
	SelectedService   *Service
	SelectedCategory  *Category
	SelectedOperation *Operation
	SelectedApproval  *aws.ApprovalAction
	ApproveAction     bool
	Summary           string
	SelectedPipeline  *aws.PipelineStatus

	// Input state
	ManualCommitID  bool
	CommitID        string
	ApprovalComment string
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
		// Initialize empty slices to avoid nil pointer issues
		Profiles:   []string{},
		Regions:    []string{},
		Approvals:  []aws.ApprovalAction{},
		Pipelines:  []aws.PipelineStatus{},
		Services:   []Service{},
		Categories: []Category{},
		Operations: []Operation{},
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	m.Regions = constants.DefaultAWSRegions
	m.Profiles = aws.GetProfiles()
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
	return &newModel
}
