// Пакет implementation — сервисный слой: бизнес-правила без знания о транспорте.
package implementation

import (
	"context"
	"errors"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"

	"example.com/restaurant-service/restaurant"
)

type service struct {
	repository restaurant.Repository
	logger     log.Logger
}

// NewService создаёт сервис ресторанов.
func NewService(repo restaurant.Repository, logger log.Logger) restaurant.Service {
	return &service{repository: repo, logger: logger}
}

func validRating(r float32) bool {
	return r >= restaurant.MinRating && r <= restaurant.MaxRating
}

// Create проверяет данные, присваивает ID и сохраняет ресторан.
func (s *service) Create(ctx context.Context, r restaurant.Restaurant) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return "", restaurant.ErrInvalidName
	}
	if !validRating(r.Rating) {
		return "", restaurant.ErrInvalidRating
	}
	r.ID = uuid.New().String()

	if err := s.repository.CreateRestaurant(ctx, r); err != nil {
		level.Error(logger).Log("err", err)
		return "", restaurant.ErrCmdRepository
	}
	level.Info(logger).Log("id", r.ID, "name", r.Name)
	return r.ID, nil
}

// GetByID возвращает ресторан по идентификатору.
func (s *service) GetByID(ctx context.Context, id string) (restaurant.Restaurant, error) {
	logger := log.With(s.logger, "method", "GetByID")

	r, err := s.repository.GetRestaurantByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("id", id, "err", err)
		return r, restaurant.ErrRestaurantNotFound
	}
	return r, nil
}

// ChangeRating — доменный метод варианта 5: изменить рейтинг ресторана.
func (s *service) ChangeRating(ctx context.Context, id string, rating float32) error {
	logger := log.With(s.logger, "method", "ChangeRating")

	if !validRating(rating) {
		return restaurant.ErrInvalidRating
	}
	if err := s.repository.ChangeRestaurantRating(ctx, id, rating); err != nil {
		level.Error(logger).Log("id", id, "err", err)
		if errors.Is(err, restaurant.ErrRestaurantNotFound) {
			return restaurant.ErrRestaurantNotFound
		}
		return restaurant.ErrCmdRepository
	}
	level.Info(logger).Log("id", id, "rating", rating)
	return nil
}
