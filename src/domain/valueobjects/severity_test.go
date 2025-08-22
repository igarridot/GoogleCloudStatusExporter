package valueobjects

import "testing"

func TestSeverity_ToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		severity Severity
		expected float64
	}{
		{
			name:     "low severity returns 1.0",
			severity: SeverityLow,
			expected: 1.0,
		},
		{
			name:     "medium severity returns 2.0",
			severity: SeverityMedium,
			expected: 2.0,
		},
		{
			name:     "high severity returns 3.0",
			severity: SeverityHigh,
			expected: 3.0,
		},
		{
			name:     "unknown severity returns 0.0",
			severity: Severity("unknown"),
			expected: 0.0,
		},
		{
			name:     "empty severity returns 0.0",
			severity: Severity(""),
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.severity.ToFloat64()
			if result != tt.expected {
				t.Errorf("ToFloat64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewSeverity(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected Severity
	}{
		{
			name:     "low string returns SeverityLow",
			value:    "low",
			expected: SeverityLow,
		},
		{
			name:     "medium string returns SeverityMedium",
			value:    "medium",
			expected: SeverityMedium,
		},
		{
			name:     "high string returns SeverityHigh",
			value:    "high",
			expected: SeverityHigh,
		},
		{
			name:     "unknown string returns empty severity",
			value:    "unknown",
			expected: Severity(""),
		},
		{
			name:     "empty string returns empty severity",
			value:    "",
			expected: Severity(""),
		},
		{
			name:     "case sensitive - Low returns empty severity",
			value:    "Low",
			expected: Severity(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSeverity(tt.value)
			if result != tt.expected {
				t.Errorf("NewSeverity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSeverityConstants(t *testing.T) {
	if SeverityLow != "low" {
		t.Errorf("SeverityLow = %v, want %v", SeverityLow, "low")
	}

	if SeverityMedium != "medium" {
		t.Errorf("SeverityMedium = %v, want %v", SeverityMedium, "medium")
	}

	if SeverityHigh != "high" {
		t.Errorf("SeverityHigh = %v, want %v", SeverityHigh, "high")
	}
}