import json
import requests
import time
import os
import argparse
import sys
from prometheus_client import start_http_server, Summary, REGISTRY, PROCESS_COLLECTOR, PLATFORM_COLLECTOR
from prometheus_client.core import GaugeMetricFamily


def parse_args():
    parser = argparse.ArgumentParser(
        description='endpoint port debug products zones all'
    )
    parser.add_argument(
        '-e', '--gcp_status_endpoint',
        required=False,
        help='Set URL to fetch GCP Incident status',
        default=os.environ.get('GCP_STATUS_ENDPOINT',
                               'https://status.cloud.google.com/incidents.json')
    )
    parser.add_argument(
        '-p', '--listen_port',
        type=int,
        required=False,
        help='Exporter listen port',
        default=int(os.environ.get('LISTEN_PORT', 9118))
    )
    parser.add_argument(
        '-d', '--debug_mode',
        required=False,
        action='store_true',
        help='Enable debug mode',
        default=bool(os.environ.get('DEBUG', False))
    )
    parser.add_argument(
        '-P', '--products',
        required=False,
        nargs="+",
        help='Products to monitor',
        default=os.environ.get('PRODUCTS', [])
    )
    parser.add_argument(
        '-z', '--zones',
        required=False,
        nargs="+",
        help='Geo zones to monitor',
        default=os.environ.get('ZONES', [])
    )
    parser.add_argument(
        '-a', '--manage_all_events',
        required=False,
        action='store_true',
        help='Monitor also resolved issues',
        default=os.environ.get('MANAGE_ALL_EVENTS', False)
    )
    parser.add_argument(
        '-u', '--last_update',
        required=False,
        action='store_true',
        help='Add extra label with last update information',
        default=os.environ.get('LAST_UPDATE', False)
    )
    return parser.parse_args()


class GCPStatusCollector(object):

    def __init__(self, gcp_status_endpoint, products, zones, manage_all_events, last_update):
        self.gcp_status_endpoint = gcp_status_endpoint
        self.products = products
        self.zones = zones
        if self.zones:
            self.zones.extend(['Global', 'global'])
        self.manage_all_events = manage_all_events
        self.last_update = last_update

    def collect(self):
        
        metric = GaugeMetricFamily(
            'gcp_incidents',
            'GCP Incident last update status',
            labels=['id', 'status', 'product', 'description', 'uri', 'last_update'])

        data = self.request_handler()

        for incident in data:
            if not self.manage_all_events:
                if incident['most_recent_update']['status'] != 'AVAILABLE' and not 'end' in incident:
                    self.incident_handler(incident, metric)
            else:
                self.incident_handler(incident, metric)

        yield metric

    def request_handler(self):
        print("Polling {}".format(self.gcp_status_endpoint))
        resp = requests.get(url=self.gcp_status_endpoint)
        print("GCP Incident status response code: {}".format(resp.status_code))
        return resp.json()

    def severity_handler(self, incident):
        if incident['most_recent_update']['status'] == 'AVAILABLE' or 'end' in incident:
            return 0
        elif incident['severity'] == 'high':
            return 2
        elif incident['severity'] == 'medium' or incident['severity'] == 'low':
            return 1

    def incident_handler(self, incident, metric):
        incident_severity = self.severity_handler(incident)
        if self.zones:
            for zone in self.zones:
                if zone in incident['external_desc']:
                    self.metric_handler(incident, metric, incident_severity)
        else:
            self.metric_handler(incident, metric, incident_severity)

    def metric_handler(self, incident, metric, incident_severity):
        for product in incident['affected_products']:
            if self.products:
                if product['title'] in self.products:
                    self.add_metric(incident, metric,
                                    incident_severity, product)
            else:
                self.add_metric(incident, metric, incident_severity, product)

    def add_metric(self, incident, metric, incident_severity, product):
        if self.last_update:
            metric.add_metric([incident['id'], incident['most_recent_update']['status'], product['title'], incident['external_desc'],
                            self.gcp_status_endpoint.removesuffix('incidents.json')+incident['uri'], incident['most_recent_update']['text']], incident_severity)
        else:
            metric.add_metric([incident['id'], incident['most_recent_update']['status'], product['title'], incident['external_desc'],
                            self.gcp_status_endpoint.removesuffix('incidents.json')+incident['uri']], incident_severity)


def main():
    try:
        args = parse_args()
        if type(args.zones) == str:
            args.zones = args.zones.split(",")
        if type(args.products) == str:
            args.products = args.products.split(",")

        if args.debug_mode:
            print("ENDPOINNT: ", args.gcp_status_endpoint)
            print("PORT: ", type(args.debug_mode), args.debug_mode)
            print("PRODUCTS: ", type(args.products),
                  len(args.products), args.products)
            print("ZONES: ", type(args.zones), len(args.zones), args.zones)
            print("EVENTS: ", type(args.manage_all_events),
                  args.manage_all_events)
            print("UPDATE: ", type(args.last_update), args.last_update)

        if not args.debug_mode:
            for coll in list(REGISTRY._collector_to_names.keys()):
                REGISTRY.unregister(coll)
        REGISTRY.register(GCPStatusCollector(args.gcp_status_endpoint,
                          args.products, args.zones, args.manage_all_events, args.last_update))
        start_http_server(args.listen_port)
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print(" Interrupted")
        exit(0)


if __name__ == "__main__":
    main()
