package aws

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Common errors
var (
	ErrConfigNotFound      = fmt.Errorf("AWS config file not found")
	ErrCredentialsNotFound = fmt.Errorf("AWS credentials file not found")
	ErrNoProfiles          = fmt.Errorf("no AWS profiles found")
	ErrNoRegions           = fmt.Errorf("no AWS regions found")
)

// getAWSProfiles returns all available AWS profiles from the user's home directory.
func getAWSProfiles() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Get profiles from config file
	configProfiles, err := parseAWSConfigFile(filepath.Join(homeDir, ".aws", "config"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Get profiles from credentials file
	credentialsProfiles, err := parseAWSConfigFile(filepath.Join(homeDir, ".aws", "credentials"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Combine profiles and remove duplicates
	profileMap := make(map[string]bool)
	for _, profile := range configProfiles {
		profileMap[profile] = true
	}
	for _, profile := range credentialsProfiles {
		profileMap[profile] = true
	}

	// Convert map to slice
	profiles := make([]string, 0, len(profileMap))
	for profile := range profileMap {
		profiles = append(profiles, profile)
	}

	if len(profiles) == 0 {
		return nil, ErrNoProfiles
	}

	return profiles, nil
}

// parseAWSConfigFile parses an AWS config file and returns all profile names.
func parseAWSConfigFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles []string
	profileRegex := regexp.MustCompile(`^\[(?:profile\s+)?([^\]]+)\]$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := profileRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			profiles = append(profiles, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read AWS config file: %w", err)
	}

	return profiles, nil
}
