package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dzwiedz90/smart-shipping-aggregator/api/shipping"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
)

type ShippingService interface {
	GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error)
}

type Handler struct {
	service ShippingService
	pb.UnimplementedShippingServiceServer
}

func NewHandler(svc ShippingService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) GetQuotes(ctx context.Context, req *pb.GetQuotesRequest) (*pb.GetQuotesResponse, error) {
	if req.Sender == nil || req.Recipient == nil {
		return nil, status.Error(codes.InvalidArgument, "sender and recipient are required")
	}

	domainReq := MapProtoToDomain(req)

	quotes, err := h.service.GetQuotes(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch quotes: %v", err)
	}

	mapped, err := MapDomainToProto(quotes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map domain models to proto")
	}

	return mapped, nil
}
