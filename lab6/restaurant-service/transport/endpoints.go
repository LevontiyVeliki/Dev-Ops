// Пакет transport — уровень эндпоинтов Go kit и DTO запросов/ответов.
package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/restaurant-service/restaurant"
)

// DTO запросов и ответов. Бизнес-ошибка передаётся клиенту строкой Error:
// значение интерфейса error кодируется в JSON как пустой объект {}.
// Исходная ошибка хранится в неэкспортируемом поле err и нужна
// транспорту, чтобы выбрать HTTP-код ответа.
type CreateRequest struct {
	Restaurant restaurant.Restaurant
}

type CreateResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	err   error
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Restaurant *restaurant.Restaurant `json:"restaurant,omitempty"`
	Error      string                 `json:"error,omitempty"`
	err        error
}

type ChangeRatingRequest struct {
	ID     string  `json:"id"`
	Rating float32 `json:"rating"`
}

type ChangeRatingResponse struct {
	Error string `json:"error,omitempty"`
	err   error
}

// Failer реализуют ответы, которые могут содержать бизнес-ошибку.
type Failer interface {
	Failed() error
}

func (r CreateResponse) Failed() error       { return r.err }
func (r GetByIDResponse) Failed() error      { return r.err }
func (r ChangeRatingResponse) Failed() error { return r.err }

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// Endpoints — набор эндпоинтов сервиса ресторанов.
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeRating endpoint.Endpoint
}

// MakeEndpoints превращает каждый метод сервиса в endpoint.Endpoint.
func MakeEndpoints(s restaurant.Service) Endpoints {
	return Endpoints{
		Create:       makeCreateEndpoint(s),
		GetByID:      makeGetByIDEndpoint(s),
		ChangeRating: makeChangeRatingEndpoint(s),
	}
}

func makeCreateEndpoint(s restaurant.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Restaurant)
		return CreateResponse{ID: id, Error: errText(err), err: err}, nil
	}
}

func makeGetByIDEndpoint(s restaurant.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		r, err := s.GetByID(ctx, req.ID)
		if err != nil {
			return GetByIDResponse{Error: errText(err), err: err}, nil
		}
		return GetByIDResponse{Restaurant: &r}, nil
	}
}

func makeChangeRatingEndpoint(s restaurant.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeRatingRequest)
		err := s.ChangeRating(ctx, req.ID, req.Rating)
		return ChangeRatingResponse{Error: errText(err), err: err}, nil
	}
}
