package prometheus

import (
	"errors"
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/prometheus/client_golang/prometheus"
)

// UseCase interface for testing
type UseCase interface {
	Execute(filterCriteria valueobjects.FilterCriteria, collectResolved bool, config ports.MetricsConfig) ([]entities.Incident, error)
}

// Mock use case for testing
type mockCollectMetricsUseCase struct {
	incidents []entities.Incident
	err       error
}

func (m *mockCollectMetricsUseCase) Execute(filterCriteria valueobjects.FilterCriteria, collectResolved bool, config ports.MetricsConfig) ([]entities.Incident, error) {
	return m.incidents, m.err
}

// Test collector with interface
type testGCPStatusCollector struct {
	useCase         UseCase
	adapter         *PrometheusAdapter
	filterCriteria  valueobjects.FilterCriteria
	collectResolved bool
	metricsConfig   ports.MetricsConfig
}

func newTestGCPStatusCollector(
	useCase UseCase,
	adapter *PrometheusAdapter,
	filterCriteria valueobjects.FilterCriteria,
	collectResolved bool,
	metricsConfig ports.MetricsConfig,
) *testGCPStatusCollector {
	return &testGCPStatusCollector{
		useCase:         useCase,
		adapter:         adapter,
		filterCriteria:  filterCriteria,
		collectResolved: collectResolved,
		metricsConfig:   metricsConfig,
	}
}

func (c *testGCPStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	c.adapter.Describe(ch)
}

func (c *testGCPStatusCollector) Collect(ch chan<- prometheus.Metric) {
	incidents, err := c.useCase.Execute(c.filterCriteria, c.collectResolved, c.metricsConfig)
	if err != nil {
		return
	}
	c.adapter.CollectMetricsWithChannel(incidents, ch)
}

func TestNewTestGCPStatusCollector(t *testing.T) {
	useCase := &mockCollectMetricsUseCase{}
	adapter := NewPrometheusAdapter(ports.MetricsConfig{SaveLastUpdate: false})
	filterCriteria := valueobjects.FilterCriteria{}
	collectResolved := true
	metricsConfig := ports.MetricsConfig{SaveLastUpdate: false}

	collector := newTestGCPStatusCollector(useCase, adapter, filterCriteria, collectResolved, metricsConfig)

	if collector == nil {
		t.Fatal("newTestGCPStatusCollector() returned nil")
	}

	if collector.useCase != useCase {
		t.Error("useCase not set correctly")
	}

	if collector.adapter != adapter {
		t.Error("adapter not set correctly")
	}

	if collector.collectResolved != collectResolved {
		t.Error("collectResolved not set correctly")
	}

	if collector.metricsConfig.SaveLastUpdate != metricsConfig.SaveLastUpdate {
		t.Error("metricsConfig not set correctly")
	}
}

func TestGCPStatusCollector_Describe(t *testing.T) {
	useCase := &mockCollectMetricsUseCase{}
	adapter := NewPrometheusAdapter(ports.MetricsConfig{SaveLastUpdate: false})
	filterCriteria := valueobjects.FilterCriteria{}
	collectResolved := false
	metricsConfig := ports.MetricsConfig{SaveLastUpdate: false}

	collector := newTestGCPStatusCollector(useCase, adapter, filterCriteria, collectResolved, metricsConfig)

	ch := make(chan *prometheus.Desc, 1)
	collector.Describe(ch)
	close(ch)

	descriptors := make([]*prometheus.Desc, 0)
	for desc := range ch {
		descriptors = append(descriptors, desc)
	}

	if len(descriptors) != 1 {
		t.Errorf("Expected 1 descriptor, got %d", len(descriptors))
	}
}

func TestGCPStatusCollector_Collect_Success(t *testing.T) {
	incidents := createTestIncidents()
	useCase := &mockCollectMetricsUseCase{
		incidents: incidents,
		err:       nil,
	}
	adapter := NewPrometheusAdapter(ports.MetricsConfig{SaveLastUpdate: false})
	filterCriteria := valueobjects.FilterCriteria{}
	collectResolved := true
	metricsConfig := ports.MetricsConfig{SaveLastUpdate: false}

	collector := newTestGCPStatusCollector(useCase, adapter, filterCriteria, collectResolved, metricsConfig)

	// Create a buffered channel to collect metrics
	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	expectedMetrics := len(incidents) // One metric per incident (each has one product)
	if len(metrics) != expectedMetrics {
		t.Errorf("Expected %d metrics, got %d", expectedMetrics, len(metrics))
	}
}

