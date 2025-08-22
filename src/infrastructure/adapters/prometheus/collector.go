package prometheus

import (
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/usecases"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type GCPStatusCollector struct {
	useCase         *usecases.CollectMetricsUseCase
	adapter         *PrometheusAdapter
	filterCriteria  valueobjects.FilterCriteria
	collectResolved bool
	metricsConfig   ports.MetricsConfig
}

func NewGCPStatusCollector(
	useCase *usecases.CollectMetricsUseCase,
	adapter *PrometheusAdapter,
	filterCriteria valueobjects.FilterCriteria,
	collectResolved bool,
	metricsConfig ports.MetricsConfig,
) *GCPStatusCollector {
	return &GCPStatusCollector{
		useCase:         useCase,
		adapter:         adapter,
		filterCriteria:  filterCriteria,
		collectResolved: collectResolved,
		metricsConfig:   metricsConfig,
	}
}

func (c *GCPStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	c.adapter.Describe(ch)
}

func (c *GCPStatusCollector) Collect(ch chan<- prometheus.Metric) {
	incidents, err := c.useCase.Execute(c.filterCriteria, c.collectResolved, c.metricsConfig)
	if err != nil {
		log.Printf("Error collecting metrics: %v", err)
		return
	}

	c.adapter.CollectMetricsWithChannel(incidents, ch)
}
