package ui

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/ezop/v2/internal/aws"
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
	approvalResultMsg struct{ err error }
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
func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DD6B20", Dark: "#ED8936"}).Italic(true)

	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.CharLimit = 156
	ti.Width = 30

	m := Model{
		currentView: ViewProviders,
		profiles:    aws.GetProfiles(),
		regions: []string{
			"us-east-1", "us-east-2", "us-west-1", "us-west-2",
			"eu-west-1", "eu-west-2", "eu-central-1",
			"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
		},
		styles:    DefaultStyles(),
		spinner:   s,
		textInput: ti,
	}
	m.updateTableForView()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case errMsg:
		newModel := m
		newModel.err = msg.err
		newModel.isLoading = false
		newModel.loadingMsg = ""
		return newModel, nil

	case approvalsMsg:
		newModel := m
		newModel.approvals = msg.approvals
		newModel.provider = msg.provider
		newModel.currentView = ViewApprovals
		newModel.isLoading = false
		newModel.loadingMsg = ""
		newModel.updateTableForView()
		return newModel, nil

	case approvalResultMsg:
		if msg.err != nil {
			return m, func() tea.Msg {
				return msg.err // This will automatically be wrapped in errMsg by the caller
			}
		}
		// First clear loading state
		newModel := m
		newModel.isLoading = false
		newModel.loadingMsg = ""
		// Then reset approval state and navigate
		newModel.currentView = ViewSelectCategory
		newModel.resetApprovalState()
		// Clear text input
		newModel.resetTextInput()
		newModel.updateTableForView()
		return newModel, nil

	case spinner.TickMsg:
		if m.isLoading {
			var cmd tea.Cmd
			newModel := m
			newModel.spinner, cmd = m.spinner.Update(msg)
			return newModel, cmd
		}
		return m, nil

	case tea.KeyMsg:
		if m.err != nil {
			switch msg.String() {
			case "esc", "q", "ctrl+c":
				return m, tea.Quit
			case "-":
				newModel := m
				newModel.err = nil
				newModel = newModel.navigateBack()
				return newModel, nil
			default:
				return m, nil
			}
		}

		// If we're loading, only handle quit
		if m.isLoading {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				return m, m.spinner.Tick
			}
		}

		// Handle text input mode first
		if m.manualInput || m.currentView == ViewSummary {
			// Only allow ctrl+c to quit in text input mode
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

			// Allow escape to cancel text input
			if msg.String() == "esc" {
				newModel := m
				if m.currentView == ViewSummary {
					newModel = newModel.navigateBack()
					return newModel, nil
				}
				newModel.manualInput = false
				newModel.resetTextInput()
				return newModel, nil
			}

			// Handle enter key specially in text input mode
			if msg.String() == "enter" {
				if m.textInput.Value() != "" {
					newModel := m
					if m.currentView == ViewSummary {
						newModel.summary = m.textInput.Value()
						newModel.textInput.Blur()
						newModel.currentView = ViewExecutingAction
						newModel.updateTableForView()
						return newModel, nil
					} else if m.awsProfile == "" {
						newModel.awsProfile = m.textInput.Value()
					} else {
						newModel.awsRegion = m.textInput.Value()
						newModel.currentView = ViewSelectService
					}
					newModel.resetTextInput()
					newModel.manualInput = false
					newModel.updateTableForView()
					return newModel, nil
				}
				return m, nil
			}

			// Handle all other keys as text input
			var tiCmd tea.Cmd
			m.textInput, tiCmd = m.textInput.Update(msg)
			return m, tiCmd
		}

		// Handle navigation and other commands when not in text input mode
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "-", "esc":
			if m.currentView > ViewProviders {
				newModel := m.navigateBack()
				return newModel, nil
			}
		case "tab":
			if m.currentView == ViewAWSConfig {
				newModel := m
				newModel.manualInput = !m.manualInput
				if newModel.manualInput {
					newModel.textInput.Focus()
					newModel.textInput.SetValue("")
				} else {
					newModel.textInput.Blur()
				}
				return newModel, nil
			}
		case "enter":
			return m.handleEnter()
		}

		// Handle table navigation for non-input views
		if !m.manualInput && m.currentView != ViewSummary {
			var tableCmd tea.Cmd
			newModel := m
			newModel.table, tableCmd = m.table.Update(msg)
			return newModel, tableCmd
		}
	}

	return m, cmd
}

