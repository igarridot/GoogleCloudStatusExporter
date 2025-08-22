package services

import (
	"testing"
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

func TestIncidentFilterService_FilterIncidents(t *testing.T) {
	service := NewIncidentFilterService()

	incidents := createTestIncidents()

	t.Run("Filter by zone", func(t *testing.T) {
		criteria := valueobjects.NewFilterCriteria("us-east1", "")
		filtered := service.FilterIncidents(incidents, criteria, true)

		if len(filtered) != 1 {
			t.Errorf("Expected 1 incident, got %d", len(filtered))
		}
	})

	t.Run("Filter by product", func(t *testing.T) {
		criteria := valueobjects.NewFilterCriteria("", "Compute")
		filtered := service.FilterIncidents(incidents, criteria, true)

		if len(filtered) != 1 {
			t.Errorf("Expected 1 incident, got %d", len(filtered))
		}
	})

	t.Run("Filter resolved incidents", func(t *testing.T) {
		criteria := valueobjects.NewFilterCriteria("", "")
		filtered := service.FilterIncidents(incidents, criteria, false)

		if len(filtered) != 1 {
			t.Errorf("Expected 1 incident (non-resolved), got %d", len(filtered))
		}
	})

	t.Run("No filter", func(t *testing.T) {
		criteria := valueobjects.NewFilterCriteria("", "")
		filtered := service.FilterIncidents(incidents, criteria, true)

		if len(filtered) != 2 {
			t.Errorf("Expected 2 incidents, got %d", len(filtered))
		}
	})
}

func createTestIncidents() []entities.Incident {
	id1, _ := valueobjects.NewIncidentID("incident-1")
	id2, _ := valueobjects.NewIncidentID("incident-2")

	return []entities.Incident{
		{
			ID:                  id1,
			ExternalDescription: "Issue in us-east1 zone",
			Severity:            valueobjects.SeverityHigh,
			AffectedProducts:    []entities.Product{{Title: "Compute Engine", ID: "compute"}},
			MostRecentUpdate:    entities.Update{UpdateStatus: "INVESTIGATING"},
		},
		{
			ID:                  id2,
			ExternalDescription: "Resolved issue",
			Severity:            valueobjects.SeverityLow,
			EndTime:             &time.Time{},
			AffectedProducts:    []entities.Product{{Title: "Cloud Storage", ID: "storage"}},
			MostRecentUpdate:    entities.Update{UpdateStatus: valueobjects.UpdateStatusAvailable},
		},
	}
}
