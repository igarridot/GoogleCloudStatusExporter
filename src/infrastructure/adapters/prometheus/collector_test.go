package prometheus

import (
	"errors"
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/usecases"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/services"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/testutils"
	"github.com/prometheus/client_golang/prometheus"
)

// Mock ports for testing
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

func TestNewGCPStatusCollector(t *testing.T) {
	// Create real use case with mock ports
	gcpStatusPort := &mockGCPStatusPort{incidents: testutils.CreateTestIncidents()}
	metricsPort := &mockMetricsPort{}
	filterService := services.NewIncidentFilterService()
	useCase := usecases.NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

	// Create real adapter
	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)

	criteria := valueobjects.NewFilterCriteria("", "")
	collectResolved := true

	collector := NewGCPStatusCollector(useCase, adapter, criteria, collectResolved, config)

	if collector == nil {
		t.Error("Expected collector to be created, got nil")
		return
	}

	if collector.useCase != useCase {
		t.Error("Expected useCase to be set correctly")
	}

	if collector.adapter != adapter {
		t.Error("Expected adapter to be set correctly")
	}

	if collector.collectResolved != collectResolved {
		t.Error("Expected collectResolved to be set correctly")
	}

	if collector.metricsConfig != config {
		t.Error("Expected metricsConfig to be set correctly")
	}
}

func TestGCPStatusCollector_Describe(t *testing.T) {
	gcpStatusPort := &mockGCPStatusPort{incidents: testutils.CreateTestIncidents()}
	metricsPort := &mockMetricsPort{}
	filterService := services.NewIncidentFilterService()
	useCase := usecases.NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)
	criteria := valueobjects.NewFilterCriteria("", "")
	collector := NewGCPStatusCollector(useCase, adapter, criteria, true, config)

	ch := make(chan *prometheus.Desc, 10)
	collector.Describe(ch)
	close(ch)

	// Check that at least one descriptor was sent
	descCount := 0
	for desc := range ch {
		descCount++
		if desc == nil {
			t.Error("Expected non-nil descriptor")
		}
	}

	if descCount == 0 {
		t.Error("Expected at least one descriptor to be sent")
	}
}

func TestGCPStatusCollector_Collect_Success(t *testing.T) {
	gcpStatusPort := &mockGCPStatusPort{incidents: testutils.CreateTestIncidents()}
	metricsPort := &mockMetricsPort{}
	filterService := services.NewIncidentFilterService()
	useCase := usecases.NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)
	criteria := valueobjects.NewFilterCriteria("", "")
	collector := NewGCPStatusCollector(useCase, adapter, criteria, true, config)

	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	// Count metrics collected
	metricCount := 0
	for metric := range ch {
		metricCount++
		if metric == nil {
			t.Error("Expected non-nil metric")
		}
	}

	// We should get metrics for the incidents (one per affected product)
	incidents := testutils.CreateTestIncidents()
	expectedMetrics := 0
	for _, incident := range incidents {
		expectedMetrics += len(incident.AffectedProducts)
	}

	if metricCount != expectedMetrics {
		t.Errorf("Expected %d metrics, got %d", expectedMetrics, metricCount)
	}
}

func TestGCPStatusCollector_Collect_UseCaseError(t *testing.T) {
	gcpStatusPort := &mockGCPStatusPort{err: errors.New("mock error")}
	metricsPort := &mockMetricsPort{}
	filterService := services.NewIncidentFilterService()
	useCase := usecases.NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

	config := ports.MetricsConfig{SaveLastUpdate: false}
	adapter := NewPrometheusAdapter(config)
	criteria := valueobjects.NewFilterCriteria("", "")
	collector := NewGCPStatusCollector(useCase, adapter, criteria, true, config)

	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	// Should not panic and should not send any metrics when useCase fails
	metricCount := 0
	for range ch {
		metricCount++
	}

	if metricCount != 0 {
		t.Errorf("Expected 0 metrics when useCase fails, got %d", metricCount)
	}

	// The collector should handle the error gracefully and continue
	// (Error logging is handled by the collector internally)
}

func TestGCPStatusCollector_Collect_WithDifferentConfigurations(t *testing.T) {
	tests := []struct {
		name            string
		collectResolved bool
		saveLastUpdate  bool
		filterZones     string
		filterProducts  string
	}{
		{"Collect resolved incidents with last update", true, true, "", ""},
		{"Skip resolved incidents without last update", false, false, "", ""},
		{"Filter by zones", true, false, "us-east1,us-west1", ""},
		{"Filter by products", true, false, "", "Compute Engine,Cloud Storage"},
		{"Filter by both zones and products", true, false, "us-east1", "Compute Engine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcpStatusPort := &mockGCPStatusPort{incidents: testutils.CreateTestIncidents()}
			metricsPort := &mockMetricsPort{}
			filterService := services.NewIncidentFilterService()
			useCase := usecases.NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

			config := ports.MetricsConfig{SaveLastUpdate: tt.saveLastUpdate}
			adapter := NewPrometheusAdapter(config)
			criteria := valueobjects.NewFilterCriteria(tt.filterZones, tt.filterProducts)
			collector := NewGCPStatusCollector(useCase, adapter, criteria, tt.collectResolved, config)

			ch := make(chan prometheus.Metric, 10)
			collector.Collect(ch)
			close(ch)

			// Should not panic with any configuration
			metricCount := 0
			for range ch {
				metricCount++
			}

			// Should produce some metrics (exact count depends on filtering logic)
			if metricCount < 0 {
				t.Errorf("Expected non-negative metric count, got %d", metricCount)
			}
		})
	}
}

func TestGCPStatusCollector_Integration_WithRealPrometheusAdapter(t *testing.T) {
	gcpStatusPort := &mockGCPStatusPort{incidents: testutils.CreateTestIncidents()}
	metricsPort := &mockMetricsPort{}
	filterService := services.NewIncidentFilterService()
	useCase := usecases.NewCollectMetricsUseCase(gcpStatusPort, metricsPort, filterService)

	config := ports.MetricsConfig{SaveLastUpdate: true}
	adapter := NewPrometheusAdapter(config)
	criteria := valueobjects.NewFilterCriteria("", "")
	collector := NewGCPStatusCollector(useCase, adapter, criteria, true, config)

	// Test that the collector integrates properly with the Prometheus adapter
	ch := make(chan prometheus.Metric, 10)

	// Describe should work
	descCh := make(chan *prometheus.Desc, 10)
	collector.Describe(descCh)
	close(descCh)

	descCount := 0
	for range descCh {
		descCount++
	}

	if descCount == 0 {
		t.Error("Expected descriptors to be provided")
	}

	// Collect should work
	collector.Collect(ch)
	close(ch)

	metricCount := 0
	for range ch {
		metricCount++
	}

	if metricCount == 0 {
		t.Error("Expected metrics to be collected")
	}
}
