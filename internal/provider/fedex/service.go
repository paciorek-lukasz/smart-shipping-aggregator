package fedex

import (
	"context"
	"fmt"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/fedex/client"
)

const (
	timeLayout  = "2006-01-02T5:04:05"
	carrierName = "fedex"
)

type apiClient interface {
	GetQuotesHome(ctx context.Context, req *client.FedexHomeApiRequest) (*client.FedexHomeApiResponse, error)
	GetQuotesPickup(ctx context.Context, req *client.FedexPickupApiRequest) (*client.FedexPickupApiResponse, error)
}

type Service struct {
	apiClient         apiClient
	clientId          string
	droppointsLimit   int32
	withinRadiusMiles int32
}

func NewService(apiClient apiClient, apiKey string, droppointsLimit, withinRadiusMiles int32) *Service {
	return &Service{
		apiClient:         apiClient,
		clientId:          apiKey,
		droppointsLimit:   droppointsLimit,
		withinRadiusMiles: withinRadiusMiles,
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

	apiReq := &client.FedexHomeApiRequest{
		FedexClientId: s.clientId,
		SourceDetails: &client.ShipperConsignee{
			Street1:     senderAddress.Address,
			PostalCode:  senderAddress.PostalCode,
			City:        senderAddress.City,
			CountryCode: senderAddress.Country,
		},
		DestinationDetails: &client.ShipperConsignee{
			Street1:     recipientAddress.Address,
			PostalCode:  recipientAddress.PostalCode,
			City:        recipientAddress.City,
			CountryCode: recipientAddress.Country,
		},
	}

	apiCtx := context.WithValue(ctx, nil, "get_quotes_home")

	resp, err := s.apiClient.GetQuotesHome(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.EstimatedPickup, resp.EstimatedDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.NetAmount,
			Currency:          resp.CurrencyCode,
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
