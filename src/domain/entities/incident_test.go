package entities

import (
	"testing"
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

func TestIncident_IsResolved(t *testing.T) {
	tests := []struct {
		name     string
		incident Incident
		expected bool
	}{
		{
			name: "incident is resolved when update status is AVAILABLE",
			incident: Incident{
				MostRecentUpdate: Update{
					UpdateStatus: valueobjects.UpdateStatusAvailable,
				},
			},
			expected: true,
		},
		{
			name: "incident is resolved when end time is set",
			incident: Incident{
				EndTime: &time.Time{},
				MostRecentUpdate: Update{
					UpdateStatus: "OTHER_STATUS",
				},
			},
			expected: true,
		},
		{
			name: "incident is not resolved when update status is not AVAILABLE and no end time",
			incident: Incident{
				EndTime: nil,
				MostRecentUpdate: Update{
					UpdateStatus: "OTHER_STATUS",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.incident.IsResolved()
			if result != tt.expected {
				t.Errorf("IsResolved() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIncident_GetSeverityValue(t *testing.T) {
	tests := []struct {
		name     string
		incident Incident
		expected float64
	}{
		{
			name: "resolved incident returns 0.0",
			incident: Incident{
				MostRecentUpdate: Update{
					UpdateStatus: valueobjects.UpdateStatusAvailable,
				},
				Severity: valueobjects.SeverityHigh,
			},
			expected: 0.0,
		},
		{
			name: "unresolved incident with high severity returns 3.0",
			incident: Incident{
				MostRecentUpdate: Update{
					UpdateStatus: "OTHER_STATUS",
				},
				Severity: valueobjects.SeverityHigh,
			},
			expected: 3.0,
		},
		{
			name: "unresolved incident with medium severity returns 2.0",
			incident: Incident{
				MostRecentUpdate: Update{
					UpdateStatus: "OTHER_STATUS",
				},
				Severity: valueobjects.SeverityMedium,
			},
			expected: 2.0,
		},
		{
			name: "unresolved incident with low severity returns 1.0",
			incident: Incident{
				MostRecentUpdate: Update{
					UpdateStatus: "OTHER_STATUS",
				},
				Severity: valueobjects.SeverityLow,
			},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.incident.GetSeverityValue()
			if result != tt.expected {
				t.Errorf("GetSeverityValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProduct(t *testing.T) {
	product := Product{
		Title: "Test Product",
		ID:    "test-id",
	}

	if product.Title != "Test Product" {
		t.Errorf("Product.Title = %v, want %v", product.Title, "Test Product")
	}

	if product.ID != "test-id" {
		t.Errorf("Product.ID = %v, want %v", product.ID, "test-id")
	}
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	update := Update{
		CreatedAt:    now,
		Status:       "INVESTIGATING",
		UpdateStatus: valueobjects.UpdateStatusAvailable,
	}

	if update.CreatedAt != now {
		t.Errorf("Update.CreatedAt = %v, want %v", update.CreatedAt, now)
	}

	if update.Status != "INVESTIGATING" {
		t.Errorf("Update.Status = %v, want %v", update.Status, "INVESTIGATING")
	}

	if update.UpdateStatus != valueobjects.UpdateStatusAvailable {
		t.Errorf("Update.UpdateStatus = %v, want %v", update.UpdateStatus, valueobjects.UpdateStatusAvailable)
	}
}