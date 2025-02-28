package aws

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/codepipeline"
)

// Common errors
var (
	ErrAWSConfigNotFound = errors.New("aws configuration not found")
	ErrProfileNotFound   = errors.New("aws profile not found")
	ErrRegionNotFound    = errors.New("aws region not found")
)

// Provider represents the AWS cloud provider.
type Provider struct {
	profile  string
	region   string
	services []cloud.Service
}

// New creates a new AWS provider.
func New() *Provider {
	return &Provider{
		services: make([]cloud.Service, 0),
	}
}

// Name returns the provider's name.
func (p *Provider) Name() string {
	return "AWS"
}

// Description returns the provider's description.
func (p *Provider) Description() string {
	return "Amazon Web Services"
}

// Services returns all available services for this provider.
func (p *Provider) Services() []cloud.Service {
	return p.services
}

// GetProfiles returns all available AWS profiles from the user's home directory.
func (p *Provider) GetProfiles() ([]string, error) {
	profiles := make([]string, 0)

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Parse AWS config file
	configProfiles, err := parseAWSConfigFile(filepath.Join(homeDir, ".aws", "config"))
	if err == nil {
		profiles = append(profiles, configProfiles...)
	}

	// Parse AWS credentials file
	credProfiles, err := parseAWSConfigFile(filepath.Join(homeDir, ".aws", "credentials"))
	if err == nil {
		profiles = append(profiles, credProfiles...)
	}

	// Remove duplicates
	uniqueProfiles := make(map[string]bool)
	for _, profile := range profiles {
		uniqueProfiles[profile] = true
	}

	// Convert back to slice
	profiles = make([]string, 0, len(uniqueProfiles))
	for profile := range uniqueProfiles {
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// LoadConfig loads AWS configuration based on the provided profile and region.
func (p *Provider) LoadConfig(profile, region string) error {
	if profile == "" {
		return ErrProfileNotFound
	}

	if region == "" {
		return ErrRegionNotFound
	}

	// Set profile and region
	p.profile = profile
	p.region = region

	// Initialize services
	p.services = make([]cloud.Service, 0)

	// Add CodePipeline service
	p.services = append(p.services, codepipeline.NewService(profile, region))

	return nil
}

// parseAWSConfigFile parses an AWS config file and returns the profile names.
func parseAWSConfigFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	profiles := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Extract profile name
			profile := line[1 : len(line)-1]
			// Always use TrimPrefix regardless of whether the prefix exists
			profile = strings.TrimPrefix(profile, "profile ")
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}
