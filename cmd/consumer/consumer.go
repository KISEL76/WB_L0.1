package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"wb_test/internal/cache"
	"wb_test/internal/handler"
	"wb_test/internal/kafka"
	"wb_test/internal/storage"
)

const (
	topic         = "orders"
	consumerGroup = "my-consumer-group"
)

func main() {
	address := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	store, err := storage.New(context.Background(), os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	c := cache.New()
	h := handler.NewDBHandler(store, c)

	consumer, err := kafka.NewConsumer(h, address, topic, consumerGroup)
	if err != nil {
		log.Fatal(err)
	}

	go consumer.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	if err := consumer.Stop(); err != nil {
		log.Printf("consumer stop error: %v", err)
	}
}
