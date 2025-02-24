package ui

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
)

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
	currentView View
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
		currentView: ViewProviders,
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case errMsg:
		return m.handleError(msg)
	case approvalsMsg:
		return m.handleApprovals(msg)
	case approvalResultMsg:
		return m.handleApprovalResult(msg)
	case pipelineExecutionMsg:
		return m.handlePipelineExecution(msg)
	case spinner.TickMsg:
		return m.handleSpinnerTick(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case pipelineStatusMsg:
		return m.handlePipelineStatus(msg)
	}
	return m, nil
}

func (m *Model) handleError(msg errMsg) (tea.Model, tea.Cmd) {
	newModel := *m
	newModel.err = msg.err
	newModel.isLoading = false
	return &newModel, nil
}

func (m *Model) handleApprovals(msg approvalsMsg) (tea.Model, tea.Cmd) {
	newModel := *m
	newModel.approvals = msg.approvals
	newModel.provider = msg.provider
	newModel.currentView = ViewApprovals
	newModel.isLoading = false
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handleApprovalResult(msg approvalResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		return m, func() tea.Msg {
			return errMsg{msg.err}
		}
	}
	// First clear loading state
	newModel := *m
	newModel.isLoading = false
	// Then reset approval state and navigate
	newModel.currentView = ViewSelectCategory
	newModel.resetApprovalState()
	// Clear text input
	newModel.resetTextInput()
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handlePipelineExecution(msg pipelineExecutionMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		return m, func() tea.Msg {
			return errMsg{msg.err}
		}
	}
	newModel := *m
	newModel.isLoading = false
	newModel.currentView = ViewSelectCategory
	newModel.selectedPipeline = nil
	newModel.selectedOperation = nil
	newModel.resetTextInput()
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handleSpinnerTick(msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	if m.isLoading {
		var cmd tea.Cmd
		newModel := *m
		newModel.spinner, cmd = m.spinner.Update(msg)
		return &newModel, cmd
	}
	return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m.handleKeyPressWithError(msg)
	}

	if m.isLoading {
		return m.handleKeyPressWhileLoading(msg)
	}

	if m.manualInput || m.currentView == ViewSummary {
		return m.handleKeyPressInTextInput(msg)
	}

	return m.handleKeyPressInNormalMode(msg)
}

func (m *Model) handleKeyPressWithError(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "ctrl+c":
		return m, tea.Quit
	case "-":
		newModel := *m
		newModel.err = nil
		return newModel.navigateBack(), nil
	}
	return m, nil
}

func (m *Model) handleKeyPressWhileLoading(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, m.spinner.Tick
}

func (m *Model) handleKeyPressInTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		newModel := *m
		if m.currentView == ViewSummary && m.selectedApproval != nil {
			// For approval summary, go back to confirmation
			newModel.currentView = ViewConfirmation
			newModel.resetTextInput()
		} else {
			newModel.manualInput = false
		}
		newModel.updateTableForView()
		return &newModel, nil
	case "enter":
		if m.textInput.Value() != "" {
			newModel := *m
			if m.currentView == ViewSummary {
				newModel.summary = m.textInput.Value()
				if m.selectedApproval != nil {
					// For approval summary, move to execution
					newModel.currentView = ViewExecutingAction
					newModel.textInput.Blur()
				} else {
					// For pipeline start summary
					newModel.textInput.Blur()
					newModel.manualInput = false
					newModel.currentView = ViewExecutingAction
				}
				newModel.updateTableForView()
				return &newModel, nil
			} else if m.awsProfile == "" {
				newModel.awsProfile = m.textInput.Value()
			} else {
				newModel.awsRegion = m.textInput.Value()
				newModel.currentView = ViewSelectService
			}
			newModel.resetTextInput()
			newModel.manualInput = false
			newModel.updateTableForView()
			return &newModel, nil
		}
		return m, nil
	default:
		var tiCmd tea.Cmd
		m.textInput, tiCmd = m.textInput.Update(msg)
		return m, tiCmd
	}
}

