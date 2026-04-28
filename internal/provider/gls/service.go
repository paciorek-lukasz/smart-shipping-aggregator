package gls

import (
	"context"
	"fmt"
	"time"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/gls/client"
)

const (
	timeLayout  = time.RFC3339
	carrierName = "gls"
)

//go:generate mockgen -destination mocks/generated.go --package mocks --source service.go

type apiClient interface {
	GetQuotesHome(ctx context.Context, req *client.GlsHomeApiRequest) (*client.GlsHomeApiResponse, error)
	GetQuotesPickup(ctx context.Context, req *client.GlsPickupApiRequest) (*client.GlsPickupApiResponse, error)
}

type Service struct {
	apiClient    apiClient
	authToken    string
	depotsLimit  int32
	searchAreaKm int32
}

func NewService(apiClient apiClient, apiKey string, depotsLimit, searchAreaKm int32) *Service {
	return &Service{
		apiClient:    apiClient,
		authToken:    apiKey,
		depotsLimit:  depotsLimit,
		searchAreaKm: searchAreaKm,
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

	apiReq := &client.GlsHomeApiRequest{
		GlsAuthToken: s.authToken,
		ShipperData: &client.ContactData{
			AddressLine1: senderAddress.Address,
			PostalCode:   senderAddress.PostalCode,
			TownCity:     senderAddress.City,
			IsoCountry:   senderAddress.Country,
		},
		ConsigneeData: &client.ContactData{
			AddressLine1: recipientAddress.Address,
			PostalCode:   recipientAddress.PostalCode,
			TownCity:     recipientAddress.City,
			IsoCountry:   recipientAddress.Country,
		},
	}

	apiCtx := context.WithValue(ctx, carrierName, "get_quotes_home")

	resp, err := s.apiClient.GetQuotesHome(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.WindowFrom, resp.WindowTo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.RateCents,
			Currency:          resp.Currency,
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
