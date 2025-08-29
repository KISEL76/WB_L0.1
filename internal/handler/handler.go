package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"wb_test/internal/cache"
	"wb_test/internal/model"
	"wb_test/internal/storage"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Handler interface {
	HandleMessage(message *ckafka.Message, offset ckafka.Offset) error
}

type DBhandler struct {
	orders *storage.OrderStorage
	cache  *cache.Cache
}

func NewDBHandler(store *storage.Storage, cache *cache.Cache) *DBhandler {
	return &DBhandler{
		orders: store.Orders(),
		cache:  cache,
	}
}

func (h *DBhandler) HandleMessage(msg *ckafka.Message, offset ckafka.Offset) error {
	var order model.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		return fmt.Errorf("[ERROR] can't unmarshall data from producer: %v", err)
	}
	if order.OrderUID == "" {
		return fmt.Errorf("[ERROR] empty order_uid")
	}

	ctx := context.Background()
	if err := h.orders.Upsert(ctx, &order); err != nil {
		return err
	}

	if h.cache != nil {
		h.cache.Set(&order)
	}
	return nil
}
