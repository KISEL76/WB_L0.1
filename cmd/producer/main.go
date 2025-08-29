package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
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

	trim := strings.TrimLeftFunc(string(data), func(r rune) bool { return r <= ' ' })
	if len(trim) == 0 {
		log.Fatal("[ERROR] file is empty")
	}

	if trim[0] == '[' {
		var orders []model.Order
		if err := json.Unmarshal(data, &orders); err != nil {
			log.Fatalf("[ERROR] Can't unmarshal orders array: %v", err)
		}
		for i := range orders {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := producer.Produce(ctx, &orders[i], topic); err != nil {
				cancel()
				log.Fatalf("[ERROR] produce failed for index %d (order_uid=%s): %v", i, orders[i].OrderUID, err)
			}
			cancel()
		}
		log.Printf("✅ %d orders sent to Kafka topic %q", len(orders), topic)
		return
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
