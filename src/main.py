import json, requests, time, os, getopt, sys
from prometheus_client import start_http_server, Summary, REGISTRY, PROCESS_COLLECTOR, PLATFORM_COLLECTOR
from prometheus_client.core import GaugeMetricFamily

debug_mode = os.environ.get('DEBUG', False)
gcp_status_endpoint = os.environ.get('GCP_STATUS_ENDPOINT', 'https://status.cloud.google.com/incidents.json')
listen_port = os.environ.get('LISTEN_PORT',9118)
manage_all_events = False

argument_list = sys.argv[1:]
short_options = 'e:p:dP:z:a'
long_options = ['endpoint=', 'port=', 'debug', 'products=', 'zones=', 'all']

try:
    arguments, values = getopt.getopt(argument_list, short_options, long_options)
except getopt.error as err:
    print (str(err))
    sys.exit(2)

for argument, value in arguments:
  if argument == '-e' or argument == '--endpoint':
    gcp_status_endpoint = value
  elif argument == '-p' or argument == '--port':
    listen_port = value
  elif argument == '-d' or argument == '--debug':
    debug_mode = True
  elif argument == '-P' or argument == '--products':
    products = value.split(',')
  elif argument == '-z' or argument == '--zones':
    zones = value.split(',')
    zones.extend(['Global', 'global'])
  elif argument == '-a' or argument == '--all':
    manage_all_events = True

def severity_handler(incident):
  if incident['most_recent_update']['status'] != 'AVAILABLE':
    return 0
  elif incident['severity'] == 'high':
    return 2
  elif incident['severity'] == 'medium':
    return 1

def add_metric(incident, metric, incident_severity, product):
      metric.add_metric([incident['id'], incident['most_recent_update']['status'], product['title'], incident['external_desc'], gcp_status_endpoint.removesuffix('incidents.json')+incident['uri']], incident_severity)

def metric_handler(incident, metric, incident_severity):
  for product in incident['affected_products']:
    if 'products' in locals():
      if product['title'] in products:
        add_metric(incident, metric, incident_severity, product)
    else:
      add_metric(incident, metric, incident_severity, product)

def incident_handler(incident, metric):
  incident_severity = severity_handler(incident)
  if 'zones' in globals():
    for zone in zones:
      if zone in incident['external_desc']:
        metric_handler(incident, metric, incident_severity)
  else:
    metric_handler(incident, metric, incident_severity)

class GCPStatusCollector(object):
  def collect(self):
    metric = GaugeMetricFamily(
        'gcp_incidents',
        'GCP Incident last update status',
        labels=['id', 'status', 'product', 'description', 'uri'])

    resp = requests.get(url=gcp_status_endpoint)
    data = resp.json()

    for incident in data:
      if not manage_all_events:
        if incident['most_recent_update']['status'] != 'AVAILABLE':
          incident_handler(incident, metric)
      else:
        incident_handler(incident, metric)

    yield metric

if __name__ == '__main__':
  if not debug_mode:
    for coll in list(REGISTRY._collector_to_names.keys()):
      REGISTRY.unregister(coll)
  REGISTRY.register(GCPStatusCollector())
  start_http_server(int(listen_port))
  while True: time.sleep(1)