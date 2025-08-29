package cache

import (
	"sync"

	"wb_test/internal/model"
)

type Cache struct {
	mu sync.RWMutex
	m  map[string]*model.Order
}

func New() *Cache { return &Cache{m: make(map[string]*model.Order)} }

func (c *Cache) Get(id string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.m[id]
	return order, ok
}

func (c *Cache) Set(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[order.OrderUID] = order
}