func (m *Model) handleKeyPressInNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "-", "esc":
		if m.currentView > ViewProviders {
			return m.navigateBack(), nil
		}
	case "tab":
		if m.currentView == ViewAWSConfig || m.currentView == ViewSummary {
			newModel := *m
			newModel.manualInput = !m.manualInput
			if newModel.manualInput {
				newModel.textInput.Focus()
				newModel.textInput.SetValue("")
			} else {
				newModel.textInput.Blur()
			}
			return &newModel, nil
		}
	case "enter":
		return m.handleEnter()
	}

	// Handle table navigation for non-input views
	if !m.manualInput && m.currentView != ViewSummary {
		var tableCmd tea.Cmd
		newModel := *m
		newModel.table, tableCmd = m.table.Update(msg)
		return &newModel, tableCmd
	}

	return m, nil
}

func (m *Model) View() string {
	if m.err != nil {
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Error.Render("Error: "+m.err.Error()),
				"\n",
				m.styles.Help.Render("q: quit • -: back"),
			),
		)
	}

	content := []string{
		m.styles.Title.Render(m.getTitleText()),
		m.styles.Context.Render(m.getContextText()),
		"",
		"",
		"", // Empty line for help text
	}

	// Add loading spinner if needed
	if m.isLoading {
		content[2] = m.spinner.View()
	}

	// For Summary view with approvals, always show text input
	if m.currentView == ViewSummary && m.selectedApproval != nil {
		content[3] = m.textInput.View()
	} else {
		// For other views, follow normal logic
		if !m.manualInput {
			content[3] = m.table.View()
		}
		if m.manualInput {
			content[3] = m.textInput.View()
		}
	}

	// Add help text
	content[4] = m.styles.Help.Render(m.getHelpText())

	// Join all content vertically with consistent spacing
	return m.styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			content...,
		),
	)
}

// handleEnter processes the enter key press based on the current view
func (m *Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewProviders:
		return m.handleProviderSelection()
	case ViewAWSConfig:
		return m.handleAWSConfigSelection()
	case ViewSelectService:
		return m.handleServiceSelection()
	case ViewSelectCategory:
		return m.handleCategorySelection()
	case ViewSelectOperation:
		return m.handleOperationSelection()
	case ViewApprovals:
		return m.handleApprovalSelection()
	case ViewConfirmation:
		return m.handleConfirmationSelection()
	case ViewSummary:
		if !m.manualInput {
			if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
				if selected := m.table.SelectedRow(); len(selected) > 0 {
					newModel := *m
					switch selected[0] {
					case "Latest Commit":
						newModel.currentView = ViewExecutingAction
						newModel.summary = "" // Empty string means use latest commit
						newModel.updateTableForView()
						return &newModel, nil
					case "Manual Input":
						newModel.manualInput = true
						newModel.textInput.Focus()
						newModel.textInput.Placeholder = "Enter commit ID"
						return &newModel, nil
					}
				}
			}
		}
		return m.handleSummaryConfirmation()
	case ViewExecutingAction:
		return m.handleExecutionSelection()
	case ViewPipelineStatus:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			newModel := *m
			for _, pipeline := range m.pipelines {
				if pipeline.Name == selected[0] {
					if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
						newModel.currentView = ViewExecutingAction
						newModel.selectedPipeline = &pipeline
						newModel.updateTableForView()
						return &newModel, nil
					}
					newModel.selectedPipeline = &pipeline
					newModel.currentView = ViewPipelineStages
					newModel.updateTableForView()
					return &newModel, nil
				}
			}
		}
	}
	return m, nil
}

func (m *Model) handleProviderSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "Amazon Web Services" {
			newModel := *m
			newModel.currentView = ViewAWSConfig
			newModel.updateTableForView()
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) handleAWSConfigSelection() (tea.Model, tea.Cmd) {
	if m.manualInput {
		return m.handleManualAWSConfig()
	}
	return m.handleTableAWSConfig()
}

func (m *Model) handleManualAWSConfig() (tea.Model, tea.Cmd) {
	if m.textInput.Value() == "" {
		return m, nil
	}
	newModel := *m
	if m.awsProfile == "" {
		newModel.awsProfile = m.textInput.Value()
	} else {
		newModel.awsRegion = m.textInput.Value()
		newModel.currentView = ViewSelectService
	}
	newModel.resetTextInput()
	newModel.manualInput = false
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handleTableAWSConfig() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		if m.awsProfile == "" {
			newModel.awsProfile = selected[0]
		} else {
			newModel.awsRegion = selected[0]
			newModel.currentView = ViewSelectService
		}
		newModel.updateTableForView()
		return &newModel, nil
	}
	return m, nil
}

