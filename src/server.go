package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startMetricServer() {
	var (
		listenAddress         = flag.String("web.listen-address", ":9118", "Address to listen on for web interface.")
		metricPath            = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
		lastUpdate            = flag.Bool("exporter.save-last-update", false, "Use flag if you want to save last incident update description. Disabled by default.")
		collectResolvedEvents = flag.Bool("exporter.collect-resolved-incidents", false, "Use flag if you want to collect already resolved incidents. Disabled by default.")
		incidentsZones        = flag.String("exporter.incidents-zones", "", "Use this flag if you want to filter by incident zones.")
		filteredProducts      = flag.String("exporter.filtered-products", "", "Use this flag if you want to filter incidents by product names.")
	)
	flag.Parse()
	if *lastUpdate {
		os.Setenv("SAVE_LAST_UPDATE", "true")
	}
	if *collectResolvedEvents {
		os.Setenv("COLLECT_RESOLVED_EVENTS", "true")
	}
	if *incidentsZones != "" {
		os.Setenv("INCIDENTS_ZONES", *incidentsZones)
	}
	if *filteredProducts != "" {
		os.Setenv("FILTERED_PRODUCTS", *filteredProducts)
	}
	log.Fatal(serverMetrics(*listenAddress, *metricPath))
}

func serverMetrics(listenAddress, metricsPath string) error {
	flag.Parse()
	Register()
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
            <html>
            <head><title>GCP status Exporter Metrics</title></head>
            <body>
            <h1>GCP Status Prometheus exporter</h1>
            <p><a href='` + metricsPath + `'>Metrics</a></p>
            </body>
            </html>
        `))
	})
	return http.ListenAndServe(listenAddress, nil)
}
