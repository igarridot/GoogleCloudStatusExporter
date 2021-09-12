package main

import (
	"log"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type gcpStatusCollector struct {
	gcpStatus *prometheus.Desc
}

func newGcpStatus() *gcpStatusCollector {
	_, save_last_update := os.LookupEnv("SAVE_LAST_UPDATE")
	if save_last_update {
		return &gcpStatusCollector{
			gcpStatus: prometheus.NewDesc(prometheus.BuildFQName("gcp", "", "incidents"),
				"GCP Incident last update status",
				[]string{"id", "status", "product", "description", "uri", "last_update"}, nil,
			),
		}
	} else {
		return &gcpStatusCollector{
			gcpStatus: prometheus.NewDesc(prometheus.BuildFQName("gcp", "", "incidents"),
				"GCP Incident last update status",
				[]string{"id", "status", "product", "description", "uri"}, nil,
			),
		}
	}
}

func severityHandler(gcpIncident incident) float64 {
	var severity float64
	if gcpIncident.MostRecentUpdate.Updatestatus == "AVAILABLE" || gcpIncident.EndsAt != "" {
		severity = 0.0
	} else if gcpIncident.Severity == "high" {
		severity = 2.0
	} else if gcpIncident.Severity == "medium" {
		severity = 1.0
	}
	return severity
}

func (collector *gcpStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.gcpStatus
}

func (collector *gcpStatusCollector) ZoneHandler(ch chan<- prometheus.Metric, gcpIncident incident) {
	incidentSeverity := severityHandler(gcpIncident)
	incidentZones, _ := os.LookupEnv("INCIDENTS_ZONES")
	if incidentZones != "" {
		incidentZoneSlice := strings.Split(incidentZones, ",")
		specialZones := []string{"Global", "global", "various", "Various"}
		for _, zone := range append(incidentZoneSlice, specialZones...) {
			found := strings.Contains(gcpIncident.ExternalDescription, zone)
			if found {
				collector.ProductHandler(ch, gcpIncident, incidentSeverity)
			}
		}
	} else {
		collector.ProductHandler(ch, gcpIncident, incidentSeverity)
	}
}

func (collector *gcpStatusCollector) ProductHandler(ch chan<- prometheus.Metric, gcpIncident incident, incidentSeverity float64) {
	incidentProducts, _ := os.LookupEnv("FILTERED_PRODUCTS")
	if incidentProducts != "" {
		incidentProductsSlice := strings.Split(incidentProducts, ",")
		for _, affectedProduct := range gcpIncident.AffectedProducts {
			for _, product := range incidentProductsSlice {
				found := strings.Contains(affectedProduct.Title, product)
				if found {
					collector.AddMetric(ch, gcpIncident, incidentSeverity, affectedProduct)
				}
			}
		}
	} else {
		for _, affectedProduct := range gcpIncident.AffectedProducts {
			collector.AddMetric(ch, gcpIncident, incidentSeverity, affectedProduct)
		}
	}
}

func (collector *gcpStatusCollector) AddMetric(ch chan<- prometheus.Metric, gcpIncident incident, incidentSeverity float64, affectedProduct product) {
	_, saveLastUpdate := os.LookupEnv("SAVE_LAST_UPDATE")
	if saveLastUpdate {
		ch <- prometheus.MustNewConstMetric(collector.gcpStatus, prometheus.GaugeValue, incidentSeverity, gcpIncident.IncidentId, gcpIncident.MostRecentUpdate.Status, affectedProduct.Title, gcpIncident.ExternalDescription, gcpIncident.URI, gcpIncident.MostRecentUpdate.Status)
	} else {
		ch <- prometheus.MustNewConstMetric(collector.gcpStatus, prometheus.GaugeValue, incidentSeverity, gcpIncident.IncidentId, gcpIncident.MostRecentUpdate.Status, affectedProduct.Title, gcpIncident.ExternalDescription, "https://status.cloud.google.com/"+gcpIncident.URI)
	}
}

func (collector *gcpStatusCollector) Collect(ch chan<- prometheus.Metric) {

	gcpCurrentStatus, err := obtainGcpStatus()
	if err != nil {
		log.Fatal("Cannot obtain GCP Status")
	}

	for _, gcpIncident := range gcpCurrentStatus {
		_, CollectResolvedIncidents := os.LookupEnv("COLLECT_RESOLVED_EVENTS")
		if !CollectResolvedIncidents {
			if gcpIncident.MostRecentUpdate.Updatestatus != "AVAILABLE" && gcpIncident.EndsAt == "" {
				collector.ZoneHandler(ch, gcpIncident)
			}
		} else {
			collector.ZoneHandler(ch, gcpIncident)
		}
	}
}

func Register() {

	gcpIncident := newGcpStatus()
	prometheus.MustRegister(gcpIncident)
}
