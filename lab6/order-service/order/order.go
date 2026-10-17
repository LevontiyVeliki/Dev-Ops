package order

import "context"

// Order представляет заказ
type Order struct {
	ID           string      `json:"id,omitempty"`
	CustomerID   string      `json:"customer_id"`
	Status       string      `json:"status"`
	CreatedOn    int64       `json:"created_on,omitempty"`
	RestaurantID string      `json:"restaurant_id"`
	OrderItems   []OrderItem `json:"order_items,omitempty"`
}

// OrderItem представляет позицию заказа
type OrderItem struct {
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

// Service описывает бизнес-логику сервиса заказов
type Service interface {
	Create(ctx context.Context, order Order) (string, error)
	GetByID(ctx context.Context, id string) (Order, error)
	ChangeStatus(ctx context.Context, id string, status string) error
}

// Repository описывает хранилище заказов
type Repository interface {
	CreateOrder(ctx context.Context, order Order) error
	GetOrderByID(ctx context.Context, id string) (Order, error)
	ChangeOrderStatus(ctx context.Context, id string, status string) error
}
