package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/HenryOwenz/ciselect/internal/aws"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	stepSelectProvider = iota
	stepProviderConfig
	stepSelectService
	stepServiceOperation
	stepSelectingApproval
	stepConfirmingAction
	stepSummaryInput
	stepExecutingAction
)

type Provider struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

var providers = []Provider{
	{
		ID:          "aws",
		Name:        "Amazon Web Services",
		Description: "AWS Cloud Services",
		Available:   true,
	},
	{
		ID:          "azure",
		Name:        "Microsoft Azure",
		Description: "Azure Cloud Platform (Coming Soon)",
		Available:   false,
	},
	{
		ID:          "gcp",
		Name:        "Google Cloud Platform",
		Description: "Google Cloud Services (Coming Soon)",
		Available:   false,
	},
}

type Service struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

var awsServices = []Service{
	{
		ID:          "codepipeline",
		Name:        "CodePipeline",
		Description: "Continuous Delivery Service",
		Available:   true,
	},
	// Placeholder for future AWS services
}

type Operation struct {
	ID          string
	Name        string
	Description string
}

var codePipelineOperations = []Operation{
	{
		ID:          "manual-approval",
		Name:        "Manual Approval",
		Description: "Manage manual approval actions",
	},
	// Placeholder for future CodePipeline operations
}

type Styles struct {
	Title       lipgloss.Style
	Selected    lipgloss.Style
	Unselected  lipgloss.Style
	Instruction lipgloss.Style
	Error       lipgloss.Style
	Disabled    lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(1, 0),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		Unselected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Instruction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Italic(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		Disabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")),
	}
}

// Model represents the application state
type Model struct {
	profiles    []string
	regions     []string
	approvals   []aws.ApprovalAction
	cursor      int
	step        int
	err         error
	styles      Styles
	manualInput bool
	inputBuffer string

	// Provider selection
	selectedProvider *Provider

	// AWS specific
	awsProfile       string
	awsRegion        string
	selectedService  *Service // UI service selection
	awsOperation     *Operation
	awsClient        *aws.Service // AWS API client
	selectedApproval *aws.ApprovalAction
	summary          string
	action           string // "approve" or "reject"
}

func initialModel(profile, region string) Model {
	m := Model{
		profiles:    getAWSProfiles(),
		regions:     []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1"},
		step:        stepSelectProvider,
		cursor:      0,
		styles:      defaultStyles(),
		manualInput: false,
		inputBuffer: "",
	}

	// If profile/region provided via flags, skip to service selection
	if profile != "" && region != "" {
		m.awsProfile = profile
		m.awsRegion = region
		m.selectedProvider = &providers[0] // AWS
		m.step = stepSelectService
	}

	return m
}

