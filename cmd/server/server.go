package main

import (
	"context"
	"log"
	"os"

	"wb_test/internal/cache"
	"wb_test/internal/server"
	"wb_test/internal/storage"
)

const (
	cacheWarmupLimit = 10
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
	warmupCache(store, c)

	server := server.New(c, store)

	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func warmupCache(store *storage.Storage, c *cache.Cache) {
	lastOrders, err := store.Orders().GetLastOrders(context.Background(), cacheWarmupLimit)
	if err != nil {
		log.Printf("cache warmup error: %v", err)
	} else {
		for _, o := range lastOrders {
			c.Set(o)
		}
		log.Printf("cache preloaded with %d orders", len(lastOrders))
	}
}
