package repository

import (
	"context"
	"sync"

	"example.com/order-service/order"
)

type memoryRepo struct {
	mu     sync.RWMutex
	orders map[string]order.Order
}

// NewMemoryRepository создаёт in-memory хранилище заказов
func NewMemoryRepository() order.Repository {
	return &memoryRepo{
		orders: make(map[string]order.Order),
	}
}

func (r *memoryRepo) CreateOrder(ctx context.Context, o order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[o.ID] = o
	return nil
}

func (r *memoryRepo) GetOrderByID(ctx context.Context, id string) (order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return order.Order{}, order.ErrOrderNotFound
	}
	return o, nil
}

func (r *memoryRepo) ChangeOrderStatus(ctx context.Context, id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[id]
	if !ok {
		return order.ErrOrderNotFound
	}
	o.Status = status
	r.orders[id] = o
	return nil
}
