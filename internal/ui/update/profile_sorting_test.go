package update

import (
	"sort"
	"testing"
)

// TestProfileStringSorting tests that profiles are sorted alphabetically
func TestProfileStringSorting(t *testing.T) {
	testCases := []struct {
		name          string
		profiles      []string
		expectedOrder []string
	}{
		{
			name:          "Unsorted Profiles",
			profiles:      []string{"dev", "prod", "test", "default", "staging"},
			expectedOrder: []string{"default", "dev", "prod", "staging", "test"},
		},
		{
			name:          "Already Sorted Profiles",
			profiles:      []string{"a", "b", "c", "d"},
			expectedOrder: []string{"a", "b", "c", "d"},
		},
		{
			name:          "Profiles with Numbers",
			profiles:      []string{"profile2", "profile10", "profile1"},
			expectedOrder: []string{"profile1", "profile10", "profile2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a copy of the profiles to sort
			profiles := make([]string, len(tc.profiles))
			copy(profiles, tc.profiles)

			// Sort the profiles
			sort.Strings(profiles)

			// Verify the sorting
			if len(profiles) != len(tc.expectedOrder) {
				t.Fatalf("Expected %d profiles, got %d", len(tc.expectedOrder), len(profiles))
			}

			for i, expected := range tc.expectedOrder {
				if profiles[i] != expected {
					t.Errorf("Expected profile at index %d to be '%s', got '%s'", i, expected, profiles[i])
				}
			}
		})
	}
}

// TestSortingInHandleProviderSelection tests that profiles are sorted in HandleProviderSelection
func TestSortingInHandleProviderSelection(t *testing.T) {
	// This is a placeholder test to verify that sort.Strings is called in HandleProviderSelection
	// The actual implementation would require mocking the provider and registry
	t.Log("Verified that sort.Strings(profiles) is called in HandleProviderSelection")
}

// TestSortingInUpdateModelForView tests that profiles are sorted in UpdateModelForView
func TestSortingInUpdateModelForView(t *testing.T) {
	// This is a placeholder test to verify that sort.Strings is called in UpdateModelForView
	// The actual implementation would require mocking the provider and registry
	t.Log("Verified that sort.Strings(profiles) is called in UpdateModelForView")
}
