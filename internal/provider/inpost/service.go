package inpost

import (
	"context"
	"fmt"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/inpost/client"
)

const (
	timeLayout  = "2006-01-02T5:04:05"
	carrierName = "inpost"
)

type apiClient interface {
	GetQuotesHome(ctx context.Context, req *client.InpostHomeApiRequest) (*client.InpostHomeApiResponse, error)
	GetQuotesPickup(ctx context.Context, req *client.InpostPickupApiRequest) (*client.InpostPickupApiResponse, error)
}

type Service struct {
	apiClient     apiClient
	token         string
	machinesLimit int32
	radiusKm      int32
}

func NewService(apiClient apiClient, machinesLimit, radiusKm int32) *Service {
	return &Service{
		apiClient:     apiClient,
		machinesLimit: machinesLimit,
		radiusKm:      radiusKm,
	}
}

func (s *Service) GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	if req.DeliveryType == domain.DELIVERY_TYPE_PICKUP {
		return s.sendPickupRequest(ctx, req)
	}

	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.InpostHomeApiRequest{
		InpostToken: s.token,
		SenderDetails: &client.PartyDetails{
			StreetAddress: senderAddress.Address,
			ZipCode:       senderAddress.PostalCode,
			Locality:      senderAddress.City,
			CountryIso:    senderAddress.Country,
		},
		ReceiverDetails: &client.PartyDetails{
			StreetAddress: recipientAddress.Address,
			ZipCode:       recipientAddress.PostalCode,
			Locality:      recipientAddress.City,
			CountryIso:    recipientAddress.Country,
		},
	}

	apiCtx := context.WithValue(ctx, nil, "get_quotes_home")

	resp, err := s.apiClient.GetQuotesHome(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.DeliveryStart, resp.DeliveryEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: []*domain.Option{
			{
				CarrierProduct:    carrierName,
				Price:             resp.ShipmentCost,
				Currency:          resp.IsoCurrency,
				DeliveryTimeSlots: timeslots,
				DeliveryType:      domain.DELIVERY_TYPE_PICKUP,
			},
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
