package valueobjects

import (
	"testing"
)

func TestNewFilterCriteria(t *testing.T) {
	tests := []struct {
		name             string
		zones            string
		products         string
		expectedZones    []string
		expectedProducts []string
	}{
		{
			name:             "Empty filters",
			zones:            "",
			products:         "",
			expectedZones:    []string{},
			expectedProducts: []string{},
		},
		{
			name:             "Single zone",
			zones:            "us-east1",
			products:         "",
			expectedZones:    []string{"us-east1"},
			expectedProducts: []string{},
		},
		{
			name:             "Multiple zones",
			zones:            "us-east1,europe-west1",
			products:         "",
			expectedZones:    []string{"us-east1", "europe-west1"},
			expectedProducts: []string{},
		},
		{
			name:             "Single product",
			zones:            "",
			products:         "Apigee",
			expectedZones:    []string{},
			expectedProducts: []string{"Apigee"},
		},
		{
			name:             "Multiple products",
			zones:            "",
			products:         "Apigee,Compute Engine",
			expectedZones:    []string{},
			expectedProducts: []string{"Apigee", "Compute Engine"},
		},
		{
			name:             "Both zones and products",
			zones:            "us-east1,europe-west1",
			products:         "Apigee,BigQuery",
			expectedZones:    []string{"us-east1", "europe-west1"},
			expectedProducts: []string{"Apigee", "BigQuery"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria(tt.zones, tt.products)

			if len(criteria.Zones) != len(tt.expectedZones) {
				t.Errorf("Expected %d zones, got %d", len(tt.expectedZones), len(criteria.Zones))
			}
			for i, expectedZone := range tt.expectedZones {
				if i >= len(criteria.Zones) || criteria.Zones[i] != expectedZone {
					t.Errorf("Expected zone[%d] '%s', got '%s'", i, expectedZone, criteria.Zones[i])
				}
			}

			if len(criteria.Products) != len(tt.expectedProducts) {
				t.Errorf("Expected %d products, got %d", len(tt.expectedProducts), len(criteria.Products))
			}
			for i, expectedProduct := range tt.expectedProducts {
				if i >= len(criteria.Products) || criteria.Products[i] != expectedProduct {
					t.Errorf("Expected product[%d] '%s', got '%s'", i, expectedProduct, criteria.Products[i])
				}
			}
		})
	}
}

func TestHasZoneFilter(t *testing.T) {
	tests := []struct {
		name     string
		zones    string
		expected bool
	}{
		{
			name:     "No zones",
			zones:    "",
			expected: false,
		},
		{
			name:     "Single zone",
			zones:    "us-east1",
			expected: true,
		},
		{
			name:     "Multiple zones",
			zones:    "us-east1,europe-west1",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria(tt.zones, "")
			result := criteria.HasZoneFilter()
			if result != tt.expected {
				t.Errorf("Expected HasZoneFilter %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHasProductFilter(t *testing.T) {
	tests := []struct {
		name     string
		products string
		expected bool
	}{
		{
			name:     "No products",
			products: "",
			expected: false,
		},
		{
			name:     "Single product",
			products: "Apigee",
			expected: true,
		},
		{
			name:     "Multiple products",
			products: "Apigee,Compute Engine",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria("", tt.products)
			result := criteria.HasProductFilter()
			if result != tt.expected {
				t.Errorf("Expected HasProductFilter %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNewIncidentID_EmptyString(t *testing.T) {
	_, err := NewIncidentID("")

	if err == nil {
		t.Error("Expected error when creating IncidentID with empty string")
	}

	expectedError := "incident ID cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewUpdateStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected UpdateStatus
	}{
		{"AVAILABLE status", "AVAILABLE", UpdateStatusAvailable},
		{"INVESTIGATING status", "INVESTIGATING", UpdateStatus("INVESTIGATING")},
		{"Empty string", "", UpdateStatus("")},
		{"Custom status", "CUSTOM_STATUS", UpdateStatus("CUSTOM_STATUS")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewUpdateStatus(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
