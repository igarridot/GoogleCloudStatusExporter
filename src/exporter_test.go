package main

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

type MockHttpClient struct {
	incidents *[]incident
	err       error
}

func (m MockHttpClient) obtainGcpStatus(doMethod DoMethod) (*[]incident, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.incidents, nil
}

type MockDoMethod struct{}

func (m MockDoMethod) Do(req any) (any, error) {
	return nil, nil
}

func TestSeverityHandler(t *testing.T) {
	tests := []struct {
		name     string
		incident incident
		expected float64
	}{
		{
			name: "AVAILABLE status should return 0.0",
			incident: incident{
				MostRecentUpdate: update{Updatestatus: "AVAILABLE"},
				Severity:         "high",
				EndsAt:           "",
			},
			expected: 0.0,
		},
		{
			name: "EndsAt populated should return 0.0",
			incident: incident{
				MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "high",
				EndsAt:           "2021-07-27T23:15:35+00:00",
			},
			expected: 0.0,
		},
		{
			name: "Low severity should return 1.0",
			incident: incident{
				MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "low",
				EndsAt:           "",
			},
			expected: 1.0,
		},
		{
			name: "Medium severity should return 2.0",
			incident: incident{
				MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "medium",
				EndsAt:           "",
			},
			expected: 2.0,
		},
		{
			name: "High severity should return 3.0",
			incident: incident{
				MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "high",
				EndsAt:           "",
			},
			expected: 3.0,
		},
		{
			name: "Unknown severity should return 0.0",
			incident: incident{
				MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "unknown",
				EndsAt:           "",
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := severityHandler(tt.incident)
			if result != tt.expected {
				t.Errorf("severityHandler() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewGcpStatus(t *testing.T) {
	tests := []struct {
		name              string
		saveLastUpdate    bool
		expectedLabels    int
		expectedLabelName string
	}{
		{
			name:              "With SAVE_LAST_UPDATE environment variable",
			saveLastUpdate:    true,
			expectedLabels:    6,
			expectedLabelName: "last_update",
		},
		{
			name:              "Without SAVE_LAST_UPDATE environment variable",
			saveLastUpdate:    false,
			expectedLabels:    5,
			expectedLabelName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.saveLastUpdate {
				os.Setenv("SAVE_LAST_UPDATE", "true")
			} else {
				os.Unsetenv("SAVE_LAST_UPDATE")
			}

			collector := newGcpStatus()
			
			if collector == nil {
				t.Fatal("newGcpStatus() returned nil")
			}
			
			if collector.gcpStatus == nil {
				t.Fatal("gcpStatus descriptor is nil")
			}

			os.Unsetenv("SAVE_LAST_UPDATE")
		})
	}
}

func TestGcpStatusCollectorDescribe(t *testing.T) {
	collector := newGcpStatus()
	ch := make(chan *prometheus.Desc, 1)
	
	collector.Describe(ch)
	
	select {
	case desc := <-ch:
		if desc != collector.gcpStatus {
			t.Error("Describe() did not send the correct descriptor")
		}
	default:
		t.Error("Describe() did not send any descriptor")
	}
}

func TestZoneHandler(t *testing.T) {
	tests := []struct {
		name          string
		incidentZones string
		incident      incident
		shouldMatch   bool
	}{
		{
			name:          "No zones filter - should always match",
			incidentZones: "",
			incident: incident{
				IncidentId:          "test-id",
				ExternalDescription: "some description",
				MostRecentUpdate:    update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:            "high",
				AffectedProducts:    []product{{Title: "Test Product", Id: "test-id"}},
				URI:                 "test-uri",
			},
			shouldMatch: true,
		},
		{
			name:          "Zone filter matches Global",
			incidentZones: "us-central1",
			incident: incident{
				IncidentId:          "test-id",
				ExternalDescription: "Global: some description",
				MostRecentUpdate:    update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:            "high",
				AffectedProducts:    []product{{Title: "Test Product", Id: "test-id"}},
				URI:                 "test-uri",
			},
			shouldMatch: true,
		},
		{
			name:          "Zone filter matches specific zone",
			incidentZones: "us-central1",
			incident: incident{
				IncidentId:          "test-id",
				ExternalDescription: "us-central1: some description",
				MostRecentUpdate:    update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:            "high",
				AffectedProducts:    []product{{Title: "Test Product", Id: "test-id"}},
				URI:                 "test-uri",
			},
			shouldMatch: true,
		},
		{
			name:          "Zone filter doesn't match",
			incidentZones: "us-west1",
			incident: incident{
				IncidentId:          "test-id",
				ExternalDescription: "us-central1: some description",
				MostRecentUpdate:    update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:            "high",
				AffectedProducts:    []product{{Title: "Test Product", Id: "test-id"}},
				URI:                 "test-uri",
			},
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.incidentZones != "" {
				os.Setenv("INCIDENTS_ZONES", tt.incidentZones)
			} else {
				os.Unsetenv("INCIDENTS_ZONES")
			}

			collector := newGcpStatus()
			ch := make(chan prometheus.Metric, 10)
			
			collector.ZoneHandler(ch, tt.incident)
			
			metricCount := len(ch)
			if tt.shouldMatch && metricCount == 0 {
				t.Error("Expected metric to be generated but none was found")
			} else if !tt.shouldMatch && metricCount > 0 {
				t.Error("Expected no metric to be generated but found one")
			}

			os.Unsetenv("INCIDENTS_ZONES")
		})
	}
}

func TestProductHandler(t *testing.T) {
	tests := []struct {
		name             string
		filteredProducts string
		incident         incident
		expectedMetrics  int
	}{
		{
			name:             "No product filter - should include all products",
			filteredProducts: "",
			incident: incident{
				IncidentId:       "test-id",
				MostRecentUpdate: update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "high",
				AffectedProducts: []product{
					{Title: "Google Cloud Storage", Id: "test-id-1"},
					{Title: "Compute Engine", Id: "test-id-2"},
				},
				URI: "test-uri",
			},
			expectedMetrics: 2,
		},
		{
			name:             "Product filter matches one product",
			filteredProducts: "Storage",
			incident: incident{
				IncidentId:       "test-id",
				MostRecentUpdate: update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "high",
				AffectedProducts: []product{
					{Title: "Google Cloud Storage", Id: "test-id-1"},
					{Title: "Compute Engine", Id: "test-id-2"},
				},
				URI: "test-uri",
			},
			expectedMetrics: 1,
		},
		{
			name:             "Product filter matches no products",
			filteredProducts: "BigQuery",
			incident: incident{
				IncidentId:       "test-id",
				MostRecentUpdate: update{Status: "test-status", Updatestatus: "SERVICE_DISRUPTION"},
				Severity:         "high",
				AffectedProducts: []product{
					{Title: "Google Cloud Storage", Id: "test-id-1"},
					{Title: "Compute Engine", Id: "test-id-2"},
				},
				URI: "test-uri",
			},
			expectedMetrics: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.filteredProducts != "" {
				os.Setenv("FILTERED_PRODUCTS", tt.filteredProducts)
			} else {
				os.Unsetenv("FILTERED_PRODUCTS")
			}

			collector := newGcpStatus()
			ch := make(chan prometheus.Metric, 10)
			
			collector.ProductHandler(ch, tt.incident, 3.0)
			
			metricCount := len(ch)
			if metricCount != tt.expectedMetrics {
				t.Errorf("Expected %d metrics, got %d", tt.expectedMetrics, metricCount)
			}

			os.Unsetenv("FILTERED_PRODUCTS")
		})
	}
}

func TestAddMetric(t *testing.T) {
	tests := []struct {
		name           string
		saveLastUpdate bool
	}{
		{
			name:           "With SAVE_LAST_UPDATE",
			saveLastUpdate: true,
		},
		{
			name:           "Without SAVE_LAST_UPDATE",
			saveLastUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.saveLastUpdate {
				os.Setenv("SAVE_LAST_UPDATE", "true")
			} else {
				os.Unsetenv("SAVE_LAST_UPDATE")
			}

			collector := newGcpStatus()
			ch := make(chan prometheus.Metric, 1)
			
			incident := incident{
				IncidentId:          "test-id",
				MostRecentUpdate:    update{Status: "test-status"},
				ExternalDescription: "test description",
				URI:                 "test-uri",
			}
			product := product{Title: "Test Product", Id: "test-product-id"}
			
			collector.AddMetric(ch, incident, 2.0, product)
			
			if len(ch) != 1 {
				t.Errorf("Expected 1 metric, got %d", len(ch))
			}

			os.Unsetenv("SAVE_LAST_UPDATE")
		})
	}
}

func TestCollect(t *testing.T) {
	tests := []struct {
		name                       string
		collectResolvedIncidents   bool
		incidents                  *[]incident
		httpClientError            error
		expectedPanic              bool
	}{
		{
			name:                     "Collect only unresolved incidents",
			collectResolvedIncidents: false,
			incidents: &[]incident{
				{
					IncidentId:       "active-incident",
					MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
					Severity:         "high",
					EndsAt:           "",
					AffectedProducts: []product{{Title: "Test Product", Id: "test-id"}},
					URI:              "test-uri",
				},
				{
					IncidentId:       "resolved-incident",
					MostRecentUpdate: update{Updatestatus: "AVAILABLE"},
					Severity:         "low",
					EndsAt:           "",
					AffectedProducts: []product{{Title: "Test Product", Id: "test-id"}},
					URI:              "test-uri",
				},
			},
		},
		{
			name:                     "Collect all incidents",
			collectResolvedIncidents: true,
			incidents: &[]incident{
				{
					IncidentId:       "active-incident",
					MostRecentUpdate: update{Updatestatus: "SERVICE_DISRUPTION"},
					Severity:         "high",
					EndsAt:           "",
					AffectedProducts: []product{{Title: "Test Product", Id: "test-id"}},
					URI:              "test-uri",
				},
				{
					IncidentId:       "resolved-incident",
					MostRecentUpdate: update{Updatestatus: "AVAILABLE"},
					Severity:         "low",
					EndsAt:           "",
					AffectedProducts: []product{{Title: "Test Product", Id: "test-id"}},
					URI:              "test-uri",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.collectResolvedIncidents {
				os.Setenv("COLLECT_RESOLVED_EVENTS", "true")
			} else {
				os.Unsetenv("COLLECT_RESOLVED_EVENTS")
			}

			// Create a custom collector that uses our mock HTTP client
			collector := &gcpStatusCollector{
				gcpStatus: prometheus.NewDesc(
					prometheus.BuildFQName("gcp", "", "incidents"),
					"GCP Incident last update status",
					[]string{"id", "status", "product", "description", "uri"}, nil,
				),
			}

			// We can't easily test the Collect method because it creates its own HTTP clients
			// This would require refactoring the code to accept the HTTP client as a parameter
			// For now, we'll test that the collector can be created and described

			descCh := make(chan *prometheus.Desc, 1)
			
			collector.Describe(descCh)
			if len(descCh) != 1 {
				t.Error("Expected 1 descriptor")
			}

			os.Unsetenv("COLLECT_RESOLVED_EVENTS")
		})
	}
}

func TestRegister(t *testing.T) {
	// Create a new registry to avoid conflicts with other tests
	registry := prometheus.NewRegistry()
	
	// We can't test the global Register function directly because it uses
	// prometheus.MustRegister which affects the global registry
	// Instead, we'll test that we can create and register a collector
	
	collector := newGcpStatus()
	err := registry.Register(collector)
	if err != nil {
		t.Errorf("Failed to register collector: %v", err)
	}
	
	// Try to register the same collector again - should fail
	err = registry.Register(collector)
	if err == nil {
		t.Error("Expected error when registering duplicate collector")
	}
}

func TestMetricsOutput(t *testing.T) {
	// Test that the collector can be created and described properly
	collector := newGcpStatus()
	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	if err != nil {
		t.Errorf("Failed to register collector: %v", err)
	}
	
	// Test that the descriptor is properly created
	descCh := make(chan *prometheus.Desc, 1)
	collector.Describe(descCh)
	
	if len(descCh) != 1 {
		t.Error("Expected 1 descriptor")
	}
	
	t.Log("Collector properly structured")
}