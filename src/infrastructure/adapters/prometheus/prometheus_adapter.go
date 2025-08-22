package prometheus

import (
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/application/ports"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/domain/entities"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusAdapter struct {
	gcpStatus *prometheus.Desc
	config    ports.MetricsConfig
}

func NewPrometheusAdapter(config ports.MetricsConfig) *PrometheusAdapter {
	var labels []string
	if config.SaveLastUpdate {
		labels = []string{"id", "status", "product", "description", "uri", "last_update"}
	} else {
		labels = []string{"id", "status", "product", "description", "uri"}
	}

	return &PrometheusAdapter{
		gcpStatus: prometheus.NewDesc(
			prometheus.BuildFQName("gcp", "", "incidents"),
			"GCP Incident last update status",
			labels,
			nil,
		),
		config: config,
	}
}

func (a *PrometheusAdapter) Describe(ch chan<- *prometheus.Desc) {
	ch <- a.gcpStatus
}

func (a *PrometheusAdapter) CollectMetrics(incidents []entities.Incident, config ports.MetricsConfig) error {
	return nil
}

func (a *PrometheusAdapter) CollectMetricsWithChannel(incidents []entities.Incident, ch chan<- prometheus.Metric) {
	for _, incident := range incidents {
		for _, product := range incident.AffectedProducts {
			a.addMetric(ch, incident, product)
		}
	}
}

func (a *PrometheusAdapter) addMetric(ch chan<- prometheus.Metric, incident entities.Incident, product entities.Product) {
	incidentSeverity := incident.GetSeverityValue()
	uri := "https://status.cloud.google.com/" + incident.URI

	if a.config.SaveLastUpdate {
		ch <- prometheus.MustNewConstMetric(
			a.gcpStatus,
			prometheus.GaugeValue,
			incidentSeverity,
			incident.ID.String(),
			incident.MostRecentUpdate.Status,
			product.Title,
			incident.ExternalDescription,
			uri,
			incident.MostRecentUpdate.Status,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			a.gcpStatus,
			prometheus.GaugeValue,
			incidentSeverity,
			incident.ID.String(),
			incident.MostRecentUpdate.Status,
			product.Title,
			incident.ExternalDescription,
			uri,
		)
	}
}
