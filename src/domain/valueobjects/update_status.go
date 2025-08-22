package valueobjects

type UpdateStatus string

const (
	UpdateStatusAvailable UpdateStatus = "AVAILABLE"
)

func NewUpdateStatus(value string) UpdateStatus {
	return UpdateStatus(value)
}