func TestGCPStatusCollector_Collect_Error(t *testing.T) {
	useCase := &mockCollectMetricsUseCase{
		incidents: nil,
		err:       errors.New("use case error"),
	}
	adapter := NewPrometheusAdapter(ports.MetricsConfig{SaveLastUpdate: false})
	filterCriteria := valueobjects.FilterCriteria{}
	collectResolved := true
	metricsConfig := ports.MetricsConfig{SaveLastUpdate: false}

	collector := newTestGCPStatusCollector(useCase, adapter, filterCriteria, collectResolved, metricsConfig)

	// Create a buffered channel to collect metrics
	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	// Should not collect any metrics when use case returns error
	if len(metrics) != 0 {
		t.Errorf("Expected 0 metrics when use case fails, got %d", len(metrics))
	}
}

func TestGCPStatusCollector_Integration(t *testing.T) {
	tests := []struct {
		name            string
		incidents       []entities.Incident
		collectResolved bool
		saveLastUpdate  bool
		expectedMetrics int
	}{
		{
			name:            "collect resolved and unresolved incidents",
			incidents:       createTestIncidents(),
			collectResolved: true,
			saveLastUpdate:  false,
			expectedMetrics: 2,
		},
		{
			name:            "incidents with save last update enabled",
			incidents:       createTestIncidents()[:1],
			collectResolved: false,
			saveLastUpdate:  true,
			expectedMetrics: 1,
		},
		{
			name:            "empty incidents list",
			incidents:       []entities.Incident{},
			collectResolved: true,
			saveLastUpdate:  false,
			expectedMetrics: 0,
		},
		{
			name:            "incident with multiple products",
			incidents:       createIncidentWithMultipleProducts(),
			collectResolved: true,
			saveLastUpdate:  false,
			expectedMetrics: 2, // Two products = two metrics
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &mockCollectMetricsUseCase{
				incidents: tt.incidents,
				err:       nil,
			}
			adapter := NewPrometheusAdapter(ports.MetricsConfig{SaveLastUpdate: tt.saveLastUpdate})
			filterCriteria := valueobjects.FilterCriteria{}
			metricsConfig := ports.MetricsConfig{SaveLastUpdate: tt.saveLastUpdate}

			collector := newTestGCPStatusCollector(useCase, adapter, filterCriteria, tt.collectResolved, metricsConfig)

			// Register collector with Prometheus registry to test end-to-end
			registry := prometheus.NewRegistry()
			registry.MustRegister(collector)

			metricFamilies, err := registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
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

func TestGCPStatusCollector_UseCaseParameters(t *testing.T) {
	// Test that the collector passes the correct parameters to the use case
	var capturedFilterCriteria valueobjects.FilterCriteria
	var capturedCollectResolved bool
	var capturedConfig ports.MetricsConfig

	mockUseCase := &useCaseParameterCapture{
		onExecute: func(filterCriteria valueobjects.FilterCriteria, collectResolved bool, config ports.MetricsConfig) ([]entities.Incident, error) {
			capturedFilterCriteria = filterCriteria
			capturedCollectResolved = collectResolved
			capturedConfig = config
			return []entities.Incident{}, nil
		},
	}

	adapter := NewPrometheusAdapter(ports.MetricsConfig{SaveLastUpdate: true})
	filterCriteria := valueobjects.FilterCriteria{
		Zones:    []string{"us-central1"},
		Products: []string{"compute"},
	}
	collectResolved := true
	metricsConfig := ports.MetricsConfig{SaveLastUpdate: true}

	collector := newTestGCPStatusCollector(mockUseCase, adapter, filterCriteria, collectResolved, metricsConfig)

	ch := make(chan prometheus.Metric, 1)
	collector.Collect(ch)
	close(ch)

	// Verify parameters were passed correctly
	if len(capturedFilterCriteria.Zones) != 1 || capturedFilterCriteria.Zones[0] != "us-central1" {
		t.Errorf("FilterCriteria.Zones not passed correctly: %v", capturedFilterCriteria.Zones)
	}

	if len(capturedFilterCriteria.Products) != 1 || capturedFilterCriteria.Products[0] != "compute" {
		t.Errorf("FilterCriteria.Products not passed correctly: %v", capturedFilterCriteria.Products)
	}

	if capturedCollectResolved != collectResolved {
		t.Errorf("collectResolved not passed correctly: %v", capturedCollectResolved)
	}

	if capturedConfig.SaveLastUpdate != metricsConfig.SaveLastUpdate {
		t.Errorf("metricsConfig not passed correctly: %v", capturedConfig)
	}
}

// Helper for capturing use case parameters
type useCaseParameterCapture struct {
	onExecute func(valueobjects.FilterCriteria, bool, ports.MetricsConfig) ([]entities.Incident, error)
}

func (u *useCaseParameterCapture) Execute(filterCriteria valueobjects.FilterCriteria, collectResolved bool, config ports.MetricsConfig) ([]entities.Incident, error) {
	return u.onExecute(filterCriteria, collectResolved, config)
}