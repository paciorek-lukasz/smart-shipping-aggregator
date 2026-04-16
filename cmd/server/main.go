package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	pb "github.com/dzwiedz90/smart-shipping-aggregator/api/shipping"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/aggregator"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/config"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/transport/rpc"
)

const (
	appName = "smart-shipping-aggregator"
)

func main() {
	var cfg config.EnvConfig

	if err := envconfig.Process(appName, &cfg); err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}

	providers := provider.InitProviders(cfg)
	if len(providers) == 0 {
		log.Fatal("no providers initialized")
	}

	aggSvc := aggregator.NewService(providers, cfg.Timeout)

	grpcHandler := rpc.NewHandler(aggSvc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterShippingServiceServer(s, grpcHandler)

	go func() {
		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down gRPC server gracefully...")
	s.GracefulStop()
	log.Println("Server stopped")
}
