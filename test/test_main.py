import requests
import unittest
import json
from unittest import mock
from src.main import GCPStatusCollector

with open('test/fixtures/small.json') as fixture:
    single_issue_fixture = json.load(fixture)
    fixture.close()

with open('test/fixtures/missing_available_status.json') as fixture:
    missing_available_status_fixture = json.load(fixture)
    fixture.close()

with open('test/fixtures/low_severity_incident.json') as fixture:
    low_severity_incident_fixture = json.load(fixture)
    fixture.close()


def mocked_single_request_handler(*args, **kwargs):
    class MockResponse:
        def __init__(self, json_data, status_code):
            self.json_data = json_data
            self.status_code = status_code

        def json(self):
            return self.json_data

    return MockResponse(single_issue_fixture, 200)


def mocked_missing_available_request_handler(*args, **kwargs):
    class MockResponse:
        def __init__(self, json_data, status_code):
            self.json_data = json_data
            self.status_code = status_code

        def json(self):
            return self.json_data

    return MockResponse(missing_available_status_fixture, 200)


def mocked_low_severity_incident_fixture_handler(*args, **kwargs):
    class MockResponse:
        def __init__(self, json_data, status_code):
            self.json_data = json_data
            self.status_code = status_code

        def json(self):
            return self.json_data

    return MockResponse(low_severity_incident_fixture, 200)


class GCPStatusCollectorTestCase(unittest.TestCase):

    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_default_parameters(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://status.cloud.google.com/incidents.json', [], [], False, False)
        iterator = exporter.collect()
        alerts = []
        for item in iterator:
            alerts.append(item)
        self.assertEqual(str(
            alerts[0]), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=3, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_custom_status_endpoint(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://my.fake.site/incidents.json', [], [], False, False)
        iterator = exporter.collect()
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://my.fake.site/incidents/MfiGCW4E26MPGRnCJ8by'}, value=3, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_product_filter(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', [
                                      "My Fake Product Name", "Google Cloud Storage"], [], False, False)
        iterator = exporter.collect()
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=3, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_zone_filter(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://status.cloud.google.com/incidents.json', [], ["My Fake Zone", "us-central1"], False, False)
        iterator = exporter.collect()
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=3, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_product_and_zone_filters(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('https://status.cloud.google.com/incidents.json', [
                                      "My Fake Product Name", "Google Cloud Storage"], ["My Fake Zone", "us-central1"], False, False)
        iterator = exporter.collect()
        alerts = []
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=3, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_manage_all_events(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://status.cloud.google.com/incidents.json', [], [], True, False)
        iterator = exporter.collect()
        self.maxDiff = None
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=3, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fEXXEicMtx5SaVZz2Gt7', 'status': 'AVAILABLE', 'product': 'Google Cloud Scheduler', 'description': 'Global: Cloud Scheduler Pub/Sub jobs fail with permission denied', 'uri': 'https://status.cloud.google.com/incidents/fEXXEicMtx5SaVZz2Gt7'}, value=0, timestamp=None, exemplar=None)])")


    @mock.patch('src.main.requests.get', side_effect=mocked_single_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_incident_with_extra_last_update_label(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://status.cloud.google.com/incidents.json', [], [], False, True)
        iterator = exporter.collect()
        self.maxDiff = None
        for item in iterator:
            self.assertRegex(str(
                item), r'^.+,\ \'last_update\'\:.+$')


    @mock.patch('src.main.requests.get', side_effect=mocked_missing_available_request_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_unique_solved_incident_without_available_state(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://status.cloud.google.com/incidents.json', [], [], True, False)
        iterator = exporter.collect()
        self.maxDiff = None
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'https://status.cloud.google.com/incidents/MfiGCW4E26MPGRnCJ8by'}, value=0, timestamp=None, exemplar=None)])")

    @mock.patch('src.main.requests.get', side_effect=mocked_low_severity_incident_fixture_handler)
    @mock.patch('src.main.print', side_effect='Mocked')
    def test_error_state(self, mock_get, mock_get2):
        exporter = GCPStatusCollector(
            'https://status.cloud.google.com/incidents.json', [], [], True, False)
        iterator = exporter.collect()
        self.maxDiff = None
        for item in iterator:
            self.assertEqual(str(
                item), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'XVq5om2XEDSqLtJZUvcH', 'status': 'SERVICE_INFORMATION', 'product': 'Google Compute Engine', 'description': 'Cooling related failure in one of our buildings that hosts zone europe-west2-a for region europe-west2. Europe-west2-b and europe-west2-c are not impacted for VMs. We have fixed the previously occurring issues when creating new Persistent Disk devices. Zonal autoscaling for europe-west2-a is impacted for customers who suffered VM terminations.', 'uri': 'https://status.cloud.google.com/incidents/XVq5om2XEDSqLtJZUvcH'}, value=1, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google BigQuery', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Virtual Private Cloud (VPC)', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Cloud Spanner', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Compute Engine', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Kubernetes Engine', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Cloud Memorystore', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Bigtable', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Persistent Disk', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Dataflow', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Networking', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'API Gateway', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Composer', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud SQL', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Cloud Filestore', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google App Engine', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Managed Service for Microsoft Active Directory (AD)', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Functions', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Cloud Data Fusion', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Vertex AI Online Prediction', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Tasks', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Dataproc', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None), Sample(name='gcp_incidents', labels={'id': 'fmEL9i2fArADKawkZAa2', 'status': 'SERVICE_DISRUPTION', 'product': 'Cloud Machine Learning', 'description': 'Multiple Cloud products experiencing elevated error rates, latencies or service unavailability in europe-west2', 'uri': 'https://status.cloud.google.com/incidents/fmEL9i2fArADKawkZAa2'}, value=2, timestamp=None, exemplar=None)])")


if __name__ == '__main__':
    unittest.main()