// Update handles messages and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if !m.manualInput && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if !m.manualInput {
				switch m.step {
				case stepSelectProvider:
					if m.cursor < len(providers)-1 {
						m.cursor++
					}
				case stepProviderConfig:
					if m.awsProfile == "" {
						if m.cursor < len(m.profiles)-1 {
							m.cursor++
						}
					} else {
						if m.cursor < len(m.regions)-1 {
							m.cursor++
						}
					}
				case stepSelectService:
					if m.selectedProvider.ID == "aws" && m.cursor < len(awsServices)-1 {
						m.cursor++
					}
				case stepServiceOperation:
					if m.selectedService.ID == "codepipeline" && m.cursor < len(codePipelineOperations)-1 {
						m.cursor++
					}
				case stepSelectingApproval:
					if m.cursor < len(m.approvals)-1 {
						m.cursor++
					}
				case stepConfirmingAction:
					if m.cursor < 2 { // Three options: Approve, Reject, Cancel
						m.cursor++
					}
				case stepExecutingAction:
					if m.cursor < 1 { // Two options: Yes, No
						m.cursor++
					}
				}
			}

		case "tab":
			if m.step == stepProviderConfig {
				m.manualInput = !m.manualInput
				m.inputBuffer = ""
				m.cursor = 0
			}

		case "enter":
			switch m.step {
			case stepSelectProvider:
				provider := providers[m.cursor]
				if !provider.Available {
					return m, nil
				}
				m.selectedProvider = &provider
				m.step = stepProviderConfig
				m.cursor = 0

			case stepProviderConfig:
				if m.selectedProvider.ID == "aws" {
					if m.awsProfile == "" {
						if m.manualInput {
							if m.inputBuffer != "" {
								m.awsProfile = m.inputBuffer
								m.inputBuffer = ""
								m.cursor = 0
							}
						} else {
							if len(m.profiles) > 0 {
								m.awsProfile = m.profiles[m.cursor]
								m.cursor = 0
							}
						}
						return m, nil
					} else if m.awsRegion == "" {
						if m.manualInput {
							if m.inputBuffer != "" {
								m.awsRegion = m.inputBuffer
								m.inputBuffer = ""
								m.manualInput = false
								m.step = stepSelectService
								m.cursor = 0
							}
						} else {
							if len(m.regions) > 0 {
								m.awsRegion = m.regions[m.cursor]
								m.step = stepSelectService
								m.cursor = 0
							}
						}
						return m, nil
					}
				}

			case stepSelectService:
				if m.selectedProvider.ID == "aws" {
					service := awsServices[m.cursor]
					if !service.Available {
						return m, nil
					}
					m.selectedService = &service
					m.step = stepServiceOperation
					m.cursor = 0
				}

			case stepServiceOperation:
				if m.selectedService.ID == "codepipeline" {
					operation := codePipelineOperations[m.cursor]
					m.awsOperation = &operation
					if operation.ID == "manual-approval" {
						// Initialize AWS client and fetch approvals
						return m.initAWS()
					}
				}

			case stepSelectingApproval:
				if len(m.approvals) > 0 {
					m.selectedApproval = &m.approvals[m.cursor]
					m.step = stepConfirmingAction
					m.cursor = 0
				}

			case stepConfirmingAction:
				switch m.cursor {
				case 0: // Approve
					m.action = "approve"
					m.step = stepSummaryInput
				case 1: // Reject
					m.action = "reject"
					m.step = stepSummaryInput
				case 2: // Cancel
					return m, tea.Quit
				}

			case stepSummaryInput:
				if m.summary != "" {
					m.step = stepExecutingAction
					m.cursor = 0
				}

			case stepExecutingAction:
				if m.cursor == 0 { // Yes
					ctx := context.Background()
					var err error
					if m.action == "approve" {
						err = m.awsClient.ApproveAction(ctx,
							m.selectedApproval.PipelineName,
							m.selectedApproval.StageName,
							m.selectedApproval.ActionName,
							m.selectedApproval.Token,
							m.summary)
					} else {
						err = m.awsClient.RejectAction(ctx,
							m.selectedApproval.PipelineName,
							m.selectedApproval.StageName,
							m.selectedApproval.ActionName,
							m.selectedApproval.Token,
							m.summary)
					}
					if err != nil {
						m.err = err
						return m, nil
					}
					return m, tea.Quit
				} else { // No
					return m, tea.Quit
				}
			}

		case "backspace":
			if m.manualInput && len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			} else if m.step == stepSummaryInput && len(m.summary) > 0 {
				m.summary = m.summary[:len(m.summary)-1]
			}

		default:
			if m.manualInput {
				m.inputBuffer += msg.String()
			} else if m.step == stepSummaryInput {
				m.summary += msg.String()
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	var s strings.Builder

	if m.err != nil {
		return m.styles.Error.Render(fmt.Sprintf("Error: %v", m.err))
	}

	switch m.step {
	case stepSelectProvider:
		s.WriteString(m.styles.Title.Render("Select Cloud Provider"))
		s.WriteString("\n")

		for i, provider := range providers {
			text := fmt.Sprintf("%s - %s", provider.Name, provider.Description)

			if m.cursor == i {
				if provider.Available {
					s.WriteString("\n> " + m.styles.Selected.Render(text))
				} else {
					s.WriteString("\n> " + m.styles.Disabled.Render(text))
				}
			} else {
				if provider.Available {
					s.WriteString("\n  " + text)
				} else {
					s.WriteString("\n  " + m.styles.Disabled.Render(text))
				}
			}
		}

	case stepProviderConfig:
		if m.selectedProvider.ID == "aws" {
			if m.awsProfile == "" {
				s.WriteString(m.styles.Title.Render("Select AWS Profile"))
				s.WriteString("\n")

				if m.manualInput {
					s.WriteString(m.styles.Instruction.Render("Enter AWS Profile: "))
					s.WriteString(m.inputBuffer)
				} else {
					for i, profile := range m.profiles {
						if m.cursor == i {
							s.WriteString("\n> " + m.styles.Selected.Render(profile))
						} else {
							s.WriteString("\n  " + profile)
						}
					}
				}
			} else {
				s.WriteString(m.styles.Title.Render("Select AWS Region"))
				s.WriteString("\n")
				s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Profile: %s", m.awsProfile)))
				s.WriteString("\n")

				if m.manualInput {
					s.WriteString(m.styles.Instruction.Render("Enter AWS Region: "))
					s.WriteString(m.inputBuffer)
				} else {
					for i, region := range m.regions {
						if m.cursor == i {
							s.WriteString("\n> " + m.styles.Selected.Render(region))
						} else {
							s.WriteString("\n  " + region)
						}
					}
				}
			}
		}

	case stepSelectService:
		if m.selectedProvider.ID == "aws" {
			s.WriteString(m.styles.Title.Render("Select AWS Service"))
			s.WriteString("\n")
			s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Profile: %s | Region: %s",
				m.awsProfile, m.awsRegion)))

			for i, service := range awsServices {
				text := fmt.Sprintf("%s - %s", service.Name, service.Description)

				if m.cursor == i {
					if service.Available {
						s.WriteString("\n> " + m.styles.Selected.Render(text))
					} else {
						s.WriteString("\n> " + m.styles.Disabled.Render(text))
					}
				} else {
					if service.Available {
						s.WriteString("\n  " + text)
					} else {
						s.WriteString("\n  " + m.styles.Disabled.Render(text))
					}
				}
			}
		}

	case stepServiceOperation:
		if m.selectedService.ID == "codepipeline" {
			s.WriteString(m.styles.Title.Render("Select Operation"))
			s.WriteString("\n")
			s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Service: %s", m.selectedService.Name)))

			for i, operation := range codePipelineOperations {
				text := fmt.Sprintf("%s - %s", operation.Name, operation.Description)

				if m.cursor == i {
					s.WriteString("\n> " + m.styles.Selected.Render(text))
				} else {
					s.WriteString("\n  " + text)
				}
			}
		}

	case stepSelectingApproval:
		s.WriteString(m.styles.Title.Render("Select Approval"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Profile: %s | Region: %s",
			m.awsProfile, m.awsRegion)))

		if len(m.approvals) == 0 {
			s.WriteString("\n\n")
			s.WriteString(m.styles.Instruction.Render("No pending approvals found"))
			s.WriteString("\n\nPress q to quit")
			return s.String()
		}

		for i, approval := range m.approvals {
			text := fmt.Sprintf("%s → %s → %s",
				approval.PipelineName,
				approval.StageName,
				approval.ActionName)

			if m.cursor == i {
				s.WriteString("\n> " + m.styles.Selected.Render(text))
			} else {
				s.WriteString("\n  " + text)
			}
		}

	case stepConfirmingAction:
		s.WriteString(m.styles.Title.Render("Choose Action"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName)))

		options := []string{"Approve", "Reject", "Cancel"}
		for i, option := range options {
			if m.cursor == i {
				s.WriteString("\n> " + m.styles.Selected.Render(option))
			} else {
				s.WriteString("\n  " + option)
			}
		}

	case stepSummaryInput:
		s.WriteString(m.styles.Title.Render("Enter Summary"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Action: %s", m.action)))
		s.WriteString("\n\nSummary: ")
		s.WriteString(m.summary)
		s.WriteString("_")

	case stepExecutingAction:
		s.WriteString(m.styles.Title.Render("Confirm Action"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf(`Pipeline: %s
Stage: %s
Action: %s
Operation: %s
Summary: %s

Are you sure?`,
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName,
			m.action,
			m.summary)))

		options := []string{"Yes", "No"}
		for i, option := range options {
			if m.cursor == i {
				s.WriteString("\n> " + m.styles.Selected.Render(option))
			} else {
				s.WriteString("\n  " + option)
			}
		}
	}

	// Help text
	s.WriteString("\n\n")
	if m.step <= 1 {
		s.WriteString(m.styles.Instruction.Render("↑/↓: Navigate • Enter: Select • Tab: Toggle Input Mode • q: Quit"))
	} else {
		s.WriteString(m.styles.Instruction.Render("↑/↓: Navigate • Enter: Select • q: Quit"))
	}

	return s.String()
}

