// Пакет http — HTTP-транспорт: маршруты, декодирование запросов и
// кодирование ответов.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/restaurant-service/restaurant"
	"example.com/restaurant-service/transport"
)

// ErrBadRequest — тело запроса не удалось разобрать.
var ErrBadRequest = errors.New("bad request body")

// NewService связывает эндпоинты Go kit с HTTP-маршрутами.
func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/restaurants").Handler(kithttp.NewServer(
		endpoints.Create, decodeCreateRequest, encodeResponse, options...))

	r.Methods("GET").Path("/restaurants/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID, decodeGetByIDRequest, encodeResponse, options...))

	r.Methods("PATCH").Path("/restaurants/{id}/rating").Handler(kithttp.NewServer(
		endpoints.ChangeRating, decodeChangeRatingRequest, encodeResponse, options...))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Restaurant); err != nil {
		return nil, ErrBadRequest
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return transport.GetByIDRequest{ID: mux.Vars(r)["id"]}, nil
}

func decodeChangeRatingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Rating *float32 `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Rating == nil {
		return nil, ErrBadRequest
	}
	return transport.ChangeRatingRequest{ID: mux.Vars(r)["id"], Rating: *body.Rating}, nil
}

// encodeResponse пишет JSON; если ответ содержит бизнес-ошибку, код ответа
// выбирается по ней (404, 400 и т.д.), иначе 200 или 201 для создания.
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if f, ok := response.(transport.Failer); ok && f.Failed() != nil {
		w.WriteHeader(codeFrom(f.Failed()))
	} else if _, ok := response.(transport.CreateResponse); ok {
		w.WriteHeader(http.StatusCreated)
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}

func codeFrom(err error) int {
	switch {
	case errors.Is(err, restaurant.ErrRestaurantNotFound):
		return http.StatusNotFound
	case errors.Is(err, restaurant.ErrInvalidName),
		errors.Is(err, restaurant.ErrInvalidRating),
		errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
