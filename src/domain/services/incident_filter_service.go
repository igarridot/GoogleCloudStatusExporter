package services

import (
	"strings"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

type IncidentFilterService struct{}

func NewIncidentFilterService() *IncidentFilterService {
	return &IncidentFilterService{}
}

func (s *IncidentFilterService) FilterIncidents(incidents []entities.Incident, criteria valueobjects.FilterCriteria, collectResolved bool) []entities.Incident {
	var filtered []entities.Incident

	for _, incident := range incidents {
		if !collectResolved && incident.IsResolved() {
			continue
		}

		if s.matchesZoneFilter(incident, criteria) && s.matchesProductFilter(incident, criteria) {
			filtered = append(filtered, incident)
		}
	}

	return filtered
}

func (s *IncidentFilterService) matchesZoneFilter(incident entities.Incident, criteria valueobjects.FilterCriteria) bool {
	if !criteria.HasZoneFilter() {
		return true
	}

	specialZones := []string{"Global", "global", "various", "Various"}
	allZones := append(criteria.Zones, specialZones...)

	for _, zone := range allZones {
		if strings.Contains(incident.ExternalDescription, zone) {
			return true
		}
	}

	return false
}

func (s *IncidentFilterService) matchesProductFilter(incident entities.Incident, criteria valueobjects.FilterCriteria) bool {
	if !criteria.HasProductFilter() {
		return true
	}

	for _, affectedProduct := range incident.AffectedProducts {
		for _, filterProduct := range criteria.Products {
			if strings.Contains(affectedProduct.Title, filterProduct) {
				return true
			}
		}
	}

	return false
}
