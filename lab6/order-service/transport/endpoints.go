package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/order-service/order"
)

// DTO для запросов и ответов
type CreateRequest struct {
	Order order.Order
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Order order.Order `json:"order"`
	Err   error       `json:"error,omitempty"`
}

type ChangeStatusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ChangeStatusResponse struct {
	Err error `json:"error,omitempty"`
}

// Endpoints — набор эндпоинтов сервиса заказов
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeStatus endpoint.Endpoint
}

// MakeEndpoints создаёт эндпоинты для сервиса заказов
func MakeEndpoints(s order.Service) Endpoints {
	return Endpoints{
		Create:       makeCreateEndpoint(s),
		GetByID:      makeGetByIDEndpoint(s),
		ChangeStatus: makeChangeStatusEndpoint(s),
	}
}

func makeCreateEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Order)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		o, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Order: o, Err: err}, nil
	}
}

func makeChangeStatusEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeStatusRequest)
		err := s.ChangeStatus(ctx, req.ID, req.Status)
		return ChangeStatusResponse{Err: err}, nil
	}
}
