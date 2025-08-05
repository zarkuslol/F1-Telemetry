package handlers

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InfluxDBHandler armazena a configuração para interagir com o InfluxDB.
type InfluxDBHandler struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
}

// NewInfluxDBHandler cria um novo handler para o InfluxDB.
func NewInfluxDBHandler(url, token, org, bucket string) *InfluxDBHandler {
	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)
	return &InfluxDBHandler{
		client:   client,
		writeAPI: writeAPI,
	}
}

// WritePoint cria um ponto de dados e o escreve no InfluxDB.
// 'measurement' é o nome da "tabela", 'tags' são índices e 'fields' são os valores.
func (h *InfluxDBHandler) WritePoint(measurement string, tags map[string]string, fields map[string]interface{}, ts time.Time) error {
	p := influxdb2.NewPoint(measurement, tags, fields, ts)

	// Usa um contexto com timeout para a operação de escrita.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.writeAPI.WritePoint(ctx, p); err != nil {
		return fmt.Errorf("falha ao escrever ponto no InfluxDB: %w", err)
	}
	return nil
}

// Close fecha o cliente do InfluxDB.
func (h *InfluxDBHandler) Close() {
	h.client.Close()
}
