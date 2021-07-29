# GoogleCloudStatusExporter
Simple exporter to monitor Google Cloud Platform issues in Prometheus format.

## Disclaimer
First of all, a couple of short disclaimers.

Due to the nature of [Google Cloud Platform Status](https://status.cloud.google.com) [JSON Schema](https://status.cloud.google.com/incidents.schema.json) this exporter can lead to [high cardinality](https://www.robustperception.io/cardinality-is-key) issues.

Is not guaranteed that "Zones" filter flag works as expected. Google used to mention the affected [zone(s)](https://cloud.google.com/docs/geography-and-regions#internal_services) in the alert brief description but is not defined in the Schema, so is not mandatory for them. In some cases, they also mentions geographical regions like ```northamerica``` or ```europe```. Special region ```Global``` will be automatically added to the filter list, they can also mention [multiregions](https://cloud.google.com/firestore/docs/locations#location-mr) under some services so, if you are decided to use the Zone filter feature we strongly recommend you to include this kind of words in the filter list.

## Key features
- Allows you to filter issues based on GCP [product names](https://status.cloud.google.com)(left column)
- Allows you to filter issues based on alert [geographical zones](https://cloud.google.com/docs/geography-and-regions#internal_services).
- Allows you to filter by only firing alerts or by all alerts including the resolved ones.
- The timeseries has a value based on the alert severity.

 
## Metrics
This Prometheus only returns one metric called ```gcp_incidents``` and the value of the metric is stablised based on the  alert severity:
  * resolved: 0
  * medium: 1
  * high: 2

Example metrics:

```
cp_incidents{description="Queries fail with RESOURCE_EXCEEDED",id="EdoHcVkqXbPQmz3qYtqb",product="Google BigQuery",status="SERVICE_DISRUPTION",uri="https://status.cloud.google.com/incidents/EdoHcVkqXbPQmz3qYtqb"} 1.0

gcp_incidents{description="We are experiencing an issue with Cloud AI in us-central1 starting at 12:13 US/Pacific.",id="UK3LcXtsL7sW9g8TZkJM",product="Cloud Machine Learning",status="SERVICE_DISRUPTION",uri="https://status.cloud.google.com/incidents/UK3LcXtsL7sW9g8TZkJM"} 1.0
```

 
## Labels
Labels are key. Each label will store the basic incident information:

| Label Name    | Value                       |
| ------------- |:---------------------------:|
| description   | brief alert description     |
| id            | unique issue identifier     |
| product       | affected product name       |
| status        | status of the incident      |
| uri           | Link to GCP incident page   |


## Configuration 

All the parameters can be introduced via environment variable or command argument. Command arguments have higher priority than environment variables:

| Short parameter    | Long Parameter           |  Environment variable name | Default Value | Format |
| -------------------|:------------------------:|:-----:|:-----:|:-----:|
| -e   | --gcp_status_endpoint | GCP_STATUS_ENDPOINT | https://status.cloud.google.com/incidents.json | URI as string |
| -p   | --listen_port | LISTEN_PORT | 9118 | integer |
| -d   | --debug_mode  | DEBUG |  | Boolean. Using the parameter without value will set is as True automatically |
| -z   | --zones | ZONES | | List of strings separated by space |
| -P   | --products | PRODUCTS | | List of strings separated by space |
| -a   | --manage_all_events | INCLUDE_RESOLVED | False | Boolean. Using the parameter without value will set is as True automatically |

## Docker deployment

You can build the image running the following target:

```
make build
```

Otherwise, the image is available in [Docker Hub](https://hub.docker.com/repository/docker/norbega/gcp-status-exporter)

## Future work
- Extend documentation to detail how to use the application without Docker
- Add more tests
- Create Grafana dashboard
- Create suggested alerts

## References
- Magnificent [Robustperception Blog](https://www.robustperception.io)
