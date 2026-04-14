package dhl

import (
	"context"
	"fmt"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dhl/client"
)

const (
	timeLayout  = "2006-01-02T5:04:05"
	carrierName = "dhl"
)

type apiClient interface {
	GetQuotesHome(ctx context.Context, req *client.DhlHomeApiRequest) (*client.DhlHomeApiResponse, error)
	GetQuotesPickup(ctx context.Context, req *client.DhlPickupApiRequest) (*client.DhlPickupApiResponse, error)
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

	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.DhlHomeApiRequest{
		DhlApiKey: s.apiKey,
		SenderAddress: &client.Party{
			Address:    senderAddress.Address,
			PostalCode: senderAddress.PostalCode,
			City:       senderAddress.City,
			Country:    senderAddress.Country,
		},
		RecipientAddress: &client.Party{
			Address:    recipientAddress.Address,
			PostalCode: recipientAddress.PostalCode,
			City:       recipientAddress.City,
			Country:    recipientAddress.Country,
		},
	}

	apiCtx := context.WithValue(ctx, nil, "get_quotes_home")

	resp, err := s.apiClient.GetQuotesHome(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.Earliest, resp.Latest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: []*domain.Option{
			{
				// OptionId: ,
				CarrierProduct:    carrierName,
				Price:             resp.Price,
				Currency:          resp.Currency,
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
