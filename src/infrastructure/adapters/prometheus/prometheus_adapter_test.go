package prometheus

import (
	"testing"
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewPrometheusAdapter(t *testing.T) {
	tests := []struct {
		name                string
		config              ports.MetricsConfig
		expectedLabelsCount int
	}{
		{
			name: "adapter with SaveLastUpdate false has 5 labels",
			config: ports.MetricsConfig{
				SaveLastUpdate: false,
			},
			expectedLabelsCount: 5,
		},
		{
			name: "adapter with SaveLastUpdate true has 6 labels",
			config: ports.MetricsConfig{
				SaveLastUpdate: true,
			},
			expectedLabelsCount: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewPrometheusAdapter(tt.config)

			if adapter == nil {
				t.Fatal("NewPrometheusAdapter() returned nil")
			}

			if adapter.config.SaveLastUpdate != tt.config.SaveLastUpdate {
				t.Errorf("config.SaveLastUpdate = %v, want %v", adapter.config.SaveLastUpdate, tt.config.SaveLastUpdate)
			}

			if adapter.gcpStatus == nil {
				t.Fatal("gcpStatus descriptor is nil")
			}
		})
	}
}

func TestPrometheusAdapter_Describe(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	ch := make(chan *prometheus.Desc, 1)
	adapter.Describe(ch)
	close(ch)

	descriptors := make([]*prometheus.Desc, 0)
	for desc := range ch {
		descriptors = append(descriptors, desc)
	}

	if len(descriptors) != 1 {
		t.Errorf("Expected 1 descriptor, got %d", len(descriptors))
	}
}

func TestPrometheusAdapter_CollectMetrics(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	incidents := createTestIncidents()
	err := adapter.CollectMetrics(incidents, config)

	if err != nil {
		t.Errorf("CollectMetrics() returned error: %v", err)
	}
}

func TestPrometheusAdapter_CollectMetricsWithChannel(t *testing.T) {
	tests := []struct {
		name              string
		config            ports.MetricsConfig
		incidents         []entities.Incident
		expectedMetrics   int
	}{
		{
			name:            "single incident with one product produces one metric",
			config:          ports.MetricsConfig{SaveLastUpdate: false},
			incidents:       createTestIncidents()[:1],
			expectedMetrics: 1,
		},
		{
			name:            "incident with multiple products produces multiple metrics",
			config:          ports.MetricsConfig{SaveLastUpdate: true},
			incidents:       createIncidentWithMultipleProducts(),
			expectedMetrics: 2,
		},
		{
			name:            "multiple incidents produce multiple metrics",
			config:          ports.MetricsConfig{SaveLastUpdate: false},
			incidents:       createTestIncidents(),
			expectedMetrics: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewPrometheusAdapter(tt.config)
			
			registry := prometheus.NewRegistry()
			collector := &testCollector{adapter: adapter, incidents: tt.incidents}
			registry.MustRegister(collector)

			metricFamilies, err := registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			if len(metricFamilies) == 0 {
				t.Fatal("No metric families found")
			}

			totalMetrics := 0
			for _, mf := range metricFamilies {
				totalMetrics += len(mf.GetMetric())
			}

			if totalMetrics != tt.expectedMetrics {
				t.Errorf("Expected %d metrics, got %d", tt.expectedMetrics, totalMetrics)
			}
		})
	}
}

