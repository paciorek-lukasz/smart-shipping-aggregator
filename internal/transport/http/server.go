package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

type shippingService interface {
	FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) *domain.GetOptionsResponse
}

type HttpServer struct {
	service shippingService
}

func NewHttpServer(service shippingService) *HttpServer {
	return &HttpServer{
		service: service,
	}
}

func (s *HttpServer) GetQuotes(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "failed to extract body from the request", http.StatusBadRequest)
		return
	}

	var domainReq domain.GetQuotesRequest
	if err := json.Unmarshal(body, &domainReq); err != nil {
		http.Error(rw, "failed to unmarshal JSON data from body", http.StatusBadRequest)
		return
	}

	resp := s.service.FetchQuotes(req.Context(), &domainReq)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(resp)
}
