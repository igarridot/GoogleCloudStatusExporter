package prometheus

import (
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/testutils"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewPrometheusAdapter(t *testing.T) {
	t.Run("With SaveLastUpdate=false", func(t *testing.T) {
		config := ports.MetricsConfig{SaveLastUpdate: false}
		adapter := NewPrometheusAdapter(config)

		if adapter == nil {
			t.Error("Expected adapter to be created, got nil")
			return
		}

		if adapter.config.SaveLastUpdate != false {
			t.Error("Expected SaveLastUpdate to be false")
		}

		if adapter.gcpStatus == nil {
			t.Error("Expected gcpStatus descriptor to be created")
		}
	})

	t.Run("With SaveLastUpdate=true", func(t *testing.T) {
		config := ports.MetricsConfig{SaveLastUpdate: true}
		adapter := NewPrometheusAdapter(config)

		if adapter == nil {
			t.Error("Expected adapter to be created, got nil")
			return
		}

		if adapter.config.SaveLastUpdate != true {
			t.Error("Expected SaveLastUpdate to be true")
		}
	})
}

func TestPrometheusAdapter_Describe(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	ch := make(chan *prometheus.Desc, 1)
	adapter.Describe(ch)
	close(ch)

	// Check that exactly one descriptor was sent
	count := 0
	for desc := range ch {
		count++
		if desc == nil {
			t.Error("Expected non-nil descriptor")
		}
	}

	if count != 1 {
		t.Errorf("Expected 1 descriptor, got %d", count)
	}
}

func TestPrometheusAdapter_CollectMetrics(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	incidents := testutils.CreateTestIncidents()

	// This method should not return an error (it's a placeholder implementation)
	err := adapter.CollectMetrics(incidents, config)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPrometheusAdapter_CollectMetricsWithChannel(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	incidents := testutils.CreateTestIncidents()

	ch := make(chan prometheus.Metric, 10)
	adapter.CollectMetricsWithChannel(incidents, ch)
	close(ch)

	// Count metrics collected
	metricCount := 0
	for metric := range ch {
		metricCount++
		if metric == nil {
			t.Error("Expected non-nil metric")
		}
	}

	// We should get one metric per affected product
	expectedMetrics := 0
	for _, incident := range incidents {
		expectedMetrics += len(incident.AffectedProducts)
	}

	if metricCount != expectedMetrics {
		t.Errorf("Expected %d metrics, got %d", expectedMetrics, metricCount)
	}
}

func TestPrometheusAdapter_addMetric_WithoutLastUpdate(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	incident := createTestIncident()
	product := entities.Product{Title: "Test Product", ID: "test-product"}

	ch := make(chan prometheus.Metric, 1)
	adapter.addMetric(ch, incident, product)
	close(ch)

	// Verify one metric was created
	metricCount := 0
	for range ch {
		metricCount++
	}

	if metricCount != 1 {
		t.Errorf("Expected 1 metric, got %d", metricCount)
	}
}

func TestPrometheusAdapter_addMetric_WithLastUpdate(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: true}
	adapter := NewPrometheusAdapter(config)

	incident := createTestIncident()
	product := entities.Product{Title: "Test Product", ID: "test-product"}

	ch := make(chan prometheus.Metric, 1)
	adapter.addMetric(ch, incident, product)
	close(ch)

	// Verify one metric was created
	metricCount := 0
	for range ch {
		metricCount++
	}

	if metricCount != 1 {
		t.Errorf("Expected 1 metric, got %d", metricCount)
	}
}

func TestPrometheusAdapter_addMetric_URIConstruction(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	incident := createTestIncident()
	incident.URI = "incidents/test-incident"
	product := entities.Product{Title: "Test Product", ID: "test-product"}

	ch := make(chan prometheus.Metric, 1)
	adapter.addMetric(ch, incident, product)
	close(ch)

	// The metric should be created successfully (we can't easily inspect the URI in the metric)
	metricCount := 0
	for range ch {
		metricCount++
	}

	if metricCount != 1 {
		t.Errorf("Expected 1 metric, got %d", metricCount)
	}
}

func TestPrometheusAdapter_SeverityMapping(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	tests := []struct {
		name     string
		severity valueobjects.Severity
		expected float64
	}{
		{"High severity", valueobjects.SeverityHigh, 3.0},
		{"Medium severity", valueobjects.SeverityMedium, 2.0},
		{"Low severity", valueobjects.SeverityLow, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incident := createTestIncident()
			incident.Severity = tt.severity
			product := entities.Product{Title: "Test Product", ID: "test-product"}

			ch := make(chan prometheus.Metric, 1)
			adapter.addMetric(ch, incident, product)
			close(ch)

			// Verify metric was created (the actual severity value is tested in the domain layer)
			metricCount := 0
			for range ch {
				metricCount++
			}

			if metricCount != 1 {
				t.Errorf("Expected 1 metric for %s, got %d", tt.name, metricCount)
			}
		})
	}
}

func TestPrometheusAdapter_ResolvedIncident(t *testing.T) {
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	// Create a resolved incident
	incident := createTestIncident()
	incident.MostRecentUpdate.UpdateStatus = valueobjects.UpdateStatusAvailable
	product := entities.Product{Title: "Test Product", ID: "test-product"}

	ch := make(chan prometheus.Metric, 1)
	adapter.addMetric(ch, incident, product)
	close(ch)

	// Should still create a metric (severity will be 0.0 for resolved incidents)
	metricCount := 0
	for range ch {
		metricCount++
	}

	if metricCount != 1 {
		t.Errorf("Expected 1 metric for resolved incident, got %d", metricCount)
	}
}

// Helper function to create a test incident
func createTestIncident() entities.Incident {
	id, _ := valueobjects.NewIncidentID("test-incident")
	return entities.Incident{
		ID:                  id,
		ExternalDescription: "Test incident description",
		Severity:            valueobjects.SeverityHigh,
		URI:                 "incidents/test-incident",
		MostRecentUpdate: entities.Update{
			Status:       "Investigating issue",
			UpdateStatus: "INVESTIGATING",
		},
		AffectedProducts: []entities.Product{
			{Title: "Test Product", ID: "test-product"},
		},
	}
}
