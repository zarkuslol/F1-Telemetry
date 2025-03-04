import struct

class TelemetryParser:
    @staticmethod
    def decode_header(data):
        """Decodifica os cabeçalhos do pacote."""
        # Desempacotando os cabeçalhos do pacote
        header_format = '<HBBBBBIffIIBB'
        header_size = struct.calcsize(header_format)
        header_data = struct.unpack(header_format, data[:header_size])

        return header_data, data[header_size:]

    @staticmethod
    def parse_car_telemetry(data):
        """Parseia os dados de telemetria apenas do carro com ID 0."""

        car_telemetry_format = '<H3fBbH2B H4H4B4B H4f4B'
        car_telemetry_size = struct.calcsize(car_telemetry_format)
        car_telemetry_data = struct.unpack(car_telemetry_format, data[:car_telemetry_size])

        # Montando um dicionário com os dados de telemetria
        car_telemetry = {
            'speed': car_telemetry_data[0],  # m_speed
            'throttle': car_telemetry_data[1],  # m_throttle
            'steer': car_telemetry_data[2],  # m_steer
            'brake': car_telemetry_data[3],  # m_brake
            'clutch': car_telemetry_data[4],  # m_clutch
            'gear': car_telemetry_data[5],  # m_gear
            'engine_rpm': car_telemetry_data[6],  # m_engineRPM
            'drs': car_telemetry_data[7],  # m_drs
            'rev_lights_percent': car_telemetry_data[8],  # m_revLightsPercent
            'rev_lights_bit_value': car_telemetry_data[9],  # m_revLightsBitValue
            'brakes_temperature': car_telemetry_data[10:14],  # m_brakesTemperature[4]
            'tyres_surface_temperature': car_telemetry_data[14:18],  # m_tyresSurfaceTemperature[4]
            'tyres_inner_temperature': car_telemetry_data[18:22],  # m_tyresInnerTemperature[4]
            'engine_temperature': car_telemetry_data[22],  # m_engineTemperature
            'tyres_pressure': car_telemetry_data[23:27],  # m_tyresPressure[4]
            'surface_type': car_telemetry_data[27:31],  # m_surfaceType[4]
        }

        return car_telemetry