func (m Model) View() string {
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

	// Replace content with table view for list-based views
	if !m.manualInput && m.currentView != ViewSummary {
		content[3] = m.table.View()
	}

	// Add input field for manual input views
	if m.manualInput || m.currentView == ViewSummary {
		content[3] = m.textInput.View()
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

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewProviders:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			if selected[0] == "Amazon Web Services" {
				newModel := m
				newModel.currentView = ViewAWSConfig
				newModel.updateTableForView()
				return newModel, nil
			}
		}
	case ViewAWSConfig:
		if m.manualInput {
			if m.textInput.Value() != "" {
				newModel := m
				if m.awsProfile == "" {
					newModel.awsProfile = m.textInput.Value()
				} else {
					newModel.awsRegion = m.textInput.Value()
					newModel.currentView = ViewSelectService
				}
				newModel.resetTextInput()
				newModel.manualInput = false
				newModel.updateTableForView()
				return newModel, nil
			}
		} else if selected := m.table.SelectedRow(); len(selected) > 0 {
			newModel := m
			if m.awsProfile == "" {
				newModel.awsProfile = selected[0]
			} else {
				newModel.awsRegion = selected[0]
				newModel.currentView = ViewSelectService
			}
			newModel.updateTableForView()
			return newModel, nil
		}
	case ViewSelectService:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			if selected[0] == "CodePipeline" {
				newModel := m
				newModel.selectedService = &Service{
					Name:        selected[0],
					Description: selected[1],
					Available:   true,
				}
				newModel.currentView = ViewSelectCategory
				newModel.updateTableForView()
				return newModel, nil
			}
		}
	case ViewSelectCategory:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			if selected[0] == "Workflows" {
				newModel := m
				newModel.selectedCategory = &Category{
					Name:        selected[0],
					Description: selected[1],
					Available:   true,
				}
				newModel.currentView = ViewSelectOperation
				newModel.updateTableForView()
				return newModel, nil
			}
		}
	case ViewSelectOperation:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			if selected[0] == "Pipeline Approvals" {
				newModel := m
				newModel.selectedOperation = &Operation{
					Name:        selected[0],
					Description: selected[1],
				}
				newModel.isLoading = true
				newModel.loadingMsg = "Fetching pipeline approvals..."
				return newModel, tea.Batch(
					newModel.spinner.Tick,
					m.initializeAWS,
				)
			}
		}
	case ViewApprovals:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			newModel := m
			for _, approval := range m.approvals {
				if approval.PipelineName == selected[0] {
					newModel.selectedApproval = &approval
					break
				}
			}
			if newModel.selectedApproval != nil {
				newModel.currentView = ViewConfirmation
				newModel.updateTableForView()
				return newModel, nil
			}
		}
	case ViewConfirmation:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			newModel := m
			switch selected[0] {
			case "Approve":
				newModel.approveAction = true
				newModel.currentView = ViewSummary
				newModel.setTextInputForApproval(true)
				return newModel, nil
			case "Reject":
				newModel.approveAction = false
				newModel.currentView = ViewSummary
				newModel.setTextInputForApproval(false)
				return newModel, nil
			}
		}
	case ViewSummary:
		if m.summary != "" {
			newModel := m
			newModel.currentView = ViewExecutingAction
			newModel.updateTableForView()
			return newModel, nil
		}
	case ViewExecutingAction:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			switch selected[0] {
			case "Execute":
				newModel := m
				newModel.isLoading = true
				newModel.loadingMsg = fmt.Sprintf("%sing pipeline...", m.approveAction)
				// Clear text input
				newModel.resetTextInput()
				return newModel, tea.Batch(
					newModel.spinner.Tick,
					func() tea.Msg {
						err := m.provider.HandleApproval(context.Background(), m.selectedApproval, m.approveAction, m.summary)
						return approvalResultMsg{err: err}
					},
				)
			case "Cancel":
				newModel := m
				newModel.currentView = ViewSelectCategory
				newModel.resetApprovalState()
				// Clear text input
				newModel.resetTextInput()
				newModel.updateTableForView()
				return newModel, nil
			}
		}
	}
	return m, nil
}

func (m Model) initializeAWS() tea.Msg {
	provider, err := aws.New(m.awsProfile, m.awsRegion)
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
				{"Pipeline Status (Coming Soon)", "View Pipeline Status"},
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
		action := "approve"
		if !m.approveAction {
			action = "reject"
		}
		return []table.Row{
			{"Execute", fmt.Sprintf("Execute %s action", action)},
			{"Cancel", "Cancel and return to main menu"},
		}
	default:
		return []table.Row{}
	}
}

func (m Model) navigateBack() Model {
	newModel := m
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
	newModel.updateTableForView()
	return newModel
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
		return fmt.Sprintf("Service: %s",
			m.selectedService.Name)
	case ViewSelectOperation:
		return fmt.Sprintf("Service: %s\nCategory: %s",
			m.selectedService.Name,
			m.selectedCategory.Name)
	case ViewApprovals:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case ViewConfirmation, ViewSummary:
		return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName)
	case ViewExecutingAction:
		return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s\nComment: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName,
			m.summary)
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
	case m.currentView == ViewSummary:
		return "enter: confirm • esc: back • ctrl+c: quit"
	default:
		return "↑/↓: navigate • enter: select • esc: back • q: quit"
	}
}
