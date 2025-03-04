import socket
from telemetry.parser import TelemetryParser as tp
from dotenv import load_dotenv
import os

load_dotenv()

class UDPListener:
    def __init__(self):
        self._UDP_IP = os.getenv("APP_IP_ADDRESS")
        self._UDP_PORT = int(os.getenv("APP_PORT"))

    def startup_listener(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)  # Permitir broadcast
        sock.bind((self._UDP_IP, self._UDP_PORT))

        print(f"Escutando pacotes na porta {self._UDP_PORT}...")

        while True:
            data, addr = sock.recvfrom(2048)  # Tamanho do buffer
            header, body = tp.decode_header(data)
            print(f"ID do pacote vindo de get_packet_id: {self._get_packet_id(header)}")

            if self._get_packet_id(header) == 6:
                print(f"Dado vindo do Car Telemetry: {tp.parse_car_telemetry(body)}")
    
    def _get_packet_id(self, header):
        """Retorna o ID do pacote."""
        return header[5]

    def _print_packet(self, header, data):
        """Processa e exibe os dados brutos do pacote recebido."""
        packet_format, game_year, game_major_version, game_minor_version, packet_version, packet_id, session_uid, session_time, frame_identifier, overall_frame_identifier, player_car_index, secondary_player_car_index = header

        print(f"📦 Pacote recebido!")
        print(f"  Formato do pacote: {packet_format}")
        print(f"  Ano do jogo: {game_year}")
        print(f"  Versão do jogo: {game_major_version}.{game_minor_version}")
        print(f"  Versão do pacote: {packet_version}")
        print(f"  ID do pacote: {packet_id}")
        print(f"  Session UID: {session_uid}")
        print(f"  Tempo da sessão: {session_time}")
        print(f"  Identificador do quadro: {frame_identifier}")
        print(f"  Identificador geral do quadro: {overall_frame_identifier}")
        print(f"  Índice do carro do jogador: {player_car_index}")
        print(f"  Índice do segundo carro do jogador (se houver): {secondary_player_car_index}")

        # Logando os primeiros 100 bytes para inspeção
        print(f"🔍 Dados brutos: {data[:100].hex()}...")
