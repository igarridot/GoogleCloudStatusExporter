package ports

import (
	"errors"
	"testing"
	"time"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/valueobjects"
)

// Mock implementation for testing interface compliance
type mockGCPStatusPort struct {
	incidents []entities.Incident
	err       error
}

func (m *mockGCPStatusPort) GetIncidents() ([]entities.Incident, error) {
	return m.incidents, m.err
}

func TestGCPStatusPortInterface(t *testing.T) {
	// Test successful case
	expectedIncidents := createTestIncidentsForPort()
	mock := &mockGCPStatusPort{
		incidents: expectedIncidents,
		err:       nil,
	}

	var port GCPStatusPort = mock

	incidents, err := port.GetIncidents()
	if err != nil {
		t.Errorf("GetIncidents() returned unexpected error: %v", err)
	}

	if len(incidents) != len(expectedIncidents) {
		t.Errorf("Expected %d incidents, got %d", len(expectedIncidents), len(incidents))
	}

	// Verify first incident details
	if len(incidents) > 0 {
		incident := incidents[0]
		expected := expectedIncidents[0]

		if incident.ID != expected.ID {
			t.Errorf("Expected incident ID %v, got %v", expected.ID, incident.ID)
		}

		if incident.Number != expected.Number {
			t.Errorf("Expected incident Number %v, got %v", expected.Number, incident.Number)
		}

		if incident.Severity != expected.Severity {
			t.Errorf("Expected incident Severity %v, got %v", expected.Severity, incident.Severity)
		}
	}
}

func TestGCPStatusPortInterface_Error(t *testing.T) {
	// Test error case
	expectedError := errors.New("failed to fetch incidents")
	mock := &mockGCPStatusPort{
		incidents: nil,
		err:       expectedError,
	}

	var port GCPStatusPort = mock

	incidents, err := port.GetIncidents()
	if err == nil {
		t.Error("Expected error from GetIncidents() but got nil")
	}

	if err.Error() != "failed to fetch incidents" {
		t.Errorf("Expected error message 'failed to fetch incidents', got '%s'", err.Error())
	}

	if incidents != nil {
		t.Errorf("Expected nil incidents when error occurs, got %v", incidents)
	}
}

func TestGCPStatusPortInterface_EmptyResult(t *testing.T) {
	// Test empty result case
	mock := &mockGCPStatusPort{
		incidents: []entities.Incident{},
		err:       nil,
	}

	var port GCPStatusPort = mock

	incidents, err := port.GetIncidents()
	if err != nil {
		t.Errorf("GetIncidents() returned unexpected error: %v", err)
	}

	if incidents == nil {
		t.Error("Expected empty slice, got nil")
	}

	if len(incidents) != 0 {
		t.Errorf("Expected 0 incidents, got %d", len(incidents))
	}
}

func TestGCPStatusPortInterface_MultipleImplementations(t *testing.T) {
	// Test that multiple implementations can coexist
	mock1 := &mockGCPStatusPort{
		incidents: createTestIncidentsForPort()[:1],
		err:       nil,
	}

	mock2 := &mockGCPStatusPort{
		incidents: createTestIncidentsForPort()[1:],
		err:       nil,
	}

	var port1 GCPStatusPort = mock1
	var port2 GCPStatusPort = mock2

	incidents1, err1 := port1.GetIncidents()
	incidents2, err2 := port2.GetIncidents()

	if err1 != nil || err2 != nil {
		t.Errorf("Unexpected errors: %v, %v", err1, err2)
	}

	if len(incidents1) != 1 {
		t.Errorf("Expected 1 incident from port1, got %d", len(incidents1))
	}

	if len(incidents2) != 1 {
		t.Errorf("Expected 1 incident from port2, got %d", len(incidents2))
	}

	// Verify they return different incidents
	if len(incidents1) > 0 && len(incidents2) > 0 {
		if incidents1[0].ID == incidents2[0].ID {
			t.Error("Different ports should return different incidents")
		}
	}
}

// Helper function to create test incidents
func createTestIncidentsForPort() []entities.Incident {
	id1, _ := valueobjects.NewIncidentID("incident-1")
	id2, _ := valueobjects.NewIncidentID("incident-2")

	return []entities.Incident{
		{
			ID:     id1,
			Number: "1",
			BeginTime: time.Now(),
			CreatedAt: time.Now(),
			ExternalDescription: "Test incident 1",
			MostRecentUpdate: entities.Update{
				Status:       "INVESTIGATING",
				UpdateStatus: "INVESTIGATING",
			},
			StatusImpact: "SERVICE_DISRUPTION",
			Severity:     valueobjects.SeverityHigh,
			AffectedProducts: []entities.Product{
				{Title: "Compute Engine", ID: "compute"},
			},
			ServiceKey:  "compute",
			ServiceName: "Compute Engine",
			URI:         "incident/incident-1",
		},
		{
			ID:     id2,
			Number: "2",
			BeginTime: time.Now(),
			CreatedAt: time.Now(),
			ExternalDescription: "Test incident 2",
			MostRecentUpdate: entities.Update{
				Status:       "AVAILABLE",
				UpdateStatus: valueobjects.UpdateStatusAvailable,
			},
			StatusImpact: "SERVICE_DISRUPTION",
			Severity:     valueobjects.SeverityMedium,
			AffectedProducts: []entities.Product{
				{Title: "Cloud Storage", ID: "storage"},
			},
			ServiceKey:  "storage",
			ServiceName: "Cloud Storage",
			URI:         "incident/incident-2",
		},
	}
}