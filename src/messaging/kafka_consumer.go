package messaging

import (
	"context"
	"encoding/json"
	"joselucas/f1-telemetry/src/handlers"
	"joselucas/f1-telemetry/src/telemetry"
	"joselucas/f1-telemetry/src/utils"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer consome dados de telemetria e os envia para o InfluxDB.
type KafkaConsumer struct {
	Reader        *kafka.Reader
	InfluxHandler *handlers.InfluxDBHandler
}

// NewKafkaConsumer cria e configura um novo consumidor Kafka.
func NewKafkaConsumer(brokers []string, topic, groupID string, influxHandler *handlers.InfluxDBHandler) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID, // GroupID garante que as mensagens sejam distribuídas entre os consumidores do mesmo grupo
		MinBytes: 10e3,    // 10KB
		MaxBytes: 10e6,    // 10MB
	})
	utils.Logger.Printf("Consumidor Kafka inicializado no tópico '%s'", topic)
	return &KafkaConsumer{
		Reader:        reader,
		InfluxHandler: influxHandler,
	}
}

// Consume e processa mensagens do Kafka em um loop infinito.
// Deve ser executado em uma goroutine.
func (kc *KafkaConsumer) ConsumeInfluxDB(ctx context.Context) {
	utils.Logger.Println("Iniciando consumo de mensagens do Kafka...")
	for {
		// O 'FetchMessage' espera por novas mensagens ou pelo cancelamento do contexto.
		msg, err := kc.Reader.FetchMessage(ctx)
		if err != nil {
			// Se o contexto foi cancelado, encerramos o loop.
			if ctx.Err() != nil {
				utils.Logger.Println("Contexto cancelado. Encerrando consumidor Kafka.")
				return
			}
			utils.Logger.Printf("Erro ao buscar mensagem do Kafka: %v", err)
			continue
		}

		var payload telemetry.TelemetryPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			utils.Logger.Printf("Erro ao decodificar JSON da mensagem (offset %d): %v", msg.Offset, err)
			continue
		}

		// Prepara os dados para o InfluxDB
		tags := map[string]string{
			"car_index": string(payload.CarIndex),
		}
		fields := map[string]interface{}{
			"speed":                 payload.Speed,
			"throttle":              payload.Throttle,
			"brake":                 payload.Brake,
			"gear":                  payload.Gear,
			"engine_rpm":            payload.EngineRPM,
			"drs":                   payload.DRS,
			"tyres_pressure_RL":     payload.TyresPressure[0],
			"tyres_pressure_RR":     payload.TyresPressure[1],
			"tyres_pressure_FL":     payload.TyresPressure[2],
			"tyres_pressure_FR":     payload.TyresPressure[3],
			"tyres_surface_temp_RL": payload.TyresSurfaceTemperature[0],
			"tyres_surface_temp_RR": payload.TyresSurfaceTemperature[1],
			"tyres_surface_temp_FL": payload.TyresSurfaceTemperature[2],
			"tyres_surface_temp_FR": payload.TyresSurfaceTemperature[3],
		}

		// Escreve os dados no InfluxDB
		measurement := "car_telemetry"
		if err := kc.InfluxHandler.WritePoint(measurement, tags, fields, time.Now()); err != nil {
			utils.Logger.Printf("Erro ao salvar dados no InfluxDB (offset %d): %v", msg.Offset, err)
		} else {
			utils.Logger.Printf("Dados de telemetria (offset %d) salvos no InfluxDB.", msg.Offset)
		}

		// Faz o "commit" da mensagem para o Kafka não a enviar novamente.
		if err := kc.Reader.CommitMessages(ctx, msg); err != nil {
			utils.Logger.Printf("Erro ao fazer commit do offset %d: %v", msg.Offset, err)
		}
	}
}

// Close fecha a conexão do leitor Kafka.
func (kc *KafkaConsumer) Close() error {
	return kc.Reader.Close()
}
