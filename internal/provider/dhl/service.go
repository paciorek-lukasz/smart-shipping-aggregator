package dhl

import (
	"context"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dhl/dhlclient"
)

const (
	timeLayout  = "2006-01-02T5:04:05"
	carrierName = "dhl"
)

type apiClient interface {
	GetQuotesHome(ctx context.Context, req *dhlclient.DhlHomeApiRequest) (*dhlclient.DhlHomeApiResponse, error)
	GetQuotesPickup(ctx context.Context, req *dhlclient.DhlPickupApiRequest) (*dhlclient.DhlPickupApiResponse, error)
}

type Service struct {
	apiClient      apiClient
	apiKey         string
	locationsLimit int32
	searchRadius   int32
}

func NewService(apiClient apiClient, locationsLimit, searchRadius int32) *Service {
	return &Service{
		apiClient:      apiClient,
		locationsLimit: locationsLimit,
		searchRadius:   searchRadius,
	}
}

func (s *Service) GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	if req.DeliveryType == domain.DELIVERY_TYPE_PICKUP {
		return s.sendPickupRequest(ctx, req)
	}

	apiReq := &dhlclient.DhlHomeApiRequest{
		DhlApiKey: s.apiKey,
		SenderAddress: &dhlclient.Party{
			Address:    "",
			PostalCode: "",
			City:       "",
			Country:    "",
		},
		RecipientAddress: &dhlclient.Party{},
	}
}
