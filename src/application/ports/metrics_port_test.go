package ports

import (
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
)

// Test MetricsConfig struct
func TestMetricsConfig(t *testing.T) {
	tests := []struct {
		name           string
		saveLastUpdate bool
	}{
		{
			name:           "config with SaveLastUpdate true",
			saveLastUpdate: true,
		},
		{
			name:           "config with SaveLastUpdate false",
			saveLastUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MetricsConfig{
				SaveLastUpdate: tt.saveLastUpdate,
			}

			if config.SaveLastUpdate != tt.saveLastUpdate {
				t.Errorf("SaveLastUpdate = %v, want %v", config.SaveLastUpdate, tt.saveLastUpdate)
			}
		})
	}
}

// Mock implementation for testing interface compliance
type mockMetricsPort struct {
	incidents []entities.Incident
	err       error
}

func (m *mockMetricsPort) CollectMetrics(incidents []entities.Incident, config MetricsConfig) error {
	m.incidents = incidents
	return m.err
}

func TestMetricsPortInterface(t *testing.T) {
	// Test that our mock implements the interface
	var port MetricsPort = &mockMetricsPort{}

	config := MetricsConfig{SaveLastUpdate: true}
	incidents := []entities.Incident{}

	err := port.CollectMetrics(incidents, config)
	if err != nil {
		t.Errorf("CollectMetrics() returned unexpected error: %v", err)
	}

	// Test with mock that returns error
	mockWithError := &mockMetricsPort{err: &testError{"test error"}}
	var portWithError MetricsPort = mockWithError

	err = portWithError.CollectMetrics(incidents, config)
	if err == nil {
		t.Error("Expected error from CollectMetrics() but got nil")
	}

	if err.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", err.Error())
	}
}

func TestMetricsConfigZeroValue(t *testing.T) {
	var config MetricsConfig

	if config.SaveLastUpdate != false {
		t.Errorf("Zero value SaveLastUpdate should be false, got %v", config.SaveLastUpdate)
	}
}

func TestMetricsConfigEquality(t *testing.T) {
	config1 := MetricsConfig{SaveLastUpdate: true}
	config2 := MetricsConfig{SaveLastUpdate: true}
	config3 := MetricsConfig{SaveLastUpdate: false}

	if config1 != config2 {
		t.Error("Configs with same values should be equal")
	}

	if config1 == config3 {
		t.Error("Configs with different values should not be equal")
	}
}

// Helper error type for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}