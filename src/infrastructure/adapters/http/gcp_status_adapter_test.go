package http

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestDefaultHTTPClient_Do(t *testing.T) {
	// We test that the DefaultHTTPClient wrapper exists and has the correct method signature
	// The actual HTTP functionality is tested through integration tests with mocked responses
	client := &DefaultHTTPClient{}

	// Verify the client implements the HTTPClient interface
	var _ HTTPClient = client

	// Test that the Do method exists and can be called (interface compliance)
	// We don't make actual HTTP calls in unit tests - that's handled by mocked responses
	// client is never nil since it's created with &DefaultHTTPClient{}
	_ = client
}

func TestNewGCPStatusAdapter(t *testing.T) {
	mockClient := &mockHTTPClient{}
	adapter := NewGCPStatusAdapter(mockClient)

	if adapter == nil {
		t.Error("Expected adapter to be created, got nil")
		return
	}

	if adapter.httpClient != mockClient {
		t.Error("Expected httpClient to be set correctly")
	}
}

func TestGCPStatusAdapter_GetIncidents_Success(t *testing.T) {
	// Use small fixture data for testing
	fixtureData := `[
		{
			"id": "test-incident-1",
			"number": "123",
			"begin": "2023-01-01T10:00:00Z",
			"created": "2023-01-01T10:00:00Z",
			"end": "2023-01-01T11:00:00Z",
			"modified": "2023-01-01T10:30:00Z",
			"external_desc": "Test incident description",
			"updates": [
				{
					"created": "2023-01-01T10:00:00Z",
					"modified": "2023-01-01T10:00:00Z",
					"text": "Investigating issue",
					"when": "2023-01-01T10:00:00Z",
					"status": "INVESTIGATING"
				}
			],
			"status_impact": "SERVICE_DISRUPTION",
			"severity": "high",
			"affected_products": [
				{
					"title": "Compute Engine",
					"id": "compute"
				}
			],
			"service_key": "compute",
			"service_name": "Compute Engine",
			"uri": "incidents/test-incident-1"
		}
	]`

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(fixtureData)),
		},
	}

	adapter := NewGCPStatusAdapter(mockClient)
	incidents, err := adapter.GetIncidents()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(incidents) != 1 {
		t.Errorf("Expected 1 incident, got %d", len(incidents))
	}

	incident := incidents[0]
	if incident.ID.String() != "test-incident-1" {
		t.Errorf("Expected incident ID 'test-incident-1', got %s", incident.ID.String())
	}

	if incident.Severity != valueobjects.SeverityHigh {
		t.Errorf("Expected high severity, got %v", incident.Severity)
	}

	if len(incident.AffectedProducts) != 1 {
		t.Errorf("Expected 1 affected product, got %d", len(incident.AffectedProducts))
	}

	if len(incident.Updates) != 1 {
		t.Errorf("Expected 1 update, got %d", len(incident.Updates))
	}

	if incident.MostRecentUpdate.Status != "Investigating issue" {
		t.Errorf("Expected 'Investigating issue', got %s", incident.MostRecentUpdate.Status)
	}
}

func TestGCPStatusAdapter_GetIncidents_HTTPRequestFails(t *testing.T) {
	mockClient := &mockHTTPClient{
		err: errors.New("network error"),
	}

	adapter := NewGCPStatusAdapter(mockClient)
	incidents, err := adapter.GetIncidents()

	if err == nil {
		t.Error("Expected error when HTTP request fails")
	}

	if incidents != nil {
		t.Error("Expected nil incidents when HTTP request fails")
	}

	expectedError := "request to GCP Status webpage failed"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGCPStatusAdapter_GetIncidents_ResponseReadFails(t *testing.T) {
	// Create a response with a body that will fail to read
	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       &failingReader{},
		},
	}

	adapter := NewGCPStatusAdapter(mockClient)
	incidents, err := adapter.GetIncidents()

	if err == nil {
		t.Error("Expected error when response body read fails")
	}

	expectedError := "cannot read the response body"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	if incidents != nil {
		t.Error("Expected nil incidents when response read fails")
	}
}

