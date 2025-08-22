package usecases

import (
	"errors"
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/services"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/testutils"
)

// Mock implementations
type mockGCPStatusPort struct {
	incidents []entities.Incident
	err       error
}

func (m *mockGCPStatusPort) GetIncidents() ([]entities.Incident, error) {
	return m.incidents, m.err
}

type mockMetricsPort struct{}

func (m *mockMetricsPort) CollectMetrics(incidents []entities.Incident, config ports.MetricsConfig) error {
	return nil
}

func TestCollectMetricsUseCase_Execute(t *testing.T) {
	tests := []struct {
		name                  string
		incidents             []entities.Incident
		filterCriteria        valueobjects.FilterCriteria
		collectResolved       bool
		expectedIncidentCount int
		expectedError         bool
		gcpStatusPortError    error
	}{
		{
			name:                  "Success case - exclude resolved",
			incidents:             testutils.CreateTestIncidents(),
			filterCriteria:        valueobjects.NewFilterCriteria("", ""),
			collectResolved:       false,
			expectedIncidentCount: 1,
			expectedError:         false,
		},
		{
			name:                  "GCP Status Port error",
			incidents:             nil,
			filterCriteria:        valueobjects.NewFilterCriteria("", ""),
			collectResolved:       false,
			expectedIncidentCount: 0,
			expectedError:         true,
			gcpStatusPortError:    errors.New("failed to fetch incidents"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			gcpStatusPort := &mockGCPStatusPort{
				incidents: tt.incidents,
				err:       tt.gcpStatusPortError,
			}
			metricsPort := &mockMetricsPort{}
			filterService := services.NewIncidentFilterService()

			// Create use case
			useCase := NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

			// Execute
			metricsConfig := ports.MetricsConfig{SaveLastUpdate: false}
			result, err := useCase.Execute(tt.filterCriteria, tt.collectResolved, metricsConfig)

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Validate result count
			if len(result) != tt.expectedIncidentCount {
				t.Errorf("Expected %d incidents, got %d.",
					tt.expectedIncidentCount, len(result))

				// Debug information
				t.Logf("Returned incidents:")
				for i, incident := range result {
					t.Logf("  [%d] ID: %s, Description: %s", i, incident.ID, incident.ExternalDescription)
					for j, product := range incident.AffectedProducts {
						t.Logf("    Product[%d]: %s", j, product.Title)
					}
				}
			}
		})
	}
}