func (m *Model) handleServiceSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "CodePipeline" {
			newModel := *m
			newModel.selectedService = &Service{
				Name:        selected[0],
				Description: selected[1],
				Available:   true,
			}
			newModel.currentView = ViewSelectCategory
			newModel.updateTableForView()
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) handleCategorySelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "Workflows" {
			newModel := *m
			newModel.selectedCategory = &Category{
				Name:        selected[0],
				Description: selected[1],
				Available:   true,
			}
			newModel.currentView = ViewSelectOperation
			newModel.updateTableForView()
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) handleOperationSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		switch selected[0] {
		case "Pipeline Approvals":
			newModel.selectedOperation = &Operation{
				Name:        selected[0],
				Description: selected[1],
			}
			newModel.isLoading = true
			return &newModel, tea.Batch(
				newModel.spinner.Tick,
				m.initializeAWS,
			)
		case "Pipeline Status":
			newModel.selectedOperation = &Operation{
				Name:        selected[0],
				Description: selected[1],
			}
			newModel.isLoading = true
			return &newModel, tea.Batch(
				newModel.spinner.Tick,
				m.initializePipelineStatus,
			)
		case "Start Pipeline":
			newModel.selectedOperation = &Operation{
				Name:        selected[0],
				Description: selected[1],
			}
			newModel.isLoading = true
			return &newModel, tea.Batch(
				newModel.spinner.Tick,
				m.initializePipelineStatus,
			)
		}
	}
	return m, nil
}

func (m *Model) handleApprovalSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		for _, approval := range m.approvals {
			if approval.PipelineName == selected[0] {
				newModel.selectedApproval = &approval
				break
			}
		}
		if newModel.selectedApproval != nil {
			newModel.currentView = ViewConfirmation
			newModel.updateTableForView()
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) handleConfirmationSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		switch selected[0] {
		case "Approve":
			newModel.approveAction = true
			newModel.currentView = ViewSummary
			newModel.setTextInputForApproval(true)
			return &newModel, nil
		case "Reject":
			newModel.approveAction = false
			newModel.currentView = ViewSummary
			newModel.setTextInputForApproval(false)
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) handleSummaryConfirmation() (tea.Model, tea.Cmd) {
	if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
		if m.selectedPipeline == nil {
			return m, nil
		}
		newModel := *m
		newModel.currentView = ViewExecutingAction
		newModel.updateTableForView()
		return &newModel, nil
	}

	if m.selectedApproval == nil {
		return m, nil
	}

	newModel := *m
	newModel.currentView = ViewExecutingAction
	newModel.isLoading = true
	newModel.updateTableForView()

	return &newModel, tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			err := m.provider.PutApprovalResult(
				context.Background(),
				*m.selectedApproval,
				m.approveAction,
				m.textInput.Value(),
			)
			if err != nil {
				return errMsg{err}
			}
			return approvalResultMsg{}
		},
	)
}

func (m *Model) handleExecutionSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		switch selected[0] {
		case "Execute":
			newModel := *m
			newModel.isLoading = true
			newModel.resetTextInput()

			if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
				if m.selectedPipeline == nil {
					return m, nil
				}
				return &newModel, tea.Batch(
					newModel.spinner.Tick,
					func() tea.Msg {
						err := m.provider.StartPipelineExecution(
							context.Background(),
							m.selectedPipeline.Name,
							"", // Always use latest commit
						)
						if err != nil {
							return errMsg{err}
						}
						return pipelineExecutionMsg{}
					},
				)
			}

			if m.selectedApproval == nil {
				return m, nil
			}
			return &newModel, tea.Batch(
				newModel.spinner.Tick,
				func() tea.Msg {
					err := m.provider.PutApprovalResult(
						context.Background(),
						*m.selectedApproval,
						m.approveAction,
						m.summary,
					)
					if err != nil {
						return errMsg{err}
					}
					return approvalResultMsg{}
				},
			)
		case "Cancel":
			newModel := *m
			newModel.currentView = ViewSelectCategory
			newModel.resetApprovalState()
			newModel.resetTextInput()
			newModel.updateTableForView()
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) initializeAWS() tea.Msg {
	provider, err := aws.New(context.Background(), m.awsProfile, m.awsRegion)
	if err != nil {
		return errMsg{err}
	}

	approvals, err := provider.GetPendingApprovals(context.Background())
	if err != nil {
		return errMsg{err}
	}

	return approvalsMsg{
		provider:  provider,
		approvals: approvals,
	}
}

