package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

func TestGetQuotes_NotPost(t *testing.T) {
	mockService := &mockShippingService{
		result: validResponse(),
	}

	server := NewHttpServer(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/quotes", nil)
	rw := httptest.NewRecorder()

	server.GetQuotes(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rw.Code)
	}
}

func TestGetQuotes_EmptyBody(t *testing.T) {
	mockService := &mockShippingService{
		result: validResponse(),
	}

	server := NewHttpServer(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/quotes", bytes.NewReader(nil))
	rw := httptest.NewRecorder()

	server.GetQuotes(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetQuotes_InvalidJSON(t *testing.T) {
	mockService := &mockShippingService{
		result: validResponse(),
	}

	server := NewHttpServer(mockService)

	invalidJSON := []byte(`{invalid json`)
	req := httptest.NewRequest(http.MethodPost, "/api/quotes", bytes.NewReader(invalidJSON))
	rw := httptest.NewRecorder()

	server.GetQuotes(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetQuotes_Success(t *testing.T) {
	mockService := &mockShippingService{
		result: validResponse(),
	}

	server := NewHttpServer(mockService)

	reqBody := validRequest()
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/quotes", bytes.NewReader(reqJSON))
	rw := httptest.NewRecorder()

	server.GetQuotes(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rw.Code)
	}

	resp := &domain.GetOptionsResponse{}
	if err := json.Unmarshal(rw.Body.Bytes(), resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(resp.Options) != 1 {
		t.Errorf("expected 1 option, got %d", len(resp.Options))
	}
}

func TestGetQuotes_NilResponseFromService(t *testing.T) {
	mockService := &mockShippingService{
		result: nil,
	}

	server := NewHttpServer(mockService)

	reqBody := validRequest()
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/quotes", bytes.NewReader(reqJSON))
	rw := httptest.NewRecorder()

	server.GetQuotes(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rw.Code)
	}
}

func TestGetQuotes_ValidationError(t *testing.T) {
	mockService := &mockShippingService{
		errs: map[string]string{"sender.name": "is required"},
	}

	server := NewHttpServer(mockService)

	reqBody := validRequest()
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/quotes", bytes.NewReader(reqJSON))
	rw := httptest.NewRecorder()

	server.GetQuotes(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}

	var errs map[string]string
	if err := json.Unmarshal(rw.Body.Bytes(), &errs); err != nil {
		t.Errorf("failed to unmarshal error response: %v", err)
	}

	if errs["sender.name"] != "is required" {
		t.Errorf("expected sender.name error, got %v", errs)
	}
}

type mockShippingService struct {
	result *domain.GetOptionsResponse
	errs   map[string]string
}

func (m *mockShippingService) FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetOptionsResponse, map[string]string) {
	if m.errs != nil {
		return nil, m.errs
	}
	return m.result, nil
}

func validRequest() *domain.GetQuotesRequest {
	return &domain.GetQuotesRequest{
		Sender: &domain.Party{
			Name: "Sender",
			Address: &domain.Address{
				Address:    "123 Main St",
				PostalCode: "12345",
				City:       "New York",
				Country:    "US",
			},
			Email: "sender@test.com",
			Phone: "123456789",
		},
		Recipient: &domain.Party{
			Name: "Recipient",
			Address: &domain.Address{
				Address:    "456 Elm St",
				PostalCode: "67890",
				City:       "Los Angeles",
				Country:    "US",
			},
			Email: "recipient@test.com",
			Phone: "987654321",
		},
		DeliveryType: domain.DELIVERY_TYPE_HOME_DELIVERY,
	}
}

func validResponse() *domain.GetOptionsResponse {
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
