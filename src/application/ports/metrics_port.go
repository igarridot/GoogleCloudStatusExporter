package ports

import "github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"

type MetricsPort interface {
	CollectMetrics(incidents []entities.Incident, config MetricsConfig) error
}

type MetricsConfig struct {
	SaveLastUpdate bool
}
