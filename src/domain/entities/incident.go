package entities

import (
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

type Incident struct {
	ID                  valueobjects.IncidentID
	Number              string
	BeginTime           time.Time
	CreatedAt           time.Time
	EndTime             *time.Time
	ModifiedAt          *time.Time
	ExternalDescription string
	Updates             []Update
	MostRecentUpdate    Update
	StatusImpact        string
	Severity            valueobjects.Severity
	AffectedProducts    []Product
	ServiceKey          string
	ServiceName         string
	URI                 string
}

type Update struct {
	CreatedAt    time.Time
	ModifiedAt   time.Time
	Status       string
	UpdatedDate  time.Time
	UpdateStatus valueobjects.UpdateStatus
}

type Product struct {
	Title string
	ID    string
}

func (i *Incident) IsResolved() bool {
	return i.MostRecentUpdate.UpdateStatus == valueobjects.UpdateStatusAvailable || i.EndTime != nil
}

func (i *Incident) GetSeverityValue() float64 {
	if i.IsResolved() {
		return 0.0
	}
	return i.Severity.ToFloat64()
}