func (m *Model) initializePipelineStatus() tea.Msg {
	provider, err := aws.New(context.Background(), m.awsProfile, m.awsRegion)
	if err != nil {
		return errMsg{err}
	}

	pipelines, err := provider.GetPipelineStatus(context.Background())
	if err != nil {
		return errMsg{err}
	}

	return pipelineStatusMsg{
		provider:  provider,
		pipelines: pipelines,
	}
}

// updateTableForView updates the table model based on the current view
func (m *Model) updateTableForView() {
	columns := m.getColumnsForView()
	rows := m.getRowsForView()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)

	t.SetStyles(m.styles.Table)
	m.table = t
}

// getColumnsForView returns the appropriate columns for the current view
func (m *Model) getColumnsForView() []table.Column {
	switch m.currentView {
	case ViewProviders:
		return []table.Column{
			{Title: "Provider", Width: 30},
			{Title: "Description", Width: 50},
		}
	case ViewAWSConfig:
		if m.awsProfile == "" {
			return []table.Column{{Title: "Profile", Width: 30}}
		}
		return []table.Column{{Title: "Region", Width: 30}}
	case ViewSelectService:
		return []table.Column{
			{Title: "Service", Width: 30},
			{Title: "Description", Width: 50},
		}
	case ViewSelectCategory:
		return []table.Column{
			{Title: "Category", Width: 30},
			{Title: "Description", Width: 50},
		}
	case ViewSelectOperation:
		return []table.Column{
			{Title: "Operation", Width: 30},
			{Title: "Description", Width: 50},
		}
	case ViewApprovals:
		return []table.Column{
			{Title: "Pipeline", Width: 40},
			{Title: "Stage", Width: 30},
			{Title: "Action", Width: 20},
		}
	case ViewConfirmation:
		return []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
	case ViewExecutingAction:
		return []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
	case ViewPipelineStatus:
		return []table.Column{
			{Title: "Pipeline", Width: 40},
			{Title: "Description", Width: 50},
		}
	case ViewPipelineStages:
		return []table.Column{
			{Title: "Stage", Width: 30},
			{Title: "Status", Width: 20},
			{Title: "Last Updated", Width: 20},
		}
	case ViewSummary:
		return []table.Column{
			{Title: "Type", Width: 30},
			{Title: "Value", Width: 50},
		}
	default:
		return []table.Column{}
	}
}

