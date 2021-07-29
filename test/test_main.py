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
    def test_inuque_incident(self, mock_get, mock_get2):
        exporter = GCPStatusCollector('', '', '', '', False)
        iterator = exporter.collect()
        alerts = []
        for item in iterator:
            alerts.append(item)
        self.assertEqual(str(alerts[0]), "Metric(gcp_incidents, GCP Incident last update status, gauge, , [Sample(name='gcp_incidents', labels={'id': 'MfiGCW4E26MPGRnCJ8by', 'status': 'SERVICE_DISRUPTION', 'product': 'Google Cloud Storage', 'description': 'us-central1: GCS is returning stale version of object for bucket', 'uri': 'incidents/MfiGCW4E26MPGRnCJ8by'}, value=0, timestamp=None, exemplar=None)])")

if __name__ == '__main__':
    unittest.main()