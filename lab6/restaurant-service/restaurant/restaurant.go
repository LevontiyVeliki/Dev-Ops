// Пакет restaurant — доменный слой сервиса «Рестораны»: сущность и интерфейсы.
package restaurant

import "context"

// Restaurant — ресторан.
type Restaurant struct {
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Rating  float32 `json:"rating"`
	Open    bool    `json:"open"`
}

// Service описывает бизнес-логику сервиса ресторанов.
type Service interface {
	Create(ctx context.Context, restaurant Restaurant) (string, error)
	GetByID(ctx context.Context, id string) (Restaurant, error)
	ChangeRating(ctx context.Context, id string, rating float32) error
}

// Repository описывает хранилище ресторанов.
type Repository interface {
	CreateRestaurant(ctx context.Context, restaurant Restaurant) error
	GetRestaurantByID(ctx context.Context, id string) (Restaurant, error)
	ChangeRestaurantRating(ctx context.Context, id string, rating float32) error
}

// MinRating и MaxRating — допустимый диапазон рейтинга.
const (
	MinRating float32 = 0
	MaxRating float32 = 5
)
