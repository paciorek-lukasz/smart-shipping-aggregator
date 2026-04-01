package dhl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dhl/dhlclient"
)

func (s *Service) sendPickupRequest(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &dhlclient.DhlPickupApiRequest{
		DhlApiKey: s.apiKey,
		SenderAddress: &dhlclient.Party{
			Address:    senderAddress.Address,
			PostalCode: senderAddress.PostalCode,
			City:       senderAddress.City,
			Country:    senderAddress.Country,
		},
		RecipientAddress: &dhlclient.Party{
			Address:    recipientAddress.Address,
			PostalCode: recipientAddress.PostalCode,
			City:       recipientAddress.City,
			Country:    recipientAddress.Country,
		},
		LocationsLimit: s.locationsLimit,
		SearchRadius:   s.searchRadius,
		LocationsType:  mapLocationTypes(req.LocationTypes),
	}
	apiCtx := context.WithValue(ctx, nil, "get_quotes_pickup")

	resp, err := s.apiClient.GetQuotesPickup(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.Earliest, resp.Latest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	locations, err := parsePickupPoints(resp.Locations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse locations: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: []*domain.Option{
			{
				// OptionId: ,
				CarrierProduct:    carrierName,
				Price:             resp.Price,
				Currency:          resp.Currency,
				DeliveryTimeSlots: timeslots,
				PickupPoints:      locations,
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

func parsePickupPoints(locs []*dhlclient.Location) ([]*domain.PickupPoint, error) {
	var result []*domain.PickupPoint

	for _, loc := range locs {
		locType, err := mapDhlLocType(loc.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to map location type: %w", err)
		}

		result = append(result, &domain.PickupPoint{
			PickupPointId: loc.Id,
			Name:          loc.Name,
			Address: &domain.Address{
				Address:    loc.AddressLine,
				PostalCode: loc.PostalCode,
				City:       loc.City,
				Country:    loc.Country,
				Longitude:  loc.Longitude,
				Latitude:   loc.Latitude,
			},
			LocationType:  locType,
			OpeningHours:  parseOpeningHours(loc.OpenTimes),
			IsOperational: loc.IsAvailable,
		})
	}

	return result, nil
}

func mapDhlLocType(dhlLocType string) (domain.LocationType, error) {
	switch dhlLocType {
	case dhlclient.DhlLocTypeLocker:
		return domain.LOCATION_TYPE_LOCKER, nil
	case dhlclient.DhlLocTypeServicePoint:
		return domain.LOCATION_TYPE_SERVICE_POINT, nil
	}

	return 0, errors.New("no matching location type")
}

func parseOpeningHours(openingHours []*dhlclient.OpenTimes) []*domain.OpeningHours {
	var result []*domain.OpeningHours

	for _, h := range openingHours {
		result = append(result, &domain.OpeningHours{
			DayOfWeek: h.DayOfWeek,
			Opens:     h.Opens,
			Closes:    h.Closes,
		})
	}

	return result
}

func mapLocationTypes(types []domain.LocationType) []dhlclient.DhlLocationType {
	var result []dhlclient.DhlLocationType

	for _, t := range types {
		if t == domain.LOCATION_TYPE_LOCKER {
			result = append(result, dhlclient.DHL_LOCATION_TYPE_LOCKER)
		}
		if t == domain.LOCATION_TYPE_SERVICE_POINT {
			result = append(result, dhlclient.DHL_LOCATION_TYPE_POSTOFFICE)
			result = append(result, dhlclient.DHL_LOCATION_TYPE_SERVICE_POINT)
		}
	}

	if len(result) == 0 {
		result = append(result, dhlclient.DHL_LOCATION_TYPE_ALL)
	}

	return result
}
