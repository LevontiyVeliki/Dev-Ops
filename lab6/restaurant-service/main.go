// Сервис «Рестораны» (вариант 5): Go kit, HTTP на :8080 и gRPC на :9090.
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpctransport "example.com/restaurant-service/grpc"
	httptransport "example.com/restaurant-service/http"
	"example.com/restaurant-service/implementation"
	"example.com/restaurant-service/pb"
	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "restaurant",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Сервисный слой и эндпоинты создаются один раз и общие для транспортов
	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// HTTP-транспорт
	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: httptransport.NewService(endpoints, logger),
		}
		errs <- server.ListenAndServe()
	}()

	// gRPC-транспорт
	go func() {
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		server := grpc.NewServer()
		pb.RegisterRestaurantServiceServer(server, grpctransport.NewGRPCServer(endpoints, logger))
		reflection.Register(server) // описание сервиса для grpcurl
		errs <- server.Serve(listener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
