from utils.logger import get_logger
from handlers.data.data_handler import DataHandler

class TelemetryProcessor:
    def __init__(self, writer: DataHandler):
        self.writer = writer
        self.logger = get_logger(self.__class__.__name__)

    def process(self, parsed_packet: dict):
        self.logger.info(f"Processando dados: {parsed_packet['measurement']}")
        self.writer.write(
            measurement=parsed_packet["measurement"],
            fields=parsed_packet["fields"],
            tags=parsed_packet.get("tags", {})
        )
