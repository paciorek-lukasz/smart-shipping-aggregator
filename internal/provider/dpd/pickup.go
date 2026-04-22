package dpd

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dpd/client"
)

func (s *Service) sendPickupRequest(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.DpdPickupApiRequest{
		SenderAddress: &client.DpdParty{
			Address:    senderAddress.Address,
			PostalCode: senderAddress.PostalCode,
			City:       senderAddress.City,
			Country:    senderAddress.Country,
		},
		RecipientAddress: &client.DpdParty{
			Address:    recipientAddress.Address,
			PostalCode: recipientAddress.PostalCode,
			City:       recipientAddress.City,
			Country:    recipientAddress.Country,
		},
		LocationsLimit: 10,
		SearchRadius:   5000,
	}

	apiCtx := context.WithValue(ctx, carrierName, "get_quotes_pickup")

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
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.Price,
			Currency:          resp.Currency,
			DeliveryTimeSlots: timeslots,
			PickupPoints:      locations,
			DeliveryType:      domain.DELIVERY_TYPE_PICKUP,
		},
	}, nil
}

func parsePickupPoints(locs []*client.DpdLocation) ([]*domain.PickupPoint, error) {
	var result []*domain.PickupPoint

	for _, loc := range locs {
		locType, err := mapDpdLocType(loc.Type)
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

func mapDpdLocType(dpdLocType string) (domain.LocationType, error) {
	switch dpdLocType {
	case client.DpdLocTypeLocker:
		return domain.LOCATION_TYPE_LOCKER, nil
	case client.DpdLocTypePackageShop:
		return domain.LOCATION_TYPE_SERVICE_POINT, nil
	}

	return 0, errors.New("no matching location type")
}

func parseOpeningHours(openingHours []*client.DpdOpenTimes) []*domain.OpeningHours {
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

func mapLocationTypes(types []domain.LocationType) []client.DpdLocationType {
	var result []client.DpdLocationType

	for _, t := range types {
		if t == domain.LOCATION_TYPE_LOCKER {
			result = append(result, client.DPD_LOCATION_TYPE_LOCKER)
		}
		if t == domain.LOCATION_TYPE_SERVICE_POINT {
			result = append(result, client.DPD_LOCATION_TYPE_PACKAGE_SHOP)
			result = append(result, client.DPD_LOCATION_TYPE_PARCEL_BOX)
		}
	}

	if len(result) == 0 {
		result = append(result, client.DPD_LOCATION_TYPE_ALL)
	}

	return result
}
