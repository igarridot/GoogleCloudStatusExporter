package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/services"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
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

func (m *mockMetricsPort) CollectMetricsWithChannel(incidents []entities.Incident, ch chan<- interface{}) {
	// Mock implementation
}

func createTestIncidents() []entities.Incident {
	id1, _ := valueobjects.NewIncidentID("incident-1")
	id2, _ := valueobjects.NewIncidentID("incident-2")
	id3, _ := valueobjects.NewIncidentID("incident-3")
	id4, _ := valueobjects.NewIncidentID("incident-4")
	id5, _ := valueobjects.NewIncidentID("incident-5")

	endTime := time.Now()

	return []entities.Incident{
		{
			ID:                  id1,
			ExternalDescription: "Issue in us-east1 zone affecting Apigee",
			Severity:            valueobjects.SeverityHigh,
			AffectedProducts: []entities.Product{
				{Title: "Apigee", ID: "apigee"},
				{Title: "Compute Engine", ID: "compute"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: "INVESTIGATING",
			},
		},
		{
			ID:                  id2,
			ExternalDescription: "Issue in europe-west1 zone",
			Severity:            valueobjects.SeverityLow,
			EndTime:             &endTime,
			AffectedProducts: []entities.Product{
				{Title: "Cloud Storage", ID: "storage"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: valueobjects.UpdateStatusAvailable,
			},
		},
		{
			ID:                  id3,
			ExternalDescription: "Worldwide issue affecting multiple products",
			Severity:            valueobjects.SeverityMedium,
			AffectedProducts: []entities.Product{
				{Title: "Apigee", ID: "apigee"},
				{Title: "BigQuery", ID: "bigquery"},
				{Title: "Cloud SQL", ID: "cloudsql"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: "SERVICE_DISRUPTION",
			},
		},
		{
			ID:                  id4,
			ExternalDescription: "us-central1 Compute Engine issue",
			Severity:            valueobjects.SeverityHigh,
			AffectedProducts: []entities.Product{
				{Title: "Compute Engine", ID: "compute"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: "INVESTIGATING",
			},
		},
		{
			ID:                  id5,
			ExternalDescription: "Global incident affecting various services",
			Severity:            valueobjects.SeverityHigh,
			AffectedProducts: []entities.Product{
				{Title: "IAM", ID: "iam"},
				{Title: "Cloud Logging", ID: "logging"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: "INVESTIGATING",
			},
		},
	}
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
			name:                  "No filters, exclude resolved",
			incidents:             createTestIncidents(),
			filterCriteria:        valueobjects.NewFilterCriteria("", ""),
			collectResolved:       false,
			expectedIncidentCount: 4,
			expectedError:         false,
		},
		{
			name:                  "Filter by product",
			incidents:             createTestIncidents(),
			filterCriteria:        valueobjects.NewFilterCriteria("", "Apigee"),
			collectResolved:       false,
			expectedIncidentCount: 2,
			expectedError:         false,
		},
		{
			name:                  "Filter by zone",
			incidents:             createTestIncidents(),
			filterCriteria:        valueobjects.NewFilterCriteria("us-east1", ""),
			collectResolved:       false,
			expectedIncidentCount: 2,
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

func TestCollectMetricsUseCase_FilterCombinations(t *testing.T) {
	// Test all possible combinations of filter criteria and collectResolved flag
	filterCombinations := []struct {
		name            string
		zones           string
		products        string
		collectResolved bool
	}{
		{"No filters, no resolved", "", "", false},
		{"No filters, with resolved", "", "", true},
		{"Zone only, no resolved", "us-east1", "", false},
		{"Zone only, with resolved", "us-east1", "", true},
		{"Product only, no resolved", "", "Apigee", false},
		{"Product only, with resolved", "", "Apigee", true},
		{"Both filters, no resolved", "us-east1", "Apigee", false},
		{"Both filters, with resolved", "us-east1", "Apigee", true},
		{"Multiple zones, no resolved", "us-east1,europe-west1", "", false},
		{"Multiple zones, with resolved", "us-east1,europe-west1", "", true},
		{"Multiple products, no resolved", "", "Apigee,Compute", false},
		{"Multiple products, with resolved", "", "Apigee,Compute", true},
		{"Multiple both, no resolved", "us-east1,Global", "Apigee,BigQuery", false},
		{"Multiple both, with resolved", "us-east1,Global", "Apigee,BigQuery", true},
	}

	incidents := createTestIncidents()

	for _, combo := range filterCombinations {
		t.Run(combo.name, func(t *testing.T) {
			// Create mocks
			gcpStatusPort := &mockGCPStatusPort{incidents: incidents}
			metricsPort := &mockMetricsPort{}
			filterService := services.NewIncidentFilterService()

			// Create use case
			useCase := NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

			// Execute
			filterCriteria := valueobjects.NewFilterCriteria(combo.zones, combo.products)
			metricsConfig := ports.MetricsConfig{SaveLastUpdate: false}
			result, err := useCase.Execute(filterCriteria, combo.collectResolved, metricsConfig)

			// Should not error
			if err != nil {
				t.Errorf("Unexpected error for combination %s: %v", combo.name, err)
			}

			// Should return some result (exact count depends on the combination)
			if result == nil {
				t.Errorf("Expected non-nil result for combination %s", combo.name)
			}

			// Log for debugging
			t.Logf("Combination %s returned %d incidents", combo.name, len(result))
		})
	}
}

func TestCollectMetricsUseCase_MetricsConfigVariations(t *testing.T) {
	metricsConfigs := []struct {
		name           string
		saveLastUpdate bool
	}{
		{"Save last update enabled", true},
		{"Save last update disabled", false},
	}

	incidents := createTestIncidents()

	for _, config := range metricsConfigs {
		t.Run(config.name, func(t *testing.T) {
			// Create mocks
			gcpStatusPort := &mockGCPStatusPort{incidents: incidents}
			metricsPort := &mockMetricsPort{}
			filterService := services.NewIncidentFilterService()

			// Create use case
			useCase := NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

			// Execute
			filterCriteria := valueobjects.NewFilterCriteria("", "")
			metricsConfig := ports.MetricsConfig{SaveLastUpdate: config.saveLastUpdate}
			result, err := useCase.Execute(filterCriteria, false, metricsConfig)

			// Should not error
			if err != nil {
				t.Errorf("Unexpected error for metrics config %s: %v", config.name, err)
			}

			// Should return expected number of incidents (4, excluding resolved)
			if len(result) != 4 {
				t.Errorf("Expected 4 incidents for config %s, got %d", config.name, len(result))
			}
		})
	}
}
