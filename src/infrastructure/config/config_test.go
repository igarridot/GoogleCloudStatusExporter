package config

import (
	"flag"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name                    string
		args                    []string
		expectedListenAddress   string
		expectedMetricsPath     string
		expectedSaveLastUpdate  bool
		expectedCollectResolved bool
	}{
		{
			name:                    "Default configuration",
			args:                    []string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
		},
		{
			name:                  "Custom configuration",
			args:                  []string{"-web.listen-address", ":8080", "-web.metrics-path", "/custom"},
			expectedListenAddress: ":8080",
			expectedMetricsPath:   "/custom",
		},
		{
			name:                    "Enable flags",
			args:                    []string{"-exporter.save-last-update", "-exporter.collect-resolved-incidents"},
			expectedSaveLastUpdate:  true,
			expectedCollectResolved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set command line arguments
			os.Args = append([]string{"test"}, tt.args...)

			config := LoadConfig()

			if tt.expectedListenAddress != "" && config.ListenAddress != tt.expectedListenAddress {
				t.Errorf("Expected ListenAddress %s, got %s", tt.expectedListenAddress, config.ListenAddress)
			}
			if tt.expectedMetricsPath != "" && config.MetricsPath != tt.expectedMetricsPath {
				t.Errorf("Expected MetricsPath %s, got %s", tt.expectedMetricsPath, config.MetricsPath)
			}
			if config.SaveLastUpdate != tt.expectedSaveLastUpdate {
				t.Errorf("Expected SaveLastUpdate %v, got %v", tt.expectedSaveLastUpdate, config.SaveLastUpdate)
			}
			if config.CollectResolved != tt.expectedCollectResolved {
				t.Errorf("Expected CollectResolved %v, got %v", tt.expectedCollectResolved, config.CollectResolved)
			}
		})
	}
}

func TestConfigToMetricsConfig(t *testing.T) {
	config := &Config{SaveLastUpdate: true}
	metricsConfig := config.ToMetricsConfig()

	if !metricsConfig.SaveLastUpdate {
		t.Error("Expected SaveLastUpdate to be true")
	}
}
