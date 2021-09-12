# Google Cloud Status Exporter
Simple exporter to monitor Google Cloud Platform issues in Prometheus format.

## Disclaimer
First of all, a couple of short disclaimers.

Due to the nature of [Google Cloud Platform Status](https://status.cloud.google.com) [JSON Schema](https://status.cloud.google.com/incidents.schema.json) this exporter can lead to [high cardinality](https://www.robustperception.io/cardinality-is-key) issues.

Is not guaranteed that "Zones" filter flag works as expected. Google used to mention the affected [zone(s)](https://cloud.google.com/docs/geography-and-regions#internal_services) in the alert brief description but is not defined in the Schema, so is not mandatory for them. In some cases, they also mentions geographical regions like ```northamerica``` or ```europe```. Special region ```Global``` will be automatically added to the filter list, they can also mention [multiregions](https://cloud.google.com/firestore/docs/locations#location-mr) under some services so, if you are decided to use the Zone filter feature we strongly recommend you to include this kind of words in the filter list.

---------------------------------

## Key features
- Allows you to filter issues based on GCP [product names](https://status.cloud.google.com) (left column)
- Allows you to filter issues based on alert [geographical zones](https://cloud.google.com/docs/geography-and-regions#internal_services).
- Allows you to filter by only firing alerts or by all alerts including the resolved ones.
- You can store last incident status as a label content if you need it.
- The timeseries has a value based on the alert severity.

---------------------------------
 
## Metrics
This Prometheus only returns one metric called ```gcp_incidents``` and the value of the metric is stablised based on the  alert severity:
  * resolved: 0
  * medium: 1
  * high: 2

Example metrics:

```
gcp_incidents{description="Queries fail with RESOURCE_EXCEEDED",id="EdoHcVkqXbPQmz3qYtqb",product="Google BigQuery",status="SERVICE_DISRUPTION",uri="https://status.cloud.google.com/incidents/EdoHcVkqXbPQmz3qYtqb"} 1.0

gcp_incidents{description="We are experiencing an issue with Cloud AI in us-central1 starting at 12:13 US/Pacific.",id="UK3LcXtsL7sW9g8TZkJM",product="Cloud Machine Learning",status="SERVICE_DISRUPTION",uri="https://status.cloud.google.com/incidents/UK3LcXtsL7sW9g8TZkJM"} 1.0
```

Example metric including last_update label:

```
gcp_incidents{description="An issue with Cloud Healthcare API in asia-east2 has been resolved",id="w1sMLXwN9R3NK46UEZAx",last_update="Cloud Healthcare API has been affected in the asia-east2 region by the Google incident https://status.cloud.google.com/incident/zall/20009 since 2020-09-17 17:02 US/Pacific. The issue was resolved for all projects as of Thursday, 2020-09-17 18:38 US/Pacific.\nWe thank you for your patience while we worked on resolving the issue.",product="Healthcare and Life Sciences",status="SERVICE_DISRUPTION",uri="https://status.cloud.google.com/incidents/w1sMLXwN9R3NK46UEZAx"} 1.0
```

---------------------------------

## Labels
Each label will store the basic incident information:

| Label Name    | Value                       | Label Type |
| ------------- |:---------------------------:|:----------:|
| description   | brief alert description     | base       |
| id            | unique issue identifier     | base       |
| product       | affected product name       | base       |
| status        | status of the incident      | base       |
| uri           | Link to GCP incident page   | base       |
| last_update   | Last incident update text   | optional   |


---------------------------------

## Configuration

Work In Progress...

All the parameters can be introduced via environment variable or command argument. Command arguments have higher priority than environment variables:

### Entrypoint parameters
| Long Param Name        |  Default Value                 | Example                                                              |
|:----------------------:|:------------------------------:|:--------------------------------------------------------------------:|
| --web.listen-address   | 9118                           | ```--web.listen-address 9118```                                                   |
| --exporter.filtered-products             |      | ```--exporter.filtered-products "Healthcare and Life Sciences,Cloud Machine Learning"```                             |
| --exporter.incidents-zones                |      | ```--exporter.incidents-zones "asia-east2,Multi-Region"```                                                            |
| --exporter.collect-resolved-incidents    | False | ```--exporter.collect-resolved-incidents```
| --exporter.save-last-update          | False | ```--exporter.save-last-update``` |
---------------------------------

## Build image
You can build the image running the following target:

```
make build
```

Otherwise, the image is available in [Docker Hub](https://hub.docker.com/layers/166885628/norbega/gcp-status-exporter/v2.0.0-rc1/images/sha256-6278c4e92de2ba6328fb7c1c8f0f271548e2eed9687b86b6bbd4c3a159157948?context=repo)

---------------------------------

## Usage outside of containers

- Run the application:

```
./gcp-status-exporter-linux
```

- Example with parameters:

```
./gcp-status-exporter-osx --web.listen-address ":9119" --web.metrics-path "/metrics" --exporter.save-last-update --exporter.collect-resolved-incidents --exporter.incidents-zones "europe-west2,multizone" --exporter.filtered-products "Google Cloud Datastore,Google Cloud DNS"
```

---------------------------------

## Grafana Dashboard

Grafana directory contains a Dashboard JSON file that looks like this:

![Captura de pantalla 2021-07-31 a las 11 56 48](https://user-images.githubusercontent.com/33375539/127736641-3196be0e-87b5-40fa-92dd-11d9b569e894.png)


---------------------------------

## Referencies
- Prometheus Go client [guide](https://prometheus.io/docs/guides/go-application/)
- Prometheus Go client [library](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)
