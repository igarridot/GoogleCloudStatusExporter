package testutils

import (
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

// CreateTestIncidents creates a standard set of test incidents for use across multiple test files.
// Returns two incidents: one high severity incident and one resolved low severity incident.
func CreateTestIncidents() []entities.Incident {
	id1, _ := valueobjects.NewIncidentID("incident-1")
	id2, _ := valueobjects.NewIncidentID("incident-2")

	endTime := time.Now()

	return []entities.Incident{
		{
			ID:                  id1,
			ExternalDescription: "High severity incident",
			Severity:            valueobjects.SeverityHigh,
			AffectedProducts: []entities.Product{
				{Title: "Apigee", ID: "apigee"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: "INVESTIGATING",
			},
		},
		{
			ID:                  id2,
			ExternalDescription: "Resolved incident",
			Severity:            valueobjects.SeverityLow,
			EndTime:             &endTime,
			AffectedProducts: []entities.Product{
				{Title: "Cloud Storage", ID: "storage"},
			},
			MostRecentUpdate: entities.Update{
				UpdateStatus: valueobjects.UpdateStatusAvailable,
			},
		},
	}
}

// CreateFilterTestIncidents creates incidents specifically for testing filter functionality.
// Returns incidents with different zones and products for filter testing.
func CreateFilterTestIncidents() []entities.Incident {
	id1, _ := valueobjects.NewIncidentID("incident-1")
	id2, _ := valueobjects.NewIncidentID("incident-2")

	return []entities.Incident{
		{
			ID:                  id1,
			ExternalDescription: "Issue in us-east1 zone",
			Severity:            valueobjects.SeverityHigh,
			AffectedProducts:    []entities.Product{{Title: "Compute Engine", ID: "compute"}},
			MostRecentUpdate:    entities.Update{UpdateStatus: "INVESTIGATING"},
		},
		{
			ID:                  id2,
			ExternalDescription: "Resolved issue",
			Severity:            valueobjects.SeverityLow,
			EndTime:             &time.Time{},
			AffectedProducts:    []entities.Product{{Title: "Cloud Storage", ID: "storage"}},
			MostRecentUpdate:    entities.Update{UpdateStatus: valueobjects.UpdateStatusAvailable},
		},
	}
}
