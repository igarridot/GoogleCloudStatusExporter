import json
import requests
import time
from prometheus_client import start_http_server
from prometheus_client.core import GaugeMetricFamily, REGISTRY

class GCPStatusCollector(object):
  def collect(self):
    metric = GaugeMetricFamily(
        'gcp_incidents',
        'GCP Incident number',
        labels=["incident_number", "product"])

    url = 'https://status.cloud.google.com/incidents.json'
    resp = requests.get(url=url)
    data = resp.json()

    for incident in data:
        if incident['most_recent_update']['status'] != 'AVAILABLE':
          incident_number = incident['number']
          for product in incident['affected_products']:
            metric.add_metric([incident_number, product['title']], 1)

    yield metric

if __name__ == "__main__":
  REGISTRY.register(GCPStatusCollector())
  start_http_server(9118)
  while True: time.sleep(1)