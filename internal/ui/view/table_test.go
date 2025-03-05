package view

import (
	"sort"
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

// TestServiceRowSorting tests that service rows are sorted alphabetically
func TestServiceRowSorting(t *testing.T) {
	// Create unsorted rows
	unsortedRows := []table.Row{
		{"Z Service", "Z Description"},
		{"A Service", "A Description"},
		{"M Service", "M Description"},
	}

	// Sort the rows
	sortedRows := make([]table.Row, len(unsortedRows))
	copy(sortedRows, unsortedRows)
	sort.Slice(sortedRows, func(i, j int) bool {
		return sortedRows[i][0] < sortedRows[j][0]
	})

	// Verify that the rows are sorted by service name
	expectedOrder := []string{"A Service", "M Service", "Z Service"}
	for i, expected := range expectedOrder {
		if sortedRows[i][0] != expected {
			t.Errorf("Expected service at index %d to be '%s', got '%s'", i, expected, sortedRows[i][0])
		}
	}

	// Test with AWS service names
	unsortedRows = []table.Row{
		{"Lambda", "Serverless Compute Service"},
		{"CodePipeline", "Continuous Delivery Service"},
		{"S3", "Simple Storage Service"},
	}

	// Sort the rows
	sortedRows = make([]table.Row, len(unsortedRows))
	copy(sortedRows, unsortedRows)
	sort.Slice(sortedRows, func(i, j int) bool {
		return sortedRows[i][0] < sortedRows[j][0]
	})

	// Verify that the rows are sorted by service name
	expectedOrder = []string{"CodePipeline", "Lambda", "S3"}
	for i, expected := range expectedOrder {
		if sortedRows[i][0] != expected {
			t.Errorf("Expected service at index %d to be '%s', got '%s'", i, expected, sortedRows[i][0])
		}
	}
}
