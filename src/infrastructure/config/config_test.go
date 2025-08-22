package config

import (
	"flag"
	"os"
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name                    string
		args                    []string
		envVars                 map[string]string
		expectedListenAddress   string
		expectedMetricsPath     string
		expectedSaveLastUpdate  bool
		expectedCollectResolved bool
		expectedZones           []string
		expectedProducts        []string
	}{
		{
			name:                    "Default configuration",
			args:                    []string{},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Custom listen address and metrics path",
			args: []string{"-web.listen-address", ":8080", "-web.metrics-path", "/custom"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":8080",
			expectedMetricsPath:     "/custom",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Enable save last update",
			args: []string{"-exporter.save-last-update"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  true,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Enable collect resolved incidents",
			args: []string{"-exporter.collect-resolved-incidents"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: true,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Both flags enabled",
			args: []string{"-exporter.save-last-update", "-exporter.collect-resolved-incidents"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  true,
			expectedCollectResolved: true,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Filter by single zone",
			args: []string{"-exporter.incidents-zones", "us-east1"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{"us-east1"},
			expectedProducts:        []string{},
		},
		{
			name: "Filter by multiple zones",
			args: []string{"-exporter.incidents-zones", "us-east1,europe-west1,asia-east1"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{"us-east1", "europe-west1", "asia-east1"},
			expectedProducts:        []string{},
		},
		{
			name: "Filter by single product",
			args: []string{"-exporter.filtered-products", "Apigee"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{"Apigee"},
		},
		{
			name: "Filter by multiple products",
			args: []string{"-exporter.filtered-products", "Apigee,Compute Engine,Cloud Storage"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{"Apigee", "Compute Engine", "Cloud Storage"},
		},
		{
			name: "Filter by zones and products",
			args: []string{"-exporter.incidents-zones", "us-east1,europe-west1", "-exporter.filtered-products", "Apigee,Compute"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{"us-east1", "europe-west1"},
			expectedProducts:        []string{"Apigee", "Compute"},
		},
		{
			name: "All flags enabled with filtering",
			args: []string{
				"-web.listen-address", ":9999",
				"-web.metrics-path", "/custom-metrics",
				"-exporter.save-last-update",
				"-exporter.collect-resolved-incidents",
				"-exporter.incidents-zones", "us-central1",
				"-exporter.filtered-products", "BigQuery",
			},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9999",
			expectedMetricsPath:     "/custom-metrics",
			expectedSaveLastUpdate:  true,
			expectedCollectResolved: true,
			expectedZones:           []string{"us-central1"},
			expectedProducts:        []string{"BigQuery"},
		},
		{
			name: "Empty zone filter",
			args: []string{"-exporter.incidents-zones", ""},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Empty product filter",
			args: []string{"-exporter.filtered-products", ""},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{},
		},
		{
			name: "Zone filter with spaces",
			args: []string{"-exporter.incidents-zones", "us-east1, europe-west1 , asia-east1"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{"us-east1", " europe-west1 ", " asia-east1"},
			expectedProducts:        []string{},
		},
		{
			name: "Product filter with spaces",
			args: []string{"-exporter.filtered-products", "Compute Engine, Cloud Storage , Apigee"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{},
			expectedProducts:        []string{"Compute Engine", " Cloud Storage ", " Apigee"},
		},
		{
			name: "Single character values",
			args: []string{"-exporter.incidents-zones", "a", "-exporter.filtered-products", "b"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{"a"},
			expectedProducts:        []string{"b"},
		},
		{
			name: "Special characters in filters",
			args: []string{"-exporter.incidents-zones", "us-east1-a", "-exporter.filtered-products", "Google Cloud Storage"},
			envVars:                 map[string]string{},
			expectedListenAddress:   ":9118",
			expectedMetricsPath:     "/metrics",
			expectedSaveLastUpdate:  false,
			expectedCollectResolved: false,
			expectedZones:           []string{"us-east1-a"},
			expectedProducts:        []string{"Google Cloud Storage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("SAVE_LAST_UPDATE")
			os.Unsetenv("COLLECT_RESOLVED_EVENTS")
			os.Unsetenv("INCIDENTS_ZONES")
			os.Unsetenv("FILTERED_PRODUCTS")

			// Set environment variables for this test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set command line arguments
			os.Args = append([]string{"test"}, tt.args...)

			config := LoadConfig()

			// Validate basic configuration
			if config.ListenAddress != tt.expectedListenAddress {
				t.Errorf("Expected ListenAddress %s, got %s", tt.expectedListenAddress, config.ListenAddress)
			}
			if config.MetricsPath != tt.expectedMetricsPath {
				t.Errorf("Expected MetricsPath %s, got %s", tt.expectedMetricsPath, config.MetricsPath)
			}
			if config.SaveLastUpdate != tt.expectedSaveLastUpdate {
				t.Errorf("Expected SaveLastUpdate %v, got %v", tt.expectedSaveLastUpdate, config.SaveLastUpdate)
			}
			if config.CollectResolved != tt.expectedCollectResolved {
				t.Errorf("Expected CollectResolved %v, got %v", tt.expectedCollectResolved, config.CollectResolved)
			}

			// Validate filter criteria
			if len(config.FilterCriteria.Zones) != len(tt.expectedZones) {
				t.Errorf("Expected %d zones, got %d", len(tt.expectedZones), len(config.FilterCriteria.Zones))
			}
			for i, expectedZone := range tt.expectedZones {
				if i >= len(config.FilterCriteria.Zones) || config.FilterCriteria.Zones[i] != expectedZone {
					t.Errorf("Expected zone[%d] %s, got %s", i, expectedZone, config.FilterCriteria.Zones[i])
				}
			}

			if len(config.FilterCriteria.Products) != len(tt.expectedProducts) {
				t.Errorf("Expected %d products, got %d", len(tt.expectedProducts), len(config.FilterCriteria.Products))
			}
			for i, expectedProduct := range tt.expectedProducts {
				if i >= len(config.FilterCriteria.Products) || config.FilterCriteria.Products[i] != expectedProduct {
					t.Errorf("Expected product[%d] %s, got %s", i, expectedProduct, config.FilterCriteria.Products[i])
				}
			}

			// Validate environment variables are set correctly
			if tt.expectedSaveLastUpdate {
				if os.Getenv("SAVE_LAST_UPDATE") != "true" {
					t.Error("Expected SAVE_LAST_UPDATE env var to be set to 'true'")
				}
			}
			if tt.expectedCollectResolved {
				if os.Getenv("COLLECT_RESOLVED_EVENTS") != "true" {
					t.Error("Expected COLLECT_RESOLVED_EVENTS env var to be set to 'true'")
				}
			}
			if len(tt.expectedZones) > 0 {
				expectedZoneEnv := ""
				for i, zone := range tt.expectedZones {
					if i > 0 {
						expectedZoneEnv += ","
					}
					expectedZoneEnv += zone
				}
				if os.Getenv("INCIDENTS_ZONES") != expectedZoneEnv {
					t.Errorf("Expected INCIDENTS_ZONES env var to be '%s', got '%s'", expectedZoneEnv, os.Getenv("INCIDENTS_ZONES"))
				}
			}
			if len(tt.expectedProducts) > 0 {
				expectedProductEnv := ""
				for i, product := range tt.expectedProducts {
					if i > 0 {
						expectedProductEnv += ","
					}
					expectedProductEnv += product
				}
				if os.Getenv("FILTERED_PRODUCTS") != expectedProductEnv {
					t.Errorf("Expected FILTERED_PRODUCTS env var to be '%s', got '%s'", expectedProductEnv, os.Getenv("FILTERED_PRODUCTS"))
				}
			}

			// Clean up environment variables
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestConfigToMetricsConfig(t *testing.T) {
	tests := []struct {
		name               string
		saveLastUpdate     bool
		expectedSaveUpdate bool
	}{
		{
			name:               "SaveLastUpdate enabled",
			saveLastUpdate:     true,
			expectedSaveUpdate: true,
		},
		{
			name:               "SaveLastUpdate disabled",
			saveLastUpdate:     false,
			expectedSaveUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				SaveLastUpdate: tt.saveLastUpdate,
			}

			metricsConfig := config.ToMetricsConfig()

			if metricsConfig.SaveLastUpdate != tt.expectedSaveUpdate {
				t.Errorf("Expected SaveLastUpdate %v, got %v", tt.expectedSaveUpdate, metricsConfig.SaveLastUpdate)
			}
		})
	}
}

func TestFilterCriteriaHelpers(t *testing.T) {
	tests := []struct {
		name                string
		zones               string
		products            string
		expectedHasZone     bool
		expectedHasProduct  bool
	}{
		{
			name:                "No filters",
			zones:               "",
			products:            "",
			expectedHasZone:     false,
			expectedHasProduct:  false,
		},
		{
			name:                "Zone filter only",
			zones:               "us-east1",
			products:            "",
			expectedHasZone:     true,
			expectedHasProduct:  false,
		},
		{
			name:                "Product filter only",
			zones:               "",
			products:            "Apigee",
			expectedHasZone:     false,
			expectedHasProduct:  true,
		},
		{
			name:                "Both filters",
			zones:               "us-east1",
			products:            "Apigee",
			expectedHasZone:     true,
			expectedHasProduct:  true,
		},
		{
			name:                "Multiple zones and products",
			zones:               "us-east1,europe-west1",
			products:            "Apigee,Compute",
			expectedHasZone:     true,
			expectedHasProduct:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := valueobjects.NewFilterCriteria(tt.zones, tt.products)

			if criteria.HasZoneFilter() != tt.expectedHasZone {
				t.Errorf("Expected HasZoneFilter %v, got %v", tt.expectedHasZone, criteria.HasZoneFilter())
			}
			if criteria.HasProductFilter() != tt.expectedHasProduct {
				t.Errorf("Expected HasProductFilter %v, got %v", tt.expectedHasProduct, criteria.HasProductFilter())
			}
		})
	}
}