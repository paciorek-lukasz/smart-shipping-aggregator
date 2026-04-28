package fedex

import (
	"context"
	"errors"
	"fmt"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/fedex/client"
)

func (s *Service) sendPickupRequest(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.FedexPickupApiRequest{
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
		DroppointsLimit:   s.droppointsLimit,
		WithinRadiusMiles: s.withinRadiusMiles,
		ServiceTypeFilter: mapLocationTypes(req.LocationTypes),
	}

	apiCtx := context.WithValue(ctx, carrierName, "get_quotes_pickup")

	resp, err := s.apiClient.GetQuotesPickup(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.EstimatedPickup, resp.EstimatedDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	droppoints, err := parseDropPoints(resp.Droppoints)
	if err != nil {
		return nil, fmt.Errorf("failed to parse drop points: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.NetAmount,
			Currency:          resp.CurrencyCode,
			DeliveryTimeSlots: timeslots,
			PickupPoints:      droppoints,
			DeliveryType:      domain.DELIVERY_TYPE_PICKUP,
		},
	}, nil
}

func parseDropPoints(drops []*client.FedexDropPoint) ([]*domain.PickupPoint, error) {
	var result []*domain.PickupPoint

	for _, d := range drops {
		locType, err := mapFedexLocationType(d.LocationType)
		if err != nil {
			return nil, fmt.Errorf("failed to map location type: %w", err)
		}

		result = append(result, &domain.PickupPoint{
			PickupPointId: d.LocationId,
			Name:          d.LocationName,
			Address: &domain.Address{
				Address:    d.Address,
				PostalCode: d.PostalCode,
				City:       d.City,
				Country:    d.CountryCode,
				Longitude:  d.LatLong,
				Latitude:   d.LatLong,
			},
			LocationType:  locType,
			OpeningHours:  parseOperationWindows(d.HoursOperation),
			IsOperational: d.CurrentlyOpen,
		})
	}

	return result, nil
}

func mapFedexLocationType(locType string) (domain.LocationType, error) {
	switch locType {
	case client.FedexLocTypeStation:
		return domain.LOCATION_TYPE_LOCKER, nil
	case client.FedexLocTypeauthorized:
		return domain.LOCATION_TYPE_SERVICE_POINT, nil
	}

	return 0, errors.New("no matching location type")
}

func parseOperationWindows(windows []*client.OperationWindow) []*domain.OpeningHours {
	var result []*domain.OpeningHours

	for _, w := range windows {
		result = append(result, &domain.OpeningHours{
			DayOfWeek: w.DaysOfOperation,
			Opens:     w.WindowOpen,
			Closes:    w.WindowClose,
		})
	}

	return result
}

func mapLocationTypes(types []domain.LocationType) []client.FedexServiceType {
	var result []client.FedexServiceType

	for _, t := range types {
		if t == domain.LOCATION_TYPE_LOCKER {
			result = append(result, client.FEDEX_SERVICE_GROUND)
		}
		if t == domain.LOCATION_TYPE_SERVICE_POINT {
			result = append(result, client.FEDEX_SERVICE_EXPRESS)
			result = append(result, client.FEDEX_SERVICE_SMART_POST)
		}
	}

	if len(result) == 0 {
		result = append(result, client.FEDEX_SERVICE_ALL)
	}

	return result
}
