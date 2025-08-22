package valueobjects

import "errors"

type IncidentID string

func NewIncidentID(value string) (IncidentID, error) {
	if value == "" {
		return "", errors.New("incident ID cannot be empty")
	}
	return IncidentID(value), nil
}

func (id IncidentID) String() string {
	return string(id)
}