func getAWSProfiles() []string {
	// Read profiles from AWS credentials file
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	// Try both config and credentials files
	configFiles := []string{
		filepath.Join(home, ".aws", "config"),
		filepath.Join(home, ".aws", "credentials"),
	}

	var profiles []string
	for _, file := range configFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Parse profiles using regex
		re := regexp.MustCompile(`\[(.*?)\]`)
		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			profile := strings.TrimSpace(match[1])
			// Remove "profile " prefix if present (used in config file)
			profile = strings.TrimPrefix(profile, "profile ")
			if profile != "" && !contains(profiles, profile) {
				profiles = append(profiles, profile)
			}
		}
	}

	sort.Strings(profiles)
	return profiles
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   "Interactive mode for managing AWS CodePipeline approvals",
	Long: `Interactive mode provides a user-friendly interface for managing AWS CodePipeline approvals.
It guides you through selecting your AWS profile, region, and approval actions with a beautiful interface.

You can either:
1. Select from a list of available profiles and regions
2. Type your own profile and region names`,
	Run: func(cmd *cobra.Command, args []string) {
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		p := tea.NewProgram(initialModel(profile, region))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	interactiveCmd.Flags().StringP("profile", "p", "", "AWS profile to use (optional)")
	interactiveCmd.Flags().StringP("region", "r", "", "AWS region to use (optional)")
	rootCmd.AddCommand(interactiveCmd)
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) initAWS() (Model, tea.Cmd) {
	ctx := context.Background()
	client, err := aws.NewService(ctx, m.awsProfile, m.awsRegion)
	if err != nil {
		m.err = err
		return m, nil
	}
	m.awsClient = client

	approvals, err := client.ListPendingApprovals(ctx)
	if err != nil {
		m.err = err
		return m, nil
	}
	m.approvals = approvals
	m.step = stepSelectingApproval
	return m, nil
}
