# main.py
from telemetry.udp_listener import UDPListener
from telemetry.packet_handler import PacketHandler
from telemetry.telemetry_processor import TelemetryProcessor
from handlers.data.influxdb_handler import InfluxDBHandler

influx_writer = InfluxDBHandler(
    host="influxdb",
    port=8086,
    org="f1-org",
    bucket="f1-bucket",
    token="meu-token-secreto"
)

processor = TelemetryProcessor(writer=influx_writer)
handler = PacketHandler()
listener = UDPListener(
    ip="0.0.0.0",
    port=20777,
    packet_handler=handler,
    processor=processor
)

listener.startup_listener()