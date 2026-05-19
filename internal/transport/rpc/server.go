package rpc

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/paciorek-lukasz/smart-shipping-aggregator/api/shipping"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

type shippingService interface {
	FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetOptionsResponse, map[string]string)
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

	quotes, errs := h.service.FetchQuotes(ctx, domainReq)
	if errs != nil {
		st := status.New(codes.InvalidArgument, "invalid request parameters")

		br := &errdetails.BadRequest{}

		for field, msg := range errs {
			v := &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: msg,
			}
			br.FieldViolations = append(br.FieldViolations, v)
		}

		st, err := st.WithDetails(br)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to update status with error details")
		}

		return nil, st.Err()
	}

	mapped, err := mapDomainToProto(quotes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map domain models to proto")
	}

	return mapped, nil
}
