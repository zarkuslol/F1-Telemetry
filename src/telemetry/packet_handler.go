// --- telemetry/packet_handler.go ---
package telemetry

import (
	"bytes"
	"encoding/binary"
	"joselucas/f1-telemetry/src/utils"
)

type KafkaProducerInterface interface {
	PublishTelemetry(data interface{}) error
}

var byteOrder = binary.LittleEndian

func ParseAndPublishPacket(data []byte, producer KafkaProducerInterface) {
	reader := bytes.NewReader(data)
	var header PacketHeader
	if err := binary.Read(reader, byteOrder, &header); err != nil {
		utils.Logger.Printf("Erro ao decodificar cabeçalho: %v", err)
		return
	}

	if header.PacketID == 6 {
		telemetryPacket, err := DecodeTelemetryPacket(reader)
		if err != nil {
			utils.Logger.Printf("Erro ao decodificar pacote de telemetria: %v", err)
			return
		}
		payload := BuildTelemetryPayload(telemetryPacket, header.PlayerCarIndex)
		if err := producer.PublishTelemetry(payload); err != nil {
			utils.Logger.Printf("Erro ao publicar no Kafka: %v", err)
		} else {
			utils.Logger.Printf("Payload publicado no Kafka para o carro %d.", header.PlayerCarIndex)
		}
	}
}
