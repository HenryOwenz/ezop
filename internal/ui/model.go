package ui

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/ezop/v2/internal/aws"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
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

// Model represents the application state
type Model struct {
	table       table.Model
	width       int
	height      int
	styles      Styles
	currentView View
	awsProfile  string
	awsRegion   string
	profiles    []string
	regions     []string
	manualInput bool
	inputBuffer string
	err         error
	approvals   []aws.ApprovalAction
	provider    *aws.Provider

	// Loading state
	isLoading  bool
	loadingMsg string
	spinner    spinner.Model

	// Service selection state
	services        []Service
	selectedService *Service

	// Category selection state
	categories       []Category
	selectedCategory *Category

	// Operation selection state
	operations        []Operation
	selectedOperation *Operation

	// Approval state
	selectedApproval *aws.ApprovalAction
	approveAction    bool
	summary          string
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

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true)

	m := Model{
		currentView: ViewProviders,
		profiles:    aws.GetProfiles(),
		regions: []string{
			"us-east-1", "us-east-2", "us-west-1", "us-west-2",
			"eu-west-1", "eu-west-2", "eu-central-1",
			"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
		},
		styles:  DefaultStyles(),
		spinner: s,
	}
	m.updateTableForView()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
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
				return errMsg{msg.err}
			}
		}
		return m, tea.Quit

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

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "-":
			if m.currentView > ViewProviders {
				newModel := m.navigateBack()
				return newModel, nil
			}
		case "tab":
			if m.currentView == ViewAWSConfig {
				newModel := m
				newModel.manualInput = !m.manualInput
				newModel.inputBuffer = ""
				return newModel, nil
			}
		case "enter":
			return m.handleEnter()
		case "backspace":
			if m.manualInput {
				if len(m.inputBuffer) > 0 {
					newModel := m
					newModel.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
					return newModel, nil
				}
			} else if m.currentView == ViewSummary {
				if len(m.summary) > 0 {
					newModel := m
					newModel.summary = m.summary[:len(m.summary)-1]
					return newModel, nil
				}
			}
		default:
			if m.manualInput {
				newModel := m
				newModel.inputBuffer += msg.String()
				return newModel, nil
			} else if m.currentView == ViewSummary {
				newModel := m
				newModel.summary += msg.String()
				return newModel, nil
			}
		}

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

	switch m.currentView {
	case ViewProviders:
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Select Cloud Provider"),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • q: quit"),
			),
		)
	case ViewAWSConfig:
		var content string
		if m.awsProfile == "" {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Select AWS Profile"),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • tab: toggle input • -: back • q: quit"),
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Select AWS Region"),
				m.styles.Context.Render("Profile: "+m.awsProfile),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • tab: toggle input • -: back • q: quit"),
			)
		}

		if m.manualInput {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				content,
				"\n",
				"Enter value: "+m.inputBuffer+"_",
			)
		}

		return m.styles.App.Render(content)
	case ViewSelectService:
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Select AWS Service"),
				m.styles.Context.Render("Profile: "+m.awsProfile+" | Region: "+m.awsRegion),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • -: back • q: quit"),
			),
		)
	case ViewSelectCategory:
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Select Category"),
				m.styles.Context.Render(fmt.Sprintf("Service: %s", m.selectedService.Name)),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • -: back • q: quit"),
			),
		)
	case ViewSelectOperation:
		var content []string
		content = append(content, m.styles.Title.Render("Select Operation"))
		if m.isLoading {
			content = append(content, m.spinner.View())
		} else {
			content = append(content, "")
		}
		content = append(content, m.styles.Context.Render(fmt.Sprintf("Service: %s | Category: %s",
			m.selectedService.Name, m.selectedCategory.Name)))
		content = append(content, "\n")
		content = append(content, m.table.View())
		content = append(content, "\n")
		content = append(content, m.styles.Help.Render("↑/↓: navigate • enter: select • -: back • q: quit"))

		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				content...,
			),
		)
	case ViewApprovals:
		var content []string
		content = append(content, m.styles.Title.Render("Pipeline Approvals"))
		if m.isLoading {
			content = append(content, m.spinner.View())
		} else {
			content = append(content, "")
		}
		content = append(content, m.styles.Context.Render("Profile: "+m.awsProfile+" | Region: "+m.awsRegion))
		content = append(content, "\n")
		content = append(content, m.table.View())
		content = append(content, "\n")
		content = append(content, m.styles.Help.Render("↑/↓: navigate • enter: select • -: back • q: quit"))

		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				content...,
			),
		)
	case ViewConfirmation:
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Execute Action"),
				m.styles.Context.Render(
					"Pipeline: "+m.selectedApproval.PipelineName+
						" | Stage: "+m.selectedApproval.StageName+
						" | Action: "+m.selectedApproval.ActionName,
				),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • -: back • q: quit"),
			),
		)
	case ViewSummary:
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Enter Summary"),
				m.styles.Context.Render(
					"Pipeline: "+m.selectedApproval.PipelineName+
						" | Stage: "+m.selectedApproval.StageName+
						" | Action: "+m.selectedApproval.ActionName,
				),
				"\n",
				"Summary: "+m.summary+"_",
				"\n",
				m.styles.Help.Render("enter: confirm • -: back • q: quit"),
			),
		)
	case ViewExecutingAction:
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Title.Render("Execute Action"),
				m.styles.Context.Render(
					"Pipeline: "+m.selectedApproval.PipelineName+
						" | Stage: "+m.selectedApproval.StageName+
						" | Action: "+m.selectedApproval.ActionName+
						" | Summary: "+m.summary,
				),
				"\n",
				m.table.View(),
				"\n",
				m.styles.Help.Render("↑/↓: navigate • enter: select • -: back • q: quit"),
			),
		)
	default:
		return "Loading..."
	}
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
			if m.inputBuffer != "" {
				newModel := m
				if m.awsProfile == "" {
					newModel.awsProfile = m.inputBuffer
				} else {
					newModel.awsRegion = m.inputBuffer
					newModel.currentView = ViewSelectService
				}
				newModel.inputBuffer = ""
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
				return newModel, nil
			case "Reject":
				newModel.approveAction = false
				newModel.currentView = ViewSummary
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
				return newModel, tea.Batch(
					newModel.spinner.Tick,
					func() tea.Msg {
						err := m.provider.HandleApproval(context.Background(), m.selectedApproval, m.approveAction, m.summary)
						return approvalResultMsg{err: err}
					},
				)
			case "Cancel":
				return m, tea.Quit
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

type errMsg struct{ err error }
type approvalsMsg struct {
	provider  *aws.Provider
	approvals []aws.ApprovalAction
}
type approvalResultMsg struct{ err error }

func (m *Model) updateTableForView() {
	var columns []table.Column
	var rows []table.Row

	switch m.currentView {
	case ViewProviders:
		columns = []table.Column{
			{Title: "Provider", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"Amazon Web Services", "AWS Cloud Services"},
			{"Microsoft Azure (Coming Soon)", "Azure Cloud Platform"},
			{"Google Cloud Platform (Coming Soon)", "Google Cloud Services"},
		}
	case ViewAWSConfig:
		if m.awsProfile == "" {
			columns = []table.Column{
				{Title: "Profile", Width: 30},
			}
			for _, profile := range m.profiles {
				rows = append(rows, table.Row{profile})
			}
		} else {
			columns = []table.Column{
				{Title: "Region", Width: 30},
			}
			for _, region := range m.regions {
				rows = append(rows, table.Row{region})
			}
		}
	case ViewSelectService:
		columns = []table.Column{
			{Title: "Service", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"CodePipeline", "Continuous Delivery Service"},
			{"CodeBuild (Coming Soon)", "Build Service"},
			{"CodeDeploy (Coming Soon)", "Deployment Service"},
		}
	case ViewSelectCategory:
		columns = []table.Column{
			{Title: "Category", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"Workflows", "Pipeline Workflows and Approvals"},
			{"Operations (Coming Soon)", "Service Operations"},
		}
	case ViewSelectOperation:
		columns = []table.Column{
			{Title: "Operation", Width: 30},
			{Title: "Description", Width: 50},
		}
		if m.selectedCategory != nil && m.selectedCategory.Name == "Workflows" {
			rows = []table.Row{
				{"Pipeline Approvals", "Manage Pipeline Approvals"},
				{"Pipeline Status (Coming Soon)", "View Pipeline Status"},
			}
		}
	case ViewApprovals:
		columns = []table.Column{
			{Title: "Pipeline", Width: 40},
			{Title: "Stage", Width: 30},
			{Title: "Action", Width: 20},
		}
		for _, approval := range m.approvals {
			rows = append(rows, table.Row{
				approval.PipelineName,
				approval.StageName,
				approval.ActionName,
			})
		}
	case ViewConfirmation:
		columns = []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"Approve", "Approve the pipeline stage"},
			{"Reject", "Reject the pipeline stage"},
		}
	case ViewExecutingAction:
		columns = []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
		action := "approve"
		if !m.approveAction {
			action = "reject"
		}
		rows = []table.Row{
			{"Execute", fmt.Sprintf("Execute %s action", action)},
			{"Cancel", "Cancel and return to main menu"},
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(m.styles.Table)
	m.table = t
}

func (m Model) navigateBack() Model {
	newModel := m
	switch m.currentView {
	case ViewAWSConfig:
		if m.awsRegion != "" {
			newModel.awsRegion = ""
		} else if m.awsProfile != "" {
			newModel.awsProfile = ""
		} else {
			newModel.currentView = ViewProviders
		}
		newModel.manualInput = false
		newModel.inputBuffer = ""
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
		newModel.approvals = nil
		newModel.provider = nil
	case ViewConfirmation:
		newModel.currentView = ViewApprovals
		newModel.selectedApproval = nil
	case ViewSummary:
		newModel.currentView = ViewConfirmation
		newModel.summary = ""
	case ViewExecutingAction:
		newModel.currentView = ViewSummary
	}
	newModel.updateTableForView()
	return newModel
}
