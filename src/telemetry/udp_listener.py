import socket
from utils.logger import get_logger

class UDPListener:
    def __init__(self, ip, port, packet_handler, processor):
        self.ip = ip
        self.port = port
        self.packet_handler = packet_handler
        self.processor = processor
        self.logger = get_logger(self.__class__.__name__)

    def startup_listener(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
        sock.bind((self.ip, self.port))

        self.logger.info(f"Escutando pacotes na porta {self.port}...")

        while True:
            data, addr = sock.recvfrom(2048)
            self.logger.debug(f"Pacote recebido de {addr}")
            parsed_data = self.packet_handler.parse(data)
            if parsed_data:
                self.processor.process(parsed_data)
