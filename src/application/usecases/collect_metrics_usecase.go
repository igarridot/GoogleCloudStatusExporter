package usecases

import (
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/services"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

type CollectMetricsUseCase struct {
	gcpStatusPort ports.GCPStatusPort
	metricsPort   ports.MetricsPort
	filterService *services.IncidentFilterService
}

func NewCollectMetricsUseCase(
	gcpStatusPort ports.GCPStatusPort,
	metricsPort ports.MetricsPort,
	filterService *services.IncidentFilterService,
) *CollectMetricsUseCase {
	return &CollectMetricsUseCase{
		gcpStatusPort: gcpStatusPort,
		metricsPort:   metricsPort,
		filterService: filterService,
	}
}

func (uc *CollectMetricsUseCase) Execute(
	filterCriteria valueobjects.FilterCriteria,
	collectResolved bool,
	metricsConfig ports.MetricsConfig,
) ([]entities.Incident, error) {
	incidents, err := uc.gcpStatusPort.GetIncidents()
	if err != nil {
		return nil, err
	}

	filteredIncidents := uc.filterService.FilterIncidents(incidents, filterCriteria, collectResolved)

	return filteredIncidents, nil
}
