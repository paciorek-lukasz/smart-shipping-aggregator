package rpc

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/paciorek-lukasz/smart-shipping-aggregator/api/shipping"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

func TestHandler_GetQuotes_Success(t *testing.T) {
	mockService := &mockShippingService{
		result: validDomainResponse(),
	}

	handler := NewGrpcServer(mockService)
	resp, err := handler.GetQuotes(context.Background(), validProtoRequest())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp == nil || len(resp.Options) != 1 {
		t.Errorf("expected 1 option, got nil or %d", len(resp.Options))
	}
}

func TestHandler_GetQuotes_MissingSender(t *testing.T) {
	mockService := &mockShippingService{
		result: validDomainResponse(),
	}

	handler := NewGrpcServer(mockService)
	req := validProtoRequest()
	req.Sender = nil

	_, err := handler.GetQuotes(context.Background(), req)

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestHandler_GetQuotes_MissingRecipient(t *testing.T) {
	mockService := &mockShippingService{
		result: validDomainResponse(),
	}

	handler := NewGrpcServer(mockService)
	req := validProtoRequest()
	req.Recipient = nil

	_, err := handler.GetQuotes(context.Background(), req)

	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestHandler_GetQuotes_ServiceError(t *testing.T) {
	mockService := &mockShippingService{
		err: errors.New("service error"),
	}

	handler := NewGrpcServer(mockService)
	_, err := handler.GetQuotes(context.Background(), validProtoRequest())

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMapProtoToDomain(t *testing.T) {
	pbReq := validProtoRequest()
	domainReq := mapProtoToDomain(pbReq)

	if domainReq.Sender.Name != "Sender" {
		t.Errorf("expected sender name, got %s", domainReq.Sender.Name)
	}
	if domainReq.Recipient.Name != "Recipient" {
		t.Errorf("expected recipient name, got %s", domainReq.Recipient.Name)
	}
}

func TestMapDomainToProto(t *testing.T) {
	resp := &domain.GetOptionsResponse{
		Options: []*domain.Option{
			{
				OptionId:       0,
				CarrierProduct: "test",
				Price:          1000,
				Currency:       "USD",
				DeliveryType:   domain.DELIVERY_TYPE_HOME_DELIVERY,
			},
		},
	}

	pbResp, err := mapDomainToProto(resp)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(pbResp.Options) != 1 {
		t.Errorf("expected 1 option, got %d", len(pbResp.Options))
	}
}

func TestMapDomainToProto_NilResponse(t *testing.T) {
	_, err := mapDomainToProto(nil)

	if err == nil {
		t.Errorf("expected error for nil response, got nil")
	}
}

type mockShippingService struct {
	result *domain.GetOptionsResponse
	err    error
}

func (m *mockShippingService) FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) *domain.GetOptionsResponse {
	if m.err != nil {
		return nil
	}
	return m.result
}

func validProtoRequest() *pb.GetQuotesRequest {
	return &pb.GetQuotesRequest{
		Sender: &pb.Party{
			Name: "Sender",
			Address: &pb.Address{
				Address:    "123 Main St",
				PostalCode: "12345",
				City:       "New York",
				Country:    "US",
			},
			Email: "sender@test.com",
			Phone: "123456789",
		},
		Recipient: &pb.Party{
			Name: "Recipient",
			Address: &pb.Address{
				Address:    "456 Elm St",
				PostalCode: "67890",
				City:       "Los Angeles",
				Country:    "US",
			},
			Email: "recipient@test.com",
			Phone: "987654321",
		},
		DeliveryType: pb.DeliveryType_DELIVERY_TYPE_HOME_DELIVERY,
	}
}

func validDomainResponse() *domain.GetOptionsResponse {
	return &domain.GetOptionsResponse{
		Options: []*domain.Option{
			{
				OptionId:       0,
				CarrierProduct: "test-carrier",
				Price:          1000,
				Currency:       "USD",
				DeliveryType:   domain.DELIVERY_TYPE_HOME_DELIVERY,
			},
		},
	}
}
