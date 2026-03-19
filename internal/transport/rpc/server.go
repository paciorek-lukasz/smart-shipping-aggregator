package rpc

import (
	"net"

	"google.golang.org/grpc"

	pb "github.com/dzwiedz90/smart-shipping-aggregator/api/shipping"
)

func RunServer(port string, handler *Handler) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, handler)

	return s.Serve(lis)
}
