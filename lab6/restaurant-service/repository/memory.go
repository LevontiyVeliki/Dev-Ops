// Пакет repository — in-memory реализация хранилища на map.
package repository

import (
	"context"
	"sync"

	"example.com/restaurant-service/restaurant"
)

type memoryRepo struct {
	mu          sync.RWMutex
	restaurants map[string]restaurant.Restaurant
}

// NewMemoryRepository создаёт in-memory хранилище ресторанов.
func NewMemoryRepository() restaurant.Repository {
	return &memoryRepo{restaurants: make(map[string]restaurant.Restaurant)}
}

func (r *memoryRepo) CreateRestaurant(_ context.Context, rest restaurant.Restaurant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.restaurants[rest.ID] = rest
	return nil
}

func (r *memoryRepo) GetRestaurantByID(_ context.Context, id string) (restaurant.Restaurant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rest, ok := r.restaurants[id]
	if !ok {
		return restaurant.Restaurant{}, restaurant.ErrRestaurantNotFound
	}
	return rest, nil
}

func (r *memoryRepo) ChangeRestaurantRating(_ context.Context, id string, rating float32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rest, ok := r.restaurants[id]
	if !ok {
		return restaurant.ErrRestaurantNotFound
	}
	rest.Rating = rating
	r.restaurants[id] = rest
	return nil
}
