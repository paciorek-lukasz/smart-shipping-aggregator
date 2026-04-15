package inpost

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/inpost/client"
)

func (s *Service) sendPickupRequest(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.InpostPickupApiRequest{
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
		MachinesLimit:   s.machinesLimit,
		RadiusKm:        s.radiusKm,
		MachineCategory: mapLocationTypes(req.LocationTypes),
	}

	apiCtx := context.WithValue(ctx, nil, "get_quotes_pickup")

	resp, err := s.apiClient.GetQuotesPickup(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.DeliveryStart, resp.DeliveryEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	machines, err := parseMachines(resp.Machines)
	if err != nil {
		return nil, fmt.Errorf("failed to parse machines: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.ShipmentCost,
			Currency:          resp.IsoCurrency,
			DeliveryTimeSlots: timeslots,
			PickupPoints:      machines,
			DeliveryType:      domain.DELIVERY_TYPE_PICKUP,
		},
	}, nil
}

func parseMachines(machines []*client.ParcelLocker) ([]*domain.PickupPoint, error) {
	var result []*domain.PickupPoint

	for _, m := range machines {
		locType, err := mapInpostStatus(m.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to map machine status: %w", err)
		}

		result = append(result, &domain.PickupPoint{
			PickupPointId: m.MachineId,
			Name:          m.Name,
			Address: &domain.Address{
				Address:    m.Address,
				PostalCode: m.ZipCode,
				City:       m.Locality,
				Country:    m.CountryIso,
				Longitude:  m.Longitude,
				Latitude:   m.Latitude,
			},
			LocationType:  locType,
			OpeningHours:  parseHours(m.Hours),
			IsOperational: m.Status == client.InpostStatusActive,
		})
	}

	return result, nil
}

func mapInpostStatus(status string) (domain.LocationType, error) {
	switch status {
	case client.InpostStatusActive:
		return domain.LOCATION_TYPE_LOCKER, nil
	case client.InpostStatusInactive:
		return domain.LOCATION_TYPE_SERVICE_POINT, nil
	}

	return 0, errors.New("no matching status")
}

func parseHours(hours []*client.Opening) []*domain.OpeningHours {
	var result []*domain.OpeningHours

	for _, h := range hours {
		result = append(result, &domain.OpeningHours{
			DayOfWeek: h.DayIndex,
			Opens:     h.OpenTime,
			Closes:    h.CloseTime,
		})
	}

	return result
}

func mapLocationTypes(types []domain.LocationType) []client.InpostCategory {
	var result []client.InpostCategory

	for _, t := range types {
		if t == domain.LOCATION_TYPE_LOCKER {
			result = append(result, client.INPOST_CATEGORY_LOCKER)
		}
		if t == domain.LOCATION_TYPE_SERVICE_POINT {
			result = append(result, client.INPOST_CATEGORY_PARTNER)
		}
	}

	if len(result) == 0 {
		result = append(result, client.INPOST_CATEGORY_ALL)
	}

	return result
}
