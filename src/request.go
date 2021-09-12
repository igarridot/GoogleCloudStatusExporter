package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type incident struct {
	IncidentId          string    `json:"id"`
	IncidentNumber      string    `json:"number"`
	BeginsAt            string    `json:"begin"`
	CreatedAt           string    `json:"created"`
	EndsAt              string    `json:"end,omitempty"`
	ModifiedAt          string    `json:"modified,omitempty"`
	ExternalDescription string    `json:"external_desc"`
	Updates             []update  `json:"updates,omitempty"`
	MostRecentUpdate    update    `json:"most_recent_update,omitempty"`
	StatusImpact        string    `json:"status_impact"`
	Severity            string    `json:"severity"`
	AffectedProducts    []product `json:"affected_products"`
	ServiceKey          string    `json:"service_key"`
	ServiceName         string    `json:"service_name"`
	URI                 string    `json:"uri"`
}

type update struct {
	CreatedAt    string `json:"created"`
	ModifiedAt   string `json:"modified"`
	Status       string `json:"text"`
	UpdatedDate  string `json:"when"`
	Updatestatus string `json:"status"`
}

type product struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

func obtainGcpStatus() ([]incident, error) {
	var status []incident

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://status.cloud.google.com/incidents.json", nil)
	if err != nil {
		return []incident{}, errors.New("request builder has failed")
	}
	req.Header.Add("Content-Type", "application/json")
	fmt.Println("Polling GCP Status webpage...")
	resp, err := client.Do(req)
	if err != nil {
		return []incident{}, errors.New("request to GCP Status webpage failed")
	}
	fmt.Println("Successfully polled GCP status webpage")
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []incident{}, errors.New("cannot read the response body")
	}
	err = json.Unmarshal(responseBody, &status)
	if err != nil {
		return []incident{}, errors.New("failed to unmarshal response json body")
	}
	return status, nil
}
