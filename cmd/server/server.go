package main

import (
	"context"
	"log"
	"os"
	"wb_test/internal/cache"
	"wb_test/internal/server"
	"wb_test/internal/storage"
)

func main() {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		log.Fatal("DSN is not setted")
	}

	store, err := storage.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	c := cache.New()
	server := server.New(c, store)

	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