// getRowsForView returns the appropriate rows for the current view
func (m *Model) getRowsForView() []table.Row {
	switch m.currentView {
	case ViewProviders:
		return []table.Row{
			{"Amazon Web Services", "AWS Cloud Services"},
			{"Microsoft Azure (Coming Soon)", "Azure Cloud Platform"},
			{"Google Cloud Platform (Coming Soon)", "Google Cloud Services"},
		}
	case ViewAWSConfig:
		if m.awsProfile == "" {
			rows := make([]table.Row, len(m.profiles))
			for i, profile := range m.profiles {
				rows[i] = table.Row{profile}
			}
			return rows
		}
		rows := make([]table.Row, len(m.regions))
		for i, region := range m.regions {
			rows[i] = table.Row{region}
		}
		return rows
	case ViewSelectService:
		return []table.Row{
			{"CodePipeline", "Continuous Delivery Service"},
		}
	case ViewSelectCategory:
		return []table.Row{
			{"Workflows", "Pipeline Workflows and Approvals"},
			{"Operations (Coming Soon)", "Service Operations"},
		}
	case ViewSelectOperation:
		if m.selectedCategory != nil && m.selectedCategory.Name == "Workflows" {
			return []table.Row{
				{"Pipeline Approvals", "Manage Pipeline Approvals"},
				{"Pipeline Status", "View Pipeline Status"},
				{"Start Pipeline", "Trigger Pipeline Execution"},
			}
		}
		return []table.Row{}
	case ViewApprovals:
		rows := make([]table.Row, len(m.approvals))
		for i, approval := range m.approvals {
			rows[i] = table.Row{
				approval.PipelineName,
				approval.StageName,
				approval.ActionName,
			}
		}
		return rows
	case ViewConfirmation:
		return []table.Row{
			{"Approve", "Approve the pipeline stage"},
			{"Reject", "Reject the pipeline stage"},
		}
	case ViewExecutingAction:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			return []table.Row{
				{"Execute", "Start pipeline with latest commit"},
				{"Cancel", "Cancel and return to main menu"},
			}
		}
		action := "approve"
		if !m.approveAction {
			action = "reject"
		}
		return []table.Row{
			{"Execute", fmt.Sprintf("Execute %s action", action)},
			{"Cancel", "Cancel and return to main menu"},
		}
	case ViewPipelineStatus:
		if m.pipelines == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.pipelines))
		for i, pipeline := range m.pipelines {
			rows[i] = table.Row{
				pipeline.Name,
				fmt.Sprintf("%d stages", len(pipeline.Stages)),
			}
		}
		return rows
	case ViewPipelineStages:
		if m.selectedPipeline == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.selectedPipeline.Stages))
		for i, stage := range m.selectedPipeline.Stages {
			rows[i] = table.Row{
				stage.Name,
				stage.Status,
				stage.LastUpdated,
			}
		}
		return rows
	case ViewSummary:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			if m.selectedPipeline == nil {
				return []table.Row{}
			}
			return []table.Row{
				{"Latest Commit", "Use latest commit from source"},
				{"Manual Input", "Enter specific commit ID"},
			}
		}
		// For approval summary, don't show any rows since we're showing text input
		return []table.Row{}
	default:
		return []table.Row{}
	}
}

func (m *Model) navigateBack() *Model {
	newModel := *m
	switch m.currentView {
	case ViewAWSConfig:
		if m.awsProfile != "" {
			// If we're in region selection, just clear region and stay in AWS config
			newModel.awsRegion = ""
			newModel.awsProfile = ""
			// Don't change the view - we'll stay in AWS config to show profiles
		} else {
			// If we're in profile selection, go back to providers
			newModel.currentView = ViewProviders
		}
		newModel.manualInput = false
		newModel.resetTextInput()
	case ViewSelectService:
		newModel.currentView = ViewAWSConfig
		newModel.selectedService = nil
	case ViewSelectCategory:
		newModel.currentView = ViewSelectService
		newModel.selectedCategory = nil
	case ViewSelectOperation:
		newModel.currentView = ViewSelectCategory
		newModel.selectedOperation = nil
	case ViewApprovals:
		newModel.currentView = ViewSelectOperation
		newModel.resetApprovalState()
	case ViewConfirmation:
		newModel.currentView = ViewApprovals
		newModel.selectedApproval = nil
	case ViewSummary:
		newModel.currentView = ViewConfirmation
		newModel.summary = ""
		newModel.resetTextInput()
	case ViewExecutingAction:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			// For pipeline start flow, go back to pipeline selection
			newModel.currentView = ViewPipelineStatus
			newModel.selectedPipeline = nil
		} else {
			// For approval flow, go back to summary
			newModel.currentView = ViewSummary
			// When going back to summary, restore the previous comment and focus
			newModel.textInput.SetValue(m.summary)
			newModel.textInput.Focus()
			if newModel.approveAction {
				newModel.textInput.Placeholder = "Enter approval comment..."
			} else {
				newModel.textInput.Placeholder = "Enter rejection comment..."
			}
		}
	case ViewPipelineStages:
		newModel.currentView = ViewPipelineStatus
		newModel.selectedPipeline = nil
	case ViewPipelineStatus:
		newModel.currentView = ViewSelectOperation
		newModel.pipelines = nil
		newModel.provider = nil
	}
	newModel.updateTableForView()
	return &newModel
}

// Helper functions for common operations
func (m *Model) resetApprovalState() {
	m.approvals = nil
	m.provider = nil
	m.selectedApproval = nil
	m.summary = ""
}

