package gls

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/gls/client"
)

func (s *Service) sendPickupRequest(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	senderAddress := req.Sender.Address
	recipientAddress := req.Recipient.Address

	apiReq := &client.GlsPickupApiRequest{
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
		DepotsLimit:       s.depotsLimit,
		SearchAreaKm:      s.searchAreaKm,
		DepotFacilityType: mapLocationTypes(req.LocationTypes),
	}

	apiCtx := context.WithValue(ctx, nil, "get_quotes_pickup")

	resp, err := s.apiClient.GetQuotesPickup(apiCtx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes from carrier API: %w", err)
	}

	timeslots, err := parseTimeslots(resp.WindowFrom, resp.WindowTo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeslots: %w", err)
	}

	depots, err := parseDepots(resp.Depots)
	if err != nil {
		return nil, fmt.Errorf("failed to parse depots: %w", err)
	}

	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct:    carrierName,
			Price:             resp.RateCents,
			Currency:          resp.Currency,
			DeliveryTimeSlots: timeslots,
			PickupPoints:      depots,
			DeliveryType:      domain.DELIVERY_TYPE_PICKUP,
		},
	}, nil
}

func parseDepots(depots []*client.DepotData) ([]*domain.PickupPoint, error) {
	var result []*domain.PickupPoint

	for _, d := range depots {
		locType, err := mapGlsFacility(d.Facility)
		if err != nil {
			return nil, fmt.Errorf("failed to map depot facility: %w", err)
		}

		result = append(result, &domain.PickupPoint{
			PickupPointId: d.DepotCode,
			Name:          d.DepotName,
			Address: &domain.Address{
				Address:    d.Street,
				PostalCode: d.PostalCode,
				City:       d.City,
				Country:    d.IsoCountry,
				Longitude:  d.GeoLon,
				Latitude:   d.GeoLat,
			},
			LocationType:  locType,
			OpeningHours:  parseTimetable(d.Timetable),
			IsOperational: d.Operational,
		})
	}

	return result, nil
}

func mapGlsFacility(facility string) (domain.LocationType, error) {
	switch facility {
	case client.GlsFacilityShop:
		return domain.LOCATION_TYPE_SERVICE_POINT, nil
	case client.GlsFacilityDepot:
		return domain.LOCATION_TYPE_LOCKER, nil
	}

	return 0, errors.New("no matching facility type")
}

func parseTimetable(timetable []*client.TimetableDay) []*domain.OpeningHours {
	var result []*domain.OpeningHours

	for _, t := range timetable {
		result = append(result, &domain.OpeningHours{
			DayOfWeek: t.DayName,
			Opens:     t.OpenHr,
			Closes:    t.CloseHr,
		})
	}

	return result
}

func mapLocationTypes(types []domain.LocationType) []client.GlsFacilityType {
	var result []client.GlsFacilityType

	for _, t := range types {
		if t == domain.LOCATION_TYPE_LOCKER {
			result = append(result, client.GLS_FACILITY_DEPOT)
		}
		if t == domain.LOCATION_TYPE_SERVICE_POINT {
			result = append(result, client.GLS_FACILITY_PARCELSHOP)
			result = append(result, client.GLS_FACILITY_PARTNER)
		}
	}

	if len(result) == 0 {
		result = append(result, client.GLS_FACILITY_ALL)
	}

	return result
}
