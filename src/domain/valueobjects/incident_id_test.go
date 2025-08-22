package valueobjects

import "testing"

func TestNewIncidentID(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    IncidentID
		expectError bool
	}{
		{
			name:        "valid id returns IncidentID",
			value:       "incident-123",
			expected:    IncidentID("incident-123"),
			expectError: false,
		},
		{
			name:        "empty string returns error",
			value:       "",
			expected:    IncidentID(""),
			expectError: true,
		},
		{
			name:        "uuid format returns IncidentID",
			value:       "550e8400-e29b-41d4-a716-446655440000",
			expected:    IncidentID("550e8400-e29b-41d4-a716-446655440000"),
			expectError: false,
		},
		{
			name:        "numeric string returns IncidentID",
			value:       "12345",
			expected:    IncidentID("12345"),
			expectError: false,
		},
		{
			name:        "alphanumeric with special chars returns IncidentID",
			value:       "incident_123-test",
			expected:    IncidentID("incident_123-test"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewIncidentID(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("NewIncidentID() expected error but got none")
				}
				if err.Error() != "incident ID cannot be empty" {
					t.Errorf("NewIncidentID() error = %v, want %v", err.Error(), "incident ID cannot be empty")
				}
			} else {
				if err != nil {
					t.Errorf("NewIncidentID() unexpected error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("NewIncidentID() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestIncidentID_String(t *testing.T) {
	tests := []struct {
		name       string
		incidentID IncidentID
		expected   string
	}{
		{
			name:       "incident id returns string representation",
			incidentID: IncidentID("incident-123"),
			expected:   "incident-123",
		},
		{
			name:       "empty incident id returns empty string",
			incidentID: IncidentID(""),
			expected:   "",
		},
		{
			name:       "uuid incident id returns string representation",
			incidentID: IncidentID("550e8400-e29b-41d4-a716-446655440000"),
			expected:   "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.incidentID.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIncidentIDIntegration(t *testing.T) {
	// Test the full workflow of creating and using an IncidentID
	originalValue := "test-incident-456"
	
	incidentID, err := NewIncidentID(originalValue)
	if err != nil {
		t.Fatalf("NewIncidentID() unexpected error = %v", err)
	}

	stringValue := incidentID.String()
	if stringValue != originalValue {
		t.Errorf("String() = %v, want %v", stringValue, originalValue)
	}

	// Test that the IncidentID can be compared
	anotherID, _ := NewIncidentID(originalValue)
	if incidentID != anotherID {
		t.Errorf("IncidentID comparison failed: %v != %v", incidentID, anotherID)
	}
}