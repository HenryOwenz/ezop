package model

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/HenryOwenz/ezop/internal/domain"
	"github.com/HenryOwenz/ezop/internal/providers/aws"
	"github.com/HenryOwenz/ezop/internal/ui/styles"
)

// Step represents the current step in the UI workflow
type Step int

const (
	StepSelectProvider Step = iota
	StepProviderConfig
	StepSelectService
	StepServiceOperation
	StepSelectingApproval
	StepConfirmingAction
	StepSummaryInput
	StepExecutingAction
)

// Model represents the application state
type Model struct {
	Profiles    []string
	Regions     []string
	Approvals   []aws.ApprovalAction
	Cursor      int
	Step        Step
	Error       error
	Styles      styles.Styles
	ManualInput bool
	InputBuffer string

	// Provider selection
	SelectedProvider *domain.Provider
	Providers        []domain.Provider

	// Service selection
	Services          []domain.Service
	SelectedService   *domain.Service
	Operations        []domain.Operation
	SelectedOperation *domain.Operation

	// AWS specific
	AWSProfile       string
	AWSRegion        string
	AWSProvider      *aws.Provider
	SelectedApproval *aws.ApprovalAction
	Summary          string
	Action           string // "approve" or "reject"
}

// NewModel creates a new Model with initial state
func NewModel() Model {
	return Model{
		Profiles:    getAWSProfiles(),
		Regions:     []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1"},
		Step:        StepSelectProvider,
		Cursor:      0,
		Styles:      styles.DefaultStyles(),
		ManualInput: false,
		InputBuffer: "",
		Providers:   domain.DefaultProviders,
	}
}

// getAWSProfiles returns a list of AWS profiles from the AWS credentials file
func getAWSProfiles() []string {
	// Read profiles from AWS credentials file
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{"default"}
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

	if len(profiles) == 0 {
		return []string{"default"}
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
