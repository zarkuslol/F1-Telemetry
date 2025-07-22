// --- telemetry/kafka_producer.go ---
package messaging

import (
	"context"
	"encoding/json"
	"joselucas/f1-telemetry/src/utils"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		Async:   true,
	})
	utils.Logger.Printf("Produtor Kafka inicializado no tópico '%s'", topic)
	return &KafkaProducer{Writer: writer}
}

func (kp *KafkaProducer) PublishTelemetry(data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   nil,
		Value: jsonData,
	}
	return kp.Writer.WriteMessages(ctx, msg)
}

func (kp *KafkaProducer) Close() error {
	return kp.Writer.Close()
}
