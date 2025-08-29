package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"wb_test/internal/kafka"
	"wb_test/internal/model"
)

const (
	topic    = "orders"
	filename = "sample.json"
)

var (
	address = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
)

func main() {
	producer, err := kafka.NewProducer(address)
	if err != nil {
		log.Fatalf("[ERROR] Didn't manage to create producer: %v", err)
	}
	defer producer.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("[ERROR] Can't read file sample.json: %v", err)
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Fatalf("[ERROR] Can't unmarshal sample order: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := producer.Produce(ctx, &order, topic); err != nil {
		log.Fatalf("[ERROR] Something wrong while producing data: %v", err)
	}

	log.Println("✅ sample order sent to Kafka topic:", topic)
}
