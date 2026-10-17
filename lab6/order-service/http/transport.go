package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/order-service/transport"
)

// NewService связывает эндпоинты Go kit с HTTP-транспортом
func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/orders").Handler(kithttp.NewServer(
		endpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/orders/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PATCH").Path("/orders/{id}/status").Handler(kithttp.NewServer(
		endpoints.ChangeStatus,
		decodeChangeStatusRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Order); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.GetByIDRequest{ID: vars["id"]}, nil
}

func decodeChangeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return transport.ChangeStatusRequest{ID: vars["id"], Status: body.Status}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
