from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import SYNCHRONOUS
from handlers.data.data_handler import DataHandler

class InfluxDBHandler(DataHandler):
    def __init__(self, host, port, org, bucket, token):
        url = f"http://{host}:{port}"
        self.client = InfluxDBClient(url=url, token=token, org=org)
        self.write_api = self.client.write_api(write_options=SYNCHRONOUS)
        self.bucket = bucket
        self.org = org

    def write(self, measurement, fields, tags=None, timestamp=None):
        point = Point(measurement)
        tags = tags or {}

        for k, v in tags.items():
            point = point.tag(k, v)

        for k, v in fields.items():
            point = point.field(k, v)

        if timestamp:
            point = point.time(timestamp, WritePrecision.NS)

        self.write_api.write(bucket=self.bucket, org=self.org, record=point)
