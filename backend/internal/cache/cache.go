package cache

import (
	"L0-wbtech/internal/model"
	"sync"
)

type Cache interface {
	Set(order *model.Order)
	Get(uid string) (*model.Order, bool)
}

type inMemoryCache struct {
	data map[string]*model.Order
	mu   sync.RWMutex
}

func NewCache() Cache {
	return &inMemoryCache{
		data: make(map[string]*model.Order),
	}
}

func (c *inMemoryCache) Set(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

func (c *inMemoryCache) Get(uid string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.data[uid]
	return order, ok
}
