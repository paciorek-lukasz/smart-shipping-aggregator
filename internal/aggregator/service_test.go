package aggregator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider"
)

func TestService_FetchQuotes_Success(t *testing.T) {
	mockProvider := &mockProvider{
		name:    "test-provider",
		result:  validQuotesResponse(),
		wantErr: nil,
	}

	svc := NewService([]provider.Provider{mockProvider}, time.Second)
	resp := svc.FetchQuotes(context.Background(), validRequest())

	if resp == nil || len(resp.Options) != 1 {
		t.Errorf("expected 1 option, got %d", len(resp.Options))
	}
}

func TestService_FetchQuotes_ProviderError(t *testing.T) {
	mockProvider := &mockProvider{
		name:    "test-provider",
		result:  nil,
		wantErr: errors.New("api error"),
	}

	svc := NewService([]provider.Provider{mockProvider}, time.Second)
	resp := svc.FetchQuotes(context.Background(), validRequest())

	if resp == nil || len(resp.Options) != 0 {
		t.Errorf("expected 0 options on error, got %d", len(resp.Options))
	}
}

func TestService_FetchQuotes_MultipleProviders(t *testing.T) {
	mockProvider1 := &mockProvider{
		name:    "provider-1",
		result:  validQuotesResponse(),
		wantErr: nil,
	}
	mockProvider2 := &mockProvider{
		name:    "provider-2",
		result:  validQuotesResponse(),
		wantErr: nil,
	}

	svc := NewService([]provider.Provider{mockProvider1, mockProvider2}, time.Second)
	resp := svc.FetchQuotes(context.Background(), validRequest())

	if resp == nil || len(resp.Options) != 2 {
		t.Errorf("expected 2 options, got %d", len(resp.Options))
	}
}

func TestService_FetchQuotes_Timeout(t *testing.T) {
	mockProvider := &mockProvider{
		name:    "slow-provider",
		result:  nil,
		wantErr: errors.New("timeout"),
	}

	svc := NewService([]provider.Provider{mockProvider}, 10*time.Millisecond)
	resp := svc.FetchQuotes(context.Background(), validRequest())

	if resp == nil || len(resp.Options) != 0 {
		t.Errorf("expected 0 options on timeout, got %d", len(resp.Options))
	}
}

func TestService_FetchQuotes_NoProviders(t *testing.T) {
	svc := NewService([]provider.Provider{}, time.Second)
	resp := svc.FetchQuotes(context.Background(), validRequest())

	if resp == nil || len(resp.Options) != 0 {
		t.Errorf("expected 0 options, got %d", len(resp.Options))
	}
}

type mockProvider struct {
	name       string
	result    *domain.GetQuotesResponse
	wantErr   error
	callCount int
}

func (m *mockProvider) GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	m.callCount++
	if m.wantErr != nil {
		return nil, m.wantErr
	}
	return m.result, nil
}

func (m *mockProvider) Name() string {
	return m.name
}

func validRequest() *domain.GetQuotesRequest {
	return &domain.GetQuotesRequest{
		Sender: &domain.Party{
			Name:  "Sender",
			Email: "sender@test.com",
			Address: &domain.Address{
				Address:    "123 Main St",
				PostalCode: "12345",
				City:      "New York",
				Country:   "US",
			},
		},
		Recipient: &domain.Party{
			Name:  "Recipient",
			Email: "recipient@test.com",
			Address: &domain.Address{
				Address:    "456 Elm St",
				PostalCode: "67890",
				City:      "Los Angeles",
				Country:   "US",
			},
		},
		DeliveryType: domain.DELIVERY_TYPE_HOME_DELIVERY,
	}
}

func validQuotesResponse() *domain.GetQuotesResponse {
	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct: "test-carrier",
			Price:         1000,
			Currency:      "USD",
			DeliveryType: domain.DELIVERY_TYPE_HOME_DELIVERY,
		},
	}
}