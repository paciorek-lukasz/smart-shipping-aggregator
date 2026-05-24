package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

type shippingService interface {
	FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetOptionsResponse, map[string]string)
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
	if req.Method == http.MethodOptions {
		rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		rw.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))

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

	resp, errs := s.service.FetchQuotes(req.Context(), &domainReq)
	if errs != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errs)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(resp)
}