func TestGCPStatusAdapter_GetIncidents_JSONUnmarshalFails(t *testing.T) {
	invalidJSON := `{"invalid": json}`

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(invalidJSON)),
		},
	}

	adapter := NewGCPStatusAdapter(mockClient)
	incidents, err := adapter.GetIncidents()

	if err == nil {
		t.Error("Expected error when JSON unmarshaling fails")
	}

	expectedError := "failed to unmarshal response json body"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	if incidents != nil {
		t.Error("Expected nil incidents when JSON unmarshal fails")
	}
}

func TestGCPStatusAdapter_convertToEntities_EmptyEndAndModified(t *testing.T) {
	// Test incident without end time and modified time
	rawIncidents := []rawIncident{
		{
			ID:                  "test-incident-2",
			Number:              "456",
			Begin:               "2023-01-01T10:00:00Z",
			Created:             "2023-01-01T10:00:00Z",
			ExternalDescription: "Test incident without end time",
			Severity:            "medium",
			Updates:             []rawUpdate{},
			AffectedProducts:    []rawProduct{},
			ServiceKey:          "storage",
			ServiceName:         "Cloud Storage",
			URI:                 "incidents/test-incident-2",
		},
	}

	adapter := &GCPStatusAdapter{}
	incidents := adapter.convertToEntities(rawIncidents)

	if len(incidents) != 1 {
		t.Errorf("Expected 1 incident, got %d", len(incidents))
	}

	incident := incidents[0]
	if incident.EndTime != nil {
		t.Error("Expected nil end time when not provided")
	}

	if incident.ModifiedAt != nil {
		t.Error("Expected nil modified time when not provided")
	}

	if incident.Severity != valueobjects.SeverityMedium {
		t.Errorf("Expected medium severity, got %v", incident.Severity)
	}
}

func TestGCPStatusAdapter_convertRawUpdateToEntity(t *testing.T) {
	adapter := &GCPStatusAdapter{}

	rawUpdate := rawUpdate{
		Created:  "2023-01-01T10:00:00Z",
		Modified: "2023-01-01T10:30:00Z",
		Text:     "Test update text",
		When:     "2023-01-01T10:00:00Z",
		Status:   "INVESTIGATING",
	}

	update := adapter.convertRawUpdateToEntity(rawUpdate)

	if update.Status != "Test update text" {
		t.Errorf("Expected 'Test update text', got %s", update.Status)
	}

	if update.UpdateStatus != valueobjects.NewUpdateStatus("INVESTIGATING") {
		t.Errorf("Expected INVESTIGATING status, got %v", update.UpdateStatus)
	}

	if update.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be parsed")
	}

	if update.ModifiedAt.IsZero() {
		t.Error("Expected ModifiedAt to be parsed")
	}

	if update.UpdatedDate.IsZero() {
		t.Error("Expected UpdatedDate to be parsed")
	}
}

func TestGCPStatusAdapter_convertRawUpdateToEntity_InvalidTimes(t *testing.T) {
	adapter := &GCPStatusAdapter{}

	rawUpdate := rawUpdate{
		Created:  "invalid-time",
		Modified: "invalid-time",
		Text:     "Test update text",
		When:     "invalid-time",
		Status:   "INVESTIGATING",
	}

	update := adapter.convertRawUpdateToEntity(rawUpdate)

	if !update.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be zero when time parsing fails")
	}

	if !update.ModifiedAt.IsZero() {
		t.Error("Expected ModifiedAt to be zero when time parsing fails")
	}

	if !update.UpdatedDate.IsZero() {
		t.Error("Expected UpdatedDate to be zero when time parsing fails")
	}
}

// Helper type for simulating read failures
type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read failure")
}

func (f *failingReader) Close() error {
	return nil
}
