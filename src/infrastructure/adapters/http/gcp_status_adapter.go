package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

type GCPStatusAdapter struct {
	httpClient HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func NewGCPStatusAdapter(httpClient HTTPClient) *GCPStatusAdapter {
	return &GCPStatusAdapter{
		httpClient: httpClient,
	}
}

func (a *GCPStatusAdapter) GetIncidents() ([]entities.Incident, error) {
	req, err := http.NewRequest("GET", "https://status.cloud.google.com/incidents.json", nil)
	if err != nil {
		return nil, errors.New("request builder has failed")
	}
	req.Header.Add("Content-Type", "application/json")

	fmt.Println("Polling GCP Status webpage...")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("request to GCP Status webpage failed")
	}
	fmt.Println("Successfully polled GCP status webpage")
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("cannot read the response body")
	}

	var rawIncidents []rawIncident
	err = json.Unmarshal(responseBody, &rawIncidents)
	if err != nil {
		return nil, errors.New("failed to unmarshal response json body")
	}

	return a.convertToEntities(rawIncidents), nil
}

func (a *GCPStatusAdapter) convertToEntities(rawIncidents []rawIncident) []entities.Incident {
	var incidents []entities.Incident

	for _, raw := range rawIncidents {
		incident := entities.Incident{
			Number:              raw.Number,
			ExternalDescription: raw.ExternalDescription,
			StatusImpact:        raw.StatusImpact,
			Severity:            valueobjects.NewSeverity(raw.Severity),
			ServiceKey:          raw.ServiceKey,
			ServiceName:         raw.ServiceName,
			URI:                 raw.URI,
		}

		if id, err := valueobjects.NewIncidentID(raw.ID); err == nil {
			incident.ID = id
		}

		if beginTime, err := time.Parse(time.RFC3339, raw.Begin); err == nil {
			incident.BeginTime = beginTime
		}

		if createdAt, err := time.Parse(time.RFC3339, raw.Created); err == nil {
			incident.CreatedAt = createdAt
		}

		if raw.End != "" {
			if endTime, err := time.Parse(time.RFC3339, raw.End); err == nil {
				incident.EndTime = &endTime
			}
		}

		if raw.Modified != "" {
			if modifiedAt, err := time.Parse(time.RFC3339, raw.Modified); err == nil {
				incident.ModifiedAt = &modifiedAt
			}
		}

		for _, rawUpdate := range raw.Updates {
			update := entities.Update{
				Status:       rawUpdate.Text,
				UpdateStatus: valueobjects.NewUpdateStatus(rawUpdate.Status),
			}

			if createdAt, err := time.Parse(time.RFC3339, rawUpdate.Created); err == nil {
				update.CreatedAt = createdAt
			}

			if modifiedAt, err := time.Parse(time.RFC3339, rawUpdate.Modified); err == nil {
				update.ModifiedAt = modifiedAt
			}

			if updatedDate, err := time.Parse(time.RFC3339, rawUpdate.When); err == nil {
				update.UpdatedDate = updatedDate
			}

			incident.Updates = append(incident.Updates, update)
		}

		if len(raw.Updates) > 0 {
			lastUpdate := raw.Updates[len(raw.Updates)-1]
			mostRecentUpdate := entities.Update{
				Status:       lastUpdate.Text,
				UpdateStatus: valueobjects.NewUpdateStatus(lastUpdate.Status),
			}

			if createdAt, err := time.Parse(time.RFC3339, lastUpdate.Created); err == nil {
				mostRecentUpdate.CreatedAt = createdAt
			}

			if modifiedAt, err := time.Parse(time.RFC3339, lastUpdate.Modified); err == nil {
				mostRecentUpdate.ModifiedAt = modifiedAt
			}

			if updatedDate, err := time.Parse(time.RFC3339, lastUpdate.When); err == nil {
				mostRecentUpdate.UpdatedDate = updatedDate
			}

			incident.MostRecentUpdate = mostRecentUpdate
		}

		for _, rawProduct := range raw.AffectedProducts {
			product := entities.Product{
				Title: rawProduct.Title,
				ID:    rawProduct.ID,
			}
			incident.AffectedProducts = append(incident.AffectedProducts, product)
		}

		incidents = append(incidents, incident)
	}

	return incidents
}

type rawIncident struct {
	ID                  string       `json:"id"`
	Number              string       `json:"number"`
	Begin               string       `json:"begin"`
	Created             string       `json:"created"`
	End                 string       `json:"end,omitempty"`
	Modified            string       `json:"modified,omitempty"`
	ExternalDescription string       `json:"external_desc"`
	Updates             []rawUpdate  `json:"updates,omitempty"`
	StatusImpact        string       `json:"status_impact"`
	Severity            string       `json:"severity"`
	AffectedProducts    []rawProduct `json:"affected_products"`
	ServiceKey          string       `json:"service_key"`
	ServiceName         string       `json:"service_name"`
	URI                 string       `json:"uri"`
}

type rawUpdate struct {
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Text     string `json:"text"`
	When     string `json:"when"`
	Status   string `json:"status"`
}

type rawProduct struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}
