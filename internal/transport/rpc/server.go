package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/paciorek-lukasz/smart-shipping-aggregator/api/shipping"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

type shippingService interface {
	FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) *domain.GetOptionsResponse
}

type GrpcServer struct {
	service shippingService
	pb.UnimplementedShippingServiceServer
}

func NewGrpcServer(svc shippingService) *GrpcServer {
	return &GrpcServer{service: svc}
}

func (h *GrpcServer) GetQuotes(ctx context.Context, req *pb.GetQuotesRequest) (*pb.GetQuotesResponse, error) {
	if req.Sender == nil || req.Recipient == nil {
		return nil, status.Error(codes.InvalidArgument, "sender and recipient are required")
	}

	domainReq := mapProtoToDomain(req)

	quotes := h.service.FetchQuotes(ctx, domainReq)

	mapped, err := mapDomainToProto(quotes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map domain models to proto")
	}

	return mapped, nil
}
