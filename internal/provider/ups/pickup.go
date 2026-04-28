package ups

import (
	"context"
	"errors"
	"fmt"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/ups/client"
)

func (s *Service) sendPickupRequest(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.UpsPickupApiRequest{
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
		DropPointsLimit: s.dropPointsLimit,
		SearchRadius:    s.searchRadius,
		DropPointType:   mapLocationTypes(req.LocationTypes),
	}

	apiCtx := context.WithValue(ctx, carrierName, "get_quotes_pickup")

	resp, err := s.apiClient.GetQuotesPickup(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.AvailableFrom, resp.AvailableTo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	locations, err := parseDropPoints(resp.DropPoints)
	if err != nil {
		return nil, fmt.Errorf("failed to parse drop points: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.FeeAmount,
			Currency:          resp.FeeCurrency,
			DeliveryTimeSlots: timeslots,
			PickupPoints:      locations,
			DeliveryType:      domain.DELIVERY_TYPE_PICKUP,
		},
	}, nil
}

func parseDropPoints(drops []*client.DropPoint) ([]*domain.PickupPoint, error) {
	var result []*domain.PickupPoint

	for _, drop := range drops {
		locType, err := mapUpsCategory(drop.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to map drop point type: %w", err)
		}

		result = append(result, &domain.PickupPoint{
			PickupPointId: drop.DropPointId,
			Name:          drop.DisplayName,
			Address: &domain.Address{
				Address:    drop.FullAddress,
				PostalCode: drop.ZipCode,
				City:       drop.City,
				Country:    drop.Country,
				Longitude:  drop.CoordLng,
				Latitude:   drop.CoordLat,
			},
			LocationType:  locType,
			OpeningHours:  parseSchedule(drop.OperationHours),
			IsOperational: drop.IsActive,
		})
	}

	return result, nil
}

func mapUpsCategory(category string) (domain.LocationType, error) {
	switch category {
	case client.UpsCategoryDropBox:
		return domain.LOCATION_TYPE_LOCKER, nil
	case client.UpsCategoryPackageShop:
		return domain.LOCATION_TYPE_SERVICE_POINT, nil
	}

	return 0, errors.New("no matching location type")
}

func parseSchedule(schedule []*client.Schedule) []*domain.OpeningHours {
	var result []*domain.OpeningHours

	for _, h := range schedule {
		result = append(result, &domain.OpeningHours{
			DayOfWeek: h.Weekday,
			Opens:     h.StartHour,
			Closes:    h.EndHour,
		})
	}

	return result
}

func mapLocationTypes(types []domain.LocationType) []client.UpsDropType {
	var result []client.UpsDropType

	for _, t := range types {
		if t == domain.LOCATION_TYPE_LOCKER {
			result = append(result, client.UPS_DROP_TYPE_DROP_BOX)
		}
		if t == domain.LOCATION_TYPE_SERVICE_POINT {
			result = append(result, client.UPS_DROP_TYPE_PACKAGE_SHOP)
			result = append(result, client.UPS_DROP_TYPE_AUTHORIZED)
		}
	}

	if len(result) == 0 {
		result = append(result, client.UPS_DROP_TYPE_ALL)
	}

	return result
}
