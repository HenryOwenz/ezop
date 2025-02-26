package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	// UI Components
	table     table.Model
	textInput textinput.Model
	spinner   spinner.Model
	styles    Styles

	// Window dimensions
	width  int
	height int

	// View state
	currentView constants.View
	manualInput bool
	err         error

	// AWS Configuration
	awsProfile string
	awsRegion  string
	profiles   []string
	regions    []string

	// Loading state
	isLoading  bool
	loadingMsg string

	// AWS Resources
	provider   *aws.Provider
	approvals  []aws.ApprovalAction
	pipelines  []aws.PipelineStatus
	services   []Service
	categories []Category
	operations []Operation

	// Selection state
	selectedService   *Service
	selectedCategory  *Category
	selectedOperation *Operation
	selectedApproval  *aws.ApprovalAction
	approveAction     bool
	summary           string
	selectedPipeline  *aws.PipelineStatus
}

// New creates and initializes a new Model
func New() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DD6B20", Dark: "#ED8936"}).Italic(true)

	ti := textinput.New()
	ti.Placeholder = "Enter comment..."
	ti.CharLimit = 100
	ti.Width = 50

	t := table.New(
		table.WithHeight(6),
		table.WithFocused(true),
	)
	t.SetStyles(DefaultStyles().Table)

	m := &Model{
		spinner:     s,
		textInput:   ti,
		table:       t,
		currentView: constants.ViewProviders,
		styles:      DefaultStyles(),
	}

	m.updateTableForView()
	return m
}

func (m *Model) Init() tea.Cmd {
	m.regions = []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"eu-west-1", "eu-west-2", "eu-central-1",
		"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
	}
	m.profiles = aws.GetProfiles()
	return m.spinner.Tick
}
