// --- main.go ---
package main

import (
	"joselucas/f1-telemetry/src/messaging"
	"joselucas/f1-telemetry/src/telemetry"
	"joselucas/f1-telemetry/src/utils"
	"log"
)

func main() {
	utils.Logger.Println("Iniciando coletor de telemetria F1...")

	producer := messaging.NewKafkaProducer([]string{"host.docker.internal:9092"}, "f1-telemetry")
	defer producer.Close()

	conn, err := telemetry.ListenUDP(20777)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buffer := make([]byte, 2048)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			utils.Logger.Printf("Erro ao ler pacote UDP: %v", err)
			continue
		}
		telemetry.ParseAndPublishPacket(buffer[:n], producer)
	}
}
