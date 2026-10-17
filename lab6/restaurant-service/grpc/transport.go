// Пакет grpc — gRPC-транспорт того же сервиса на базе
// github.com/go-kit/kit/transport/grpc. Эндпоинты и сервисный слой
// используются те же, что и в HTTP-транспорте.
package grpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"example.com/restaurant-service/pb"
	"example.com/restaurant-service/restaurant"
	"example.com/restaurant-service/transport"
)

type grpcServer struct {
	pb.UnimplementedRestaurantServiceServer
	create       grpctransport.Handler
	getByID      grpctransport.Handler
	changeRating grpctransport.Handler
}

// NewGRPCServer создаёт реализацию pb.RestaurantServiceServer: каждый RPC
// обслуживается своим grpctransport.Handler с функциями декодирования
// protobuf-сообщения в DTO и кодирования DTO обратно.
func NewGRPCServer(endpoints transport.Endpoints, logger log.Logger) pb.RestaurantServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		create: grpctransport.NewServer(
			endpoints.Create, decodeCreateRequest, encodeCreateResponse, options...),
		getByID: grpctransport.NewServer(
			endpoints.GetByID, decodeGetByIDRequest, encodeGetByIDResponse, options...),
		changeRating: grpctransport.NewServer(
			endpoints.ChangeRating, decodeChangeRatingRequest, encodeChangeRatingResponse, options...),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateReply), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDReply, error) {
	_, resp, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDReply), nil
}

func (s *grpcServer) ChangeRating(ctx context.Context, req *pb.ChangeRatingRequest) (*pb.ChangeRatingReply, error) {
	_, resp, err := s.changeRating.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ChangeRatingReply), nil
}

func toPB(r restaurant.Restaurant) *pb.Restaurant {
	return &pb.Restaurant{Id: r.ID, Name: r.Name, Address: r.Address, Rating: r.Rating, Open: r.Open}
}

func fromPB(r *pb.Restaurant) restaurant.Restaurant {
	if r == nil {
		return restaurant.Restaurant{}
	}
	return restaurant.Restaurant{ID: r.Id, Name: r.Name, Address: r.Address, Rating: r.Rating, Open: r.Open}
}

func decodeCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateRequest)
	return transport.CreateRequest{Restaurant: fromPB(req.Restaurant)}, nil
}

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.CreateResponse)
	return &pb.CreateReply{Id: resp.ID, Error: resp.Error}, nil
}

func decodeGetByIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	return transport.GetByIDRequest{ID: request.(*pb.GetByIDRequest).Id}, nil
}

func encodeGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.GetByIDResponse)
	reply := &pb.GetByIDReply{Error: resp.Error}
	if resp.Restaurant != nil {
		reply.Restaurant = toPB(*resp.Restaurant)
	}
	return reply, nil
}

func decodeChangeRatingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ChangeRatingRequest)
	return transport.ChangeRatingRequest{ID: req.Id, Rating: req.Rating}, nil
}

func encodeChangeRatingResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &pb.ChangeRatingReply{Error: response.(transport.ChangeRatingResponse).Error}, nil
}
