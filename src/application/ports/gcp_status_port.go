package ports

import "github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"

type GCPStatusPort interface {
	GetIncidents() ([]entities.Incident, error)
}