func (m *Model) resetTextInput() {
	m.textInput.SetValue("")
	m.textInput.Blur()
}

func (m *Model) setTextInputForApproval(isApproval bool) {
	m.textInput.Focus()
	if isApproval {
		m.textInput.Placeholder = "Enter approval comment..."
	} else {
		m.textInput.Placeholder = "Enter rejection comment..."
	}
}

// getContextText returns the appropriate context text for the current view
func (m *Model) getContextText() string {
	switch m.currentView {
	case ViewProviders:
		return "A simple tool to manage your cloud resources"
	case ViewAWSConfig:
		if m.awsProfile == "" {
			return "Amazon Web Services"
		}
		return fmt.Sprintf("Profile: %s", m.awsProfile)
	case ViewSelectService:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case ViewSelectCategory:
		if m.selectedService == nil {
			return ""
		}
		return fmt.Sprintf("Service: %s",
			m.selectedService.Name)
	case ViewSelectOperation:
		if m.selectedService == nil || m.selectedCategory == nil {
			return ""
		}
		return fmt.Sprintf("Service: %s\nCategory: %s",
			m.selectedService.Name,
			m.selectedCategory.Name)
	case ViewApprovals:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case ViewConfirmation, ViewSummary:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			if m.selectedPipeline == nil {
				return ""
			}
			return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
				m.awsProfile,
				m.awsRegion,
				m.selectedPipeline.Name)
		}
		if m.selectedApproval == nil {
			return ""
		}
		return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName)
	case ViewExecutingAction:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			if m.selectedPipeline == nil {
				return ""
			}
			return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s\nRevisionID: Latest commit",
				m.awsProfile,
				m.awsRegion,
				m.selectedPipeline.Name)
		}
		if m.selectedApproval == nil {
			return ""
		}
		return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s\nComment: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName,
			m.summary)
	case ViewPipelineStatus:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case ViewPipelineStages:
		if m.selectedPipeline == nil {
			return ""
		}
		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
			m.awsProfile,
			m.awsRegion,
			m.selectedPipeline.Name)
	default:
		return ""
	}
}

// getTitleText returns the appropriate title for the current view
func (m *Model) getTitleText() string {
	switch m.currentView {
	case ViewProviders:
		return "Select Cloud Provider"
	case ViewAWSConfig:
		if m.awsProfile == "" {
			return "Select AWS Profile"
		}
		return "Select AWS Region"
	case ViewSelectService:
		return "Select AWS Service"
	case ViewSelectCategory:
		return "Select Category"
	case ViewSelectOperation:
		return "Select Operation"
	case ViewApprovals:
		return "Pipeline Approvals"
	case ViewConfirmation:
		return "Execute Action"
	case ViewSummary:
		return "Enter Comment"
	case ViewExecutingAction:
		return "Execute Action"
	case ViewPipelineStatus:
		return "Select Pipeline"
	case ViewPipelineStages:
		return "Pipeline Stages"
	default:
		return ""
	}
}

// getHelpText returns the appropriate help text for the current view
func (m *Model) getHelpText() string {
	switch {
	case m.currentView == ViewProviders:
		return "↑/↓: navigate • enter: select • q: quit"
	case m.currentView == ViewAWSConfig && m.manualInput:
		return "enter: confirm • esc: cancel • ctrl+c: quit"
	case m.currentView == ViewAWSConfig:
		return "↑/↓: navigate • enter: select • tab: toggle input • esc: back • q: quit"
	case m.currentView == ViewSummary && m.manualInput:
		return "enter: confirm • esc: cancel • ctrl+c: quit"
	case m.currentView == ViewSummary:
		return "↑/↓: navigate • enter: select • tab: toggle input • esc: back • q: quit"
	default:
		return "↑/↓: navigate • enter: select • esc: back • q: quit"
	}
}

func (m *Model) handlePipelineStatus(msg pipelineStatusMsg) (tea.Model, tea.Cmd) {
	newModel := *m
	newModel.pipelines = msg.pipelines
	newModel.provider = msg.provider
	newModel.currentView = ViewPipelineStatus
	newModel.isLoading = false
	newModel.updateTableForView()
	return &newModel, nil
}
