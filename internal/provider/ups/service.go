package ups

import (
	"context"
	"fmt"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/ups/client"
)

const (
	timeLayout  = time.RFC3339
	carrierName = "ups"
)

type apiClient interface {
	GetQuotesHome(ctx context.Context, req *client.UpsHomeApiRequest) (*client.UpsHomeApiResponse, error)
	GetQuotesPickup(ctx context.Context, req *client.UpsPickupApiRequest) (*client.UpsPickupApiResponse, error)
}

type Service struct {
	apiClient       apiClient
	apiKey          string
	dropPointsLimit int32
	searchRadius    int32
}

//go:generate mockgen -destination mocks/generated.go --package mocks --source service.go

func NewService(apiClient apiClient, apiKey string, dropPointsLimit, searchRadius int32) *Service {
	return &Service{
		apiClient:       apiClient,
		apiKey:          apiKey,
		dropPointsLimit: dropPointsLimit,
		searchRadius:    searchRadius,
	}
}

func (s *Service) Name() string {
	return carrierName
}

func (s *Service) GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	if req.DeliveryType == domain.DELIVERY_TYPE_PICKUP {
		return s.sendPickupRequest(ctx, req)
	}

	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.UpsHomeApiRequest{
		UpsApiKey: s.apiKey,
		SenderAddress: &client.AddressInfo{
			Street:      senderAddress.Address,
			PostalCode:  senderAddress.PostalCode,
			CityName:    senderAddress.City,
			CountryCode: senderAddress.Country,
		},
		RecipientInfo: &client.AddressInfo{
			Street:      recipientAddress.Address,
			PostalCode:  recipientAddress.PostalCode,
			CityName:    recipientAddress.City,
			CountryCode: recipientAddress.Country,
		},
	}

	apiCtx := context.WithValue(ctx, carrierName, "get_quotes_home")

	resp, err := s.apiClient.GetQuotesHome(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.AvailableFrom, resp.AvailableTo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.FeeAmount,
			Currency:          resp.FeeCurrency,
			DeliveryTimeSlots: timeslots,
			DeliveryType:      domain.DELIVERY_TYPE_HOME_DELIVERY,
		},
	}, nil
}

func parseTimeslots(start, end string) ([]*domain.DeliveryTimeSlot, error) {
	s, err := time.Parse(timeLayout, start)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start timeslot: %w", err)
	}

	e, err := time.Parse(timeLayout, end)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end timeslot: %w", err)
	}

	return []*domain.DeliveryTimeSlot{
		{
			Start: s,
			End:   e,
		},
	}, nil
}
