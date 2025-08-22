package valueobjects

import "testing"

func TestNewUpdateStatus(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected UpdateStatus
	}{
		{
			name:     "AVAILABLE string returns UpdateStatusAvailable",
			value:    "AVAILABLE",
			expected: UpdateStatusAvailable,
		},
		{
			name:     "INVESTIGATING string returns UpdateStatus with value",
			value:    "INVESTIGATING",
			expected: UpdateStatus("INVESTIGATING"),
		},
		{
			name:     "MONITORING string returns UpdateStatus with value",
			value:    "MONITORING",
			expected: UpdateStatus("MONITORING"),
		},
		{
			name:     "empty string returns empty UpdateStatus",
			value:    "",
			expected: UpdateStatus(""),
		},
		{
			name:     "arbitrary string returns UpdateStatus with value",
			value:    "CUSTOM_STATUS",
			expected: UpdateStatus("CUSTOM_STATUS"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewUpdateStatus(tt.value)
			if result != tt.expected {
				t.Errorf("NewUpdateStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUpdateStatusConstants(t *testing.T) {
	if UpdateStatusAvailable != "AVAILABLE" {
		t.Errorf("UpdateStatusAvailable = %v, want %v", UpdateStatusAvailable, "AVAILABLE")
	}
}

func TestUpdateStatusEquality(t *testing.T) {
	status1 := NewUpdateStatus("AVAILABLE")
	status2 := UpdateStatusAvailable

	if status1 != status2 {
		t.Errorf("NewUpdateStatus(\"AVAILABLE\") != UpdateStatusAvailable")
	}
}