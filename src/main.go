package main

import (
	"context"
	"joselucas/f1-telemetry/src/handlers"
	"joselucas/f1-telemetry/src/messaging"
	"joselucas/f1-telemetry/src/telemetry"
	"joselucas/f1-telemetry/src/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	utils.Logger.Println("Iniciando serviço de telemetria F1...")

	// --- Configurações (Valores do seu docker-compose.yml) ---
	kafkaBrokers := []string{"kafka:9092"}
	kafkaTopic := "f1-telemetry"
	kafkaGroupID := "influxdb-writer-group"

	influxURL := "http://influxdb:8086"
	influxToken := "meu-token-secreto"
	influxOrg := "f1-org"
	influxBucket := "f1-bucket"

	// --- Contexto para Gerenciamento de Desligamento (Graceful Shutdown) ---
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup // WaitGroup para esperar as goroutines terminarem

	// --- Inicializa o Handler do InfluxDB ---
	influxHandler := handlers.NewInfluxDBHandler(influxURL, influxToken, influxOrg, influxBucket)
	defer influxHandler.Close()

	// --- Inicia o Consumidor Kafka em uma Goroutine ---
	consumer := messaging.NewKafkaConsumer(kafkaBrokers, kafkaTopic, kafkaGroupID, influxHandler)
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.ConsumeInfluxDB(ctx)
	}()
	defer consumer.Close()

	// --- Inicia o Produtor e o Listener UDP ---
	producer := messaging.NewKafkaProducer(kafkaBrokers, kafkaTopic)
	defer producer.Close()

	conn, err := telemetry.ListenUDP(20777)
	if err != nil {
		log.Fatalf("Falha fatal ao iniciar listener UDP: %v", err)
	}
	defer conn.Close()

	// Inicia o listener UDP em outra goroutine para não bloquear o shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 2048)
		for {
			select {
			case <-ctx.Done(): // Verifica se o contexto foi cancelado
				utils.Logger.Println("Encerrando listener UDP...")
				return
			default:
				n, _, err := conn.ReadFromUDP(buffer)
				if err != nil {
					// Ignora o erro se o contexto foi cancelado
					if ctx.Err() == nil {
						utils.Logger.Printf("Erro ao ler pacote UDP: %v", err)
					}
					continue
				}
				telemetry.ParseAndPublishPacket(buffer[:n], producer)
			}
		}
	}()

	// --- Espera por um sinal de interrupção (Ctrl+C) ---
	utils.Logger.Println("Serviço em execução. Pressione Ctrl+C para sair.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // Bloqueia até receber o sinal

	utils.Logger.Println("Recebido sinal de desligamento. Encerrando serviços...")
	cancel() // Cancela o contexto, sinalizando para as goroutines pararem

	// Espera as goroutines finalizarem de forma limpa
	wg.Wait()
	utils.Logger.Println("Todos os serviços foram encerrados.")
}
