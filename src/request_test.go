package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

const jsonFixture = `[{
	"id": "MfiGCW4E26MPGRnCJ8by",
	"number": "13858523345881857527",
	"begin": "2021-07-27T17:35:27+00:00",
	"created": "2021-07-27T18:00:25+00:00",
	"modified": "2021-07-27T23:15:35+00:00",
	"external_desc": "us-central1: GCS is returning stale version of object for bucket",
	"updates": [
		{
			"created": "2021-07-27T23:15:28+00:00",
			"modified": "2021-07-27T23:15:34+00:00",
			"when": "2021-07-27T23:15:28+00:00",
			"text": "The issue with Google Cloud Storage is believed to be affecting a very small number of projects and our Engineering Team is working on it.\nIf you have questions or are impacted, please open a case with the Support Team and we will work with you until this issue is resolved.\nNo further updates will be provided here.\nWe thank you for your patience while we're working on resolving the issue.",
			"status": "SERVICE_DISRUPTION"
		}
	],
	"most_recent_update": {
		"created": "2021-07-27T23:15:28+00:00",
		"modified": "2021-07-27T23:15:34+00:00",
		"when": "2021-07-27T23:15:28+00:00",
		"text": "The issue with Google Cloud Storage is believed to be affecting a very small number of projects and our Engineering Team is working on it.\nIf you have questions or are impacted, please open a case with the Support Team and we will work with you until this issue is resolved.\nNo further updates will be provided here.\nWe thank you for your patience while we're working on resolving the issue.",
		"status": "SERVICE_DISRUPTION"
	},
	"status_impact": "SERVICE_DISRUPTION",
	"severity": "high",
	"service_key": "UwaYoXQ5bHYHG6EdiPB8",
	"service_name": "Google Cloud Storage",
	"affected_products": [
		{
			"title": "Google Cloud Storage",
			"id": "UwaYoXQ5bHYHG6EdiPB8"
		}
	],
	"uri": "incidents/MfiGCW4E26MPGRnCJ8by"
}]
`

type FakeDoMethod struct{}

func (f FakeDoMethod) Do(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(jsonFixture)),
	}
	return response, nil
}
func TestObtainGcpStatus(t *testing.T) {
	client := &RealHttpClient{}
	doMethod := &FakeDoMethod{}
	got, _ := client.obtainGcpStatus(doMethod)
	want := &[]incident{
		{
			IncidentId:          "MfiGCW4E26MPGRnCJ8by",
			IncidentNumber:      "13858523345881857527",
			BeginsAt:            "2021-07-27T17:35:27+00:00",
			CreatedAt:           "2021-07-27T18:00:25+00:00",
			ModifiedAt:          "2021-07-27T23:15:35+00:00",
			ExternalDescription: "us-central1: GCS is returning stale version of object for bucket",
			Updates: []update{
				{
					CreatedAt:    "2021-07-27T23:15:28+00:00",
					ModifiedAt:   "2021-07-27T23:15:34+00:00",
					UpdatedDate:  "2021-07-27T23:15:28+00:00",
					Status:       "The issue with Google Cloud Storage is believed to be affecting a very small number of projects and our Engineering Team is working on it.\nIf you have questions or are impacted, please open a case with the Support Team and we will work with you until this issue is resolved.\nNo further updates will be provided here.\nWe thank you for your patience while we're working on resolving the issue.",
					Updatestatus: "SERVICE_DISRUPTION",
				},
			},
			MostRecentUpdate: update{
				CreatedAt:    "2021-07-27T23:15:28+00:00",
				ModifiedAt:   "2021-07-27T23:15:34+00:00",
				UpdatedDate:  "2021-07-27T23:15:28+00:00",
				Status:       "The issue with Google Cloud Storage is believed to be affecting a very small number of projects and our Engineering Team is working on it.\nIf you have questions or are impacted, please open a case with the Support Team and we will work with you until this issue is resolved.\nNo further updates will be provided here.\nWe thank you for your patience while we're working on resolving the issue.",
				Updatestatus: "SERVICE_DISRUPTION",
			},
			StatusImpact: "SERVICE_DISRUPTION",
			Severity:     "high",
			ServiceKey:   "UwaYoXQ5bHYHG6EdiPB8",
			ServiceName:  "Google Cloud Storage",
			AffectedProducts: []product{
				{
					Title: "Google Cloud Storage",
					Id:    "UwaYoXQ5bHYHG6EdiPB8",
				},
			},
			URI: "incidents/MfiGCW4E26MPGRnCJ8by",
		},
	}
	checkExpectedHttpResponseBody(t, got, want)

}

func checkExpectedHttpResponseBody(t *testing.T, got, want *[]incident) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
