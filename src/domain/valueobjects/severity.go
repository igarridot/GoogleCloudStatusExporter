package valueobjects

type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

func (s Severity) ToFloat64() float64 {
	switch s {
	case SeverityLow:
		return 1.0
	case SeverityMedium:
		return 2.0
	case SeverityHigh:
		return 3.0
	default:
		return 0.0
	}
}

func NewSeverity(value string) Severity {
	switch value {
	case "low":
		return SeverityLow
	case "medium":
		return SeverityMedium
	case "high":
		return SeverityHigh
	default:
		return ""
	}
}
