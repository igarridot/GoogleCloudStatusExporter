package container

import (
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/usecases"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/services"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/infrastructure/adapters/http"
	prometheusAdapter "github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/infrastructure/adapters/prometheus"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/infrastructure/config"
	"github.com/prometheus/client_golang/prometheus"
)

type Container struct {
	config    *config.Config
	collector prometheus.Collector
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: cfg,
	}
}

func (c *Container) GetCollector() prometheus.Collector {
	if c.collector == nil {
		httpClient := &http.DefaultHTTPClient{}
		gcpStatusAdapter := http.NewGCPStatusAdapter(httpClient)

		metricsConfig := c.config.ToMetricsConfig()
		promAdapter := prometheusAdapter.NewPrometheusAdapter(metricsConfig)

		filterService := services.NewIncidentFilterService()

		useCase := usecases.NewCollectMetricsUseCase(
			gcpStatusAdapter,
			promAdapter,
			filterService,
		)

		c.collector = prometheusAdapter.NewGCPStatusCollector(
			useCase,
			promAdapter,
			c.config.FilterCriteria,
			c.config.CollectResolved,
			metricsConfig,
		)
	}

	return c.collector
}
