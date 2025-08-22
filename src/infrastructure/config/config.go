package config

import (
	"flag"
	"os"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

type Config struct {
	ListenAddress   string
	MetricsPath     string
	SaveLastUpdate  bool
	CollectResolved bool
	FilterCriteria  valueobjects.FilterCriteria
}

func LoadConfig() *Config {
	var (
		listenAddress         = flag.String("web.listen-address", ":9118", "Address to listen on for web interface.")
		metricPath            = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
		lastUpdate            = flag.Bool("exporter.save-last-update", false, "Use flag if you want to save last incident update description. Disabled by default.")
		collectResolvedEvents = flag.Bool("exporter.collect-resolved-incidents", false, "Use flag if you want to collect already resolved incidents. Disabled by default.")
		incidentsZones        = flag.String("exporter.incidents-zones", "", "Use this flag if you want to filter by incident zones.")
		filteredProducts      = flag.String("exporter.filtered-products", "", "Use this flag if you want to filter incidents by product names.")
	)
	flag.Parse()

	if *lastUpdate {
		os.Setenv("SAVE_LAST_UPDATE", "true")
	}
	if *collectResolvedEvents {
		os.Setenv("COLLECT_RESOLVED_EVENTS", "true")
	}
	if *incidentsZones != "" {
		os.Setenv("INCIDENTS_ZONES", *incidentsZones)
	}
	if *filteredProducts != "" {
		os.Setenv("FILTERED_PRODUCTS", *filteredProducts)
	}

	return &Config{
		ListenAddress:   *listenAddress,
		MetricsPath:     *metricPath,
		SaveLastUpdate:  *lastUpdate,
		CollectResolved: *collectResolvedEvents,
		FilterCriteria:  valueobjects.NewFilterCriteria(*incidentsZones, *filteredProducts),
	}
}

func (c *Config) ToMetricsConfig() ports.MetricsConfig {
	return ports.MetricsConfig{
		SaveLastUpdate: c.SaveLastUpdate,
	}
}
