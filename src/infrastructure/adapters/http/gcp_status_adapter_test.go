package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockHTTPClient struct {
	ResponseBody string
	StatusCode   int
	Error        error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	response := &http.Response{
		StatusCode: m.StatusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(m.ResponseBody)),
	}

	return response, nil
}

func TestGCPStatusAdapter_GetIncidents(t *testing.T) {
	jsonResponse := `[
		{
			"id": "test-id",
			"number": "123",
			"begin": "2023-01-01T10:00:00Z",
			"created": "2023-01-01T10:00:00Z",
			"external_desc": "Test incident",
			"severity": "high",
			"affected_products": [
				{"title": "Test Product", "id": "test-product"}
			],
			"updates": [
				{
					"created": "2023-01-01T10:00:00Z",
					"modified": "2023-01-01T10:00:00Z",
					"text": "Investigating",
					"when": "2023-01-01T10:00:00Z",
					"status": "INVESTIGATING"
				}
			],
			"service_key": "test-service",
			"service_name": "Test Service",
			"uri": "/incident/test-id"
		}
	]`

	mockClient := &MockHTTPClient{
		ResponseBody: jsonResponse,
		StatusCode:   200,
	}

	adapter := NewGCPStatusAdapter(mockClient)
	incidents, err := adapter.GetIncidents()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(incidents) != 1 {
		t.Fatalf("Expected 1 incident, got %d", len(incidents))
	}

	incident := incidents[0]
	if incident.ID.String() != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", incident.ID.String())
	}

	if incident.Number != "123" {
		t.Errorf("Expected Number '123', got '%s'", incident.Number)
	}

	if len(incident.AffectedProducts) != 1 {
		t.Fatalf("Expected 1 affected product, got %d", len(incident.AffectedProducts))
	}

	if incident.AffectedProducts[0].Title != "Test Product" {
		t.Errorf("Expected product title 'Test Product', got '%s'", incident.AffectedProducts[0].Title)
	}
}
