package implementation

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"

	"example.com/order-service/order"
)

type service struct {
	repository order.Repository
	logger     log.Logger
}

// NewService создаёт новый экземпляр сервиса заказов
func NewService(repo order.Repository, logger log.Logger) order.Service {
	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, o order.Order) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	id := uuid.New().String()
	o.ID = id
	o.Status = "Pending"
	o.CreatedOn = time.Now().Unix()

	if err := s.repository.CreateOrder(ctx, o); err != nil {
		level.Error(logger).Log("err", err)
		return "", order.ErrCmdRepository
	}
	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (order.Order, error) {
	logger := log.With(s.logger, "method", "GetByID")

	o, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return o, order.ErrOrderNotFound
	}
	return o, nil
}

func (s *service) ChangeStatus(ctx context.Context, id string, status string) error {
	logger := log.With(s.logger, "method", "ChangeStatus")

	if err := s.repository.ChangeOrderStatus(ctx, id, status); err != nil {
		level.Error(logger).Log("err", err)
		return order.ErrCmdRepository
	}
	return nil
}
