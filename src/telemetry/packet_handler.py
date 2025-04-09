from telemetry.telemetry_decoder import TelemetryDecoder as td
from utils.logger import get_logger

class PacketHandler:
    def __init__(self):
        self.logger = get_logger(self.__class__.__name__)

    def parse(self, data):
        header, body = td.decode_header(data)
        packet_id = self._get_packet_id(header)

        self.logger.debug(f"Pacote com ID {packet_id} identificado.")

        if packet_id == 6:
            parsed = td.parse_car_telemetry(body)

            brakes = parsed.pop("brakes_temperature", None)
            if brakes:
                parsed["brakes_temperature_FL"] = brakes[0]
                parsed["brakes_temperature_FR"] = brakes[1]
                parsed["brakes_temperature_RL"] = brakes[2]
                parsed["brakes_temperature_RR"] = brakes[3]

            tyres_surface = parsed.pop("tyres_surface_temperature", None)
            if tyres_surface:
                parsed["tyres_surface_temperature_FL"] = tyres_surface[0]
                parsed["tyres_surface_temperature_FR"] = tyres_surface[1]
                parsed["tyres_surface_temperature_RL"] = tyres_surface[2]
                parsed["tyres_surface_temperature_RR"] = tyres_surface[3]

            tyres_inner = parsed.pop("tyres_inner_temperature", None)
            if tyres_inner:
                parsed["tyres_inner_temperature_FL"] = tyres_inner[0]
                parsed["tyres_inner_temperature_FR"] = tyres_inner[1]
                parsed["tyres_inner_temperature_RL"] = tyres_inner[2]
                parsed["tyres_inner_temperature_RR"] = tyres_inner[3]

            tyres_pressure = parsed.pop("tyres_pressure", None)
            if tyres_pressure:
                parsed["tyres_pressure_FL"] = tyres_pressure[0]
                parsed["tyres_pressure_FR"] = tyres_pressure[1]
                parsed["tyres_pressure_RL"] = tyres_pressure[2]
                parsed["tyres_pressure_RR"] = tyres_pressure[3]

            surface_type = parsed.pop("surface_type", None)
            if surface_type:
                parsed["surface_type_FL"] = surface_type[0]
                parsed["surface_type_FR"] = surface_type[1]
                parsed["surface_type_RL"] = surface_type[2]
                parsed["surface_type_RR"] = surface_type[3]
        
            return {
                "measurement": "car_telemetry",
                "fields": parsed,
                "tags": {"source": "game"}
            }
        return None
    
    def _get_packet_id(self, header):
        return header[5]
