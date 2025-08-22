package services

import (
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/testutils"
)

func TestIncidentFilterService_FilterIncidents(t *testing.T) {
	service := NewIncidentFilterService()

	incidents := testutils.CreateFilterTestIncidents()

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