func TestPrometheusAdapter_MetricValues(t *testing.T) {
	tests := []struct {
		name          string
		config        ports.MetricsConfig
		incident      entities.Incident
		expectedValue float64
	}{
		{
			name:   "resolved incident returns 0.0",
			config: ports.MetricsConfig{SaveLastUpdate: false},
			incident: entities.Incident{
				ID: mustCreateIncidentID("test-1"),
				MostRecentUpdate: entities.Update{
					Status:       "AVAILABLE",
					UpdateStatus: valueobjects.UpdateStatusAvailable,
				},
				Severity: valueobjects.SeverityHigh,
				AffectedProducts: []entities.Product{
					{Title: "Test Product", ID: "test-product"},
				},
				ExternalDescription: "Test incident",
				URI:                "incident/test-1",
			},
			expectedValue: 0.0,
		},
		{
			name:   "high severity unresolved incident returns 3.0",
			config: ports.MetricsConfig{SaveLastUpdate: false},
			incident: entities.Incident{
				ID: mustCreateIncidentID("test-2"),
				MostRecentUpdate: entities.Update{
					Status:       "INVESTIGATING",
					UpdateStatus: "INVESTIGATING",
				},
				Severity: valueobjects.SeverityHigh,
				AffectedProducts: []entities.Product{
					{Title: "Test Product", ID: "test-product"},
				},
				ExternalDescription: "Test incident",
				URI:                "incident/test-2",
			},
			expectedValue: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewPrometheusAdapter(tt.config)
			
			registry := prometheus.NewRegistry()
			collector := &testCollector{adapter: adapter, incidents: []entities.Incident{tt.incident}}
			registry.MustRegister(collector)

			metricFamilies, err := registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			if len(metricFamilies) == 0 {
				t.Fatal("No metric families found")
			}

			metric := metricFamilies[0].GetMetric()[0]
			actualValue := metric.GetGauge().GetValue()

			if actualValue != tt.expectedValue {
				t.Errorf("Expected metric value %v, got %v", tt.expectedValue, actualValue)
			}
		})
	}
}

// Helper types and functions for testing

type testCollector struct {
	adapter   *PrometheusAdapter
	incidents []entities.Incident
}

func (tc *testCollector) Describe(ch chan<- *prometheus.Desc) {
	tc.adapter.Describe(ch)
}

func (tc *testCollector) Collect(ch chan<- prometheus.Metric) {
	tc.adapter.CollectMetricsWithChannel(tc.incidents, ch)
}

func createTestIncidents() []entities.Incident {
	return []entities.Incident{
		{
			ID:     mustCreateIncidentID("incident-1"),
			Number: "1",
			BeginTime: time.Now(),
			CreatedAt: time.Now(),
			ExternalDescription: "Test incident 1",
			MostRecentUpdate: entities.Update{
				Status:       "INVESTIGATING",
				UpdateStatus: "INVESTIGATING",
			},
			StatusImpact: "SERVICE_DISRUPTION",
			Severity:     valueobjects.SeverityHigh,
			AffectedProducts: []entities.Product{
				{Title: "Compute Engine", ID: "compute"},
			},
			ServiceKey:  "compute",
			ServiceName: "Compute Engine",
			URI:         "incident/incident-1",
		},
		{
			ID:     mustCreateIncidentID("incident-2"),
			Number: "2",
			BeginTime: time.Now(),
			CreatedAt: time.Now(),
			ExternalDescription: "Test incident 2",
			MostRecentUpdate: entities.Update{
				Status:       "AVAILABLE",
				UpdateStatus: valueobjects.UpdateStatusAvailable,
			},
			StatusImpact: "SERVICE_DISRUPTION",
			Severity:     valueobjects.SeverityMedium,
			AffectedProducts: []entities.Product{
				{Title: "Cloud Storage", ID: "storage"},
			},
			ServiceKey:  "storage",
			ServiceName: "Cloud Storage",
			URI:         "incident/incident-2",
		},
	}
}

func createIncidentWithMultipleProducts() []entities.Incident {
	return []entities.Incident{
		{
			ID:     mustCreateIncidentID("incident-multi"),
			Number: "multi",
			BeginTime: time.Now(),
			CreatedAt: time.Now(),
			ExternalDescription: "Multi-product incident",
			MostRecentUpdate: entities.Update{
				Status:       "INVESTIGATING",
				UpdateStatus: "INVESTIGATING",
			},
			StatusImpact: "SERVICE_DISRUPTION",
			Severity:     valueobjects.SeverityLow,
			AffectedProducts: []entities.Product{
				{Title: "Compute Engine", ID: "compute"},
				{Title: "Cloud Storage", ID: "storage"},
			},
			ServiceKey:  "multi",
			ServiceName: "Multiple Services",
			URI:         "incident/incident-multi",
		},
	}
}

func mustCreateIncidentID(value string) valueobjects.IncidentID {
	id, err := valueobjects.NewIncidentID(value)
	if err != nil {
		panic(err)
	}
	return id
}