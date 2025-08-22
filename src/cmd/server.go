package cmd

import (
	"net/http"

	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/infrastructure/config"
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/infrastructure/container"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetricServer() error {
	cfg := config.LoadConfig()

	containerInstance := container.NewContainer(cfg)
	collector := containerInstance.GetCollector()

	prometheus.MustRegister(collector)

	return serverMetrics(cfg.ListenAddress, cfg.MetricsPath)
}

func serverMetrics(listenAddress, metricsPath string) error {
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
