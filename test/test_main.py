import requests
import unittest
import json
from unittest import mock
from src.main import GCPStatusCollector

with open('test/fixtures/small.json') as fixture:
    single_issue_fixture = json.load(fixture)
    fixture.close()


def mocked_request_handler(*args, **kwargs):
    class MockResponse:
        def __init__(self, json_data, status_code):
            self.json_data = json_data
            self.status_code = status_code

        def json(self):
            return self.json_data

    return MockResponse(single_issue_fixture, 200)


class GCPStatusCollectorTestCase(unittest.TestCase):

    @mock.patch('src.main.requests.get', side_effect=mocked_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_default_parameters(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', [], [], False)
        iterator = exporter.collect()
        alerts = []
        for item in iterator:
            alerts.append(item)
        self.assertEqual(str(
            alerts[0]), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=2, timestamp=None, exemplar=None)])")


    @mock.patch('src.main.requests.get', side_effect=mocked_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_custom_status_endpoint(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://my.fake.site/incidents.json', [], [], False)
        iterator = exporter.collect()
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://my.fake.site/incidents/MfiGCW4E26MPGRnCJ8by'}, value=2, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_product_filter(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', ["My Fake Product Name", "Google Cloud Storage"], [], False)
        iterator = exporter.collect()
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=2, timestamp=None, exemplar=None)])")


    @mock.patch('src.main.requests.get', side_effect=mocked_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_zone_filter(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', [], ["My Fake Zone", "us-central1"], False)
        iterator = exporter.collect()
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=2, timestamp=None, exemplar=None)])")


    @mock.patch('src.main.requests.get', side_effect=mocked_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_product_and_zone_filters(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', ["My Fake Product Name", "Google Cloud Storage"], ["My Fake Zone", "us-central1"], False)
        iterator = exporter.collect()
        alerts = []
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=2, timestamp=None, exemplar=None)])")


    @mock.patch('src.main.requests.get', side_effect=mocked_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_default_parameters(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', [], [], True)
        iterator = exporter.collect()
        self.maxDiff = None
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fEXXEicMtx5SaVZz2Gt7', 'status': 'AVAILABLE', 'product': 'Google Cloud Scheduler', 'description': 'Global: Cloud Scheduler Pub/Sub jobs fail with permission denied', 'uri': 'https://status.cloud.google.com/incidents/fEXXEicMtx5SaVZz2Gt7'}, value=0, timestamp=None, exemplar=None)])")


if __name__ == '__main__':
    unittest.main()