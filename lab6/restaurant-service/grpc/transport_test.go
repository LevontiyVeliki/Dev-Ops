package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	grpctransport "example.com/restaurant-service/grpc"
	"example.com/restaurant-service/implementation"
	"example.com/restaurant-service/pb"
	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/transport"
)

// Сервер поднимается в памяти через bufconn — без реального порта.
func newClient(t *testing.T) pb.RestaurantServiceClient {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	logger := log.NewNopLogger()
	svc := implementation.NewService(repository.NewMemoryRepository(), logger)
	server := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(server, grpctransport.NewGRPCServer(transport.MakeEndpoints(svc), logger))
	go server.Serve(lis)
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return pb.NewRestaurantServiceClient(conn)
}

func TestGRPCFlow(t *testing.T) {
	ctx := context.Background()
	client := newClient(t)

	created, err := client.Create(ctx, &pb.CreateRequest{Restaurant: &pb.Restaurant{
		Name: "Пиццерия", Address: "пр. Кулакова, 5", Rating: 4.1}})
	if err != nil || created.Id == "" || created.Error != "" {
		t.Fatalf("Create: %v %v", created, err)
	}
	if r, err := client.ChangeRating(ctx, &pb.ChangeRatingRequest{Id: created.Id, Rating: 4.7}); err != nil || r.Error != "" {
		t.Fatalf("ChangeRating: %v %v", r, err)
	}
	got, err := client.GetByID(ctx, &pb.GetByIDRequest{Id: created.Id})
	if err != nil || got.Restaurant.GetName() != "Пиццерия" || got.Restaurant.GetRating() != 4.7 {
		t.Fatalf("GetByID: %v %v", got, err)
	}
}

func TestGRPCErrors(t *testing.T) {
	ctx := context.Background()
	client := newClient(t)
	if r, _ := client.GetByID(ctx, &pb.GetByIDRequest{Id: "nope"}); r.Error != "restaurant not found" || r.Restaurant != nil {
		t.Errorf("GetByID: %v", r)
	}
	if r, _ := client.Create(ctx, &pb.CreateRequest{}); r.Error != "restaurant name is required" {
		t.Errorf("Create empty: %v", r)
	}
	if r, _ := client.ChangeRating(ctx, &pb.ChangeRatingRequest{Id: "nope", Rating: 9}); r.Error != "rating must be between 0 and 5" {
		t.Errorf("ChangeRating: %v", r)
	}
}
