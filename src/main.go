package main

import (
	"github.com/igarridot/GoogleCloudStatusExporter/v2.0.0/cmd"
	"log"
)

func main() {
	log.Fatal(cmd.StartMetricServer())
}
