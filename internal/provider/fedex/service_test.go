package fedex

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/fedex/client"
	mocks "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/fedex/mocks"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/testutilities"
)

func TestService_GetQuotes(t *testing.T) {
	testCases := []struct {
		name             string
		req              *domain.GetQuotesRequest
		prepareApiClient func(recorder *mocks.MockapiClientMockRecorder)
		checks           []testutilities.Check
	}{
		{
			name: "success home delivery",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_HOME_DELIVERY, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesHome(gomock.Any(), gomock.Any()).Return(validApiHomeResponse(), nil)
			},
			checks: []testutilities.Check{
				testutilities.GotSuccess,
				testutilities.GotResult(validGetQuotesHomeResponse(), testutilities.IgnoreUnexportedFields()),
			},
		},
		{
			name: "success pickup delivery",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(validApiPickupResponse(), nil)
			},
			checks: []testutilities.Check{
				testutilities.GotSuccess,
			},
		},
		{
			name: "fail, api returned error",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(nil, errors.New("some api error"))
			},
			checks: []testutilities.Check{
				testutilities.GotExactErrorMessage("failed to get quotes from carrier API: some api error"),
			},
		},
		{
			name: "fail home delivery, api returned error",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_HOME_DELIVERY, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesHome(gomock.Any(), gomock.Any()).Return(nil, errors.New("connection refused"))
			},
			checks: []testutilities.Check{
				testutilities.GotExactErrorMessage("failed to get quotes from carrier API: connection refused"),
			},
		},
		{
			name: "fail pickup, api returned error",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(nil, errors.New("timeout"))
			},
			checks: []testutilities.Check{
				testutilities.GotExactErrorMessage("failed to get quotes from carrier API: timeout"),
			},
		},
		{
			name: "fail home delivery, invalid timeslot",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_HOME_DELIVERY, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesHome(gomock.Any(), gomock.Any()).Return(invalidTimeslotResponse(), nil)
			},
			checks: []testutilities.Check{
				testutilities.GotContainsErrorMessage("failed to parse start timeslot"),
			},
		},
		{
			name: "fail pickup, invalid timeslot",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(invalidPickupTimeslotResponse(), nil)
			},
			checks: []testutilities.Check{
				testutilities.GotContainsErrorMessage("failed to parse start timeslot"),
			},
		},
		{
			name: "success pickup, empty droppoints",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(emptyDroppointsResponse(), nil)
			},
			checks: []testutilities.Check{
				testutilities.GotSuccess,
			},
		},
		{
			name: "fail pickup, unknown location type",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(unknownLocationTypeResponse(), nil)
			},
			checks: []testutilities.Check{
				testutilities.GotContainsErrorMessage("failed to map location type"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			apiCli := mocks.NewMockapiClient(ctrl)
			if tc.prepareApiClient != nil {
				tc.prepareApiClient(apiCli.EXPECT())
			}

			svc := NewService(apiCli, "apiKey", 10, 10)

			resp, err := svc.GetQuotes(context.Background(), tc.req)
			for _, ch := range tc.checks {
				ch(resp, err, t)
			}
		})
	}
}

func validGetQuotesRequest(deliveryType domain.DeliveryType, locationTypes []domain.LocationType) *domain.GetQuotesRequest {
	return &domain.GetQuotesRequest{
		Sender: &domain.Party{
			Name:  "Something Testy",
			Email: "warehouse@somethingtesty.com",
			Phone: "1234567890",
			Address: &domain.Address{
				Address:    "123 Main St",
				PostalCode: "12345",
				City:       "New York",
				Country:    "US",
				Longitude:  "31.1234",
				Latitude:   "33.1234",
			},
		},
		Recipient: &domain.Party{
			Name:  "James Kirk",
			Email: "james.kirk@federation.com",
			Phone: "1234567890",
			Address: &domain.Address{
				Address:    "456 Elm St",
				PostalCode: "67890",
				City:       "Los Angeles",
				Country:    "US",
				Longitude:  "31.1234",
				Latitude:   "33.1234",
			},
		},
		Package: &domain.Package{
			Items: []*domain.Item{
				{
					ItemID:   12345,
					Sku:      "WPI-1234",
					Name:     "Warp Plasma Intercooler",
					Quantity: 1,
					Price:    1299,
				},
			},
		},
		DeliveryType:  deliveryType,
		LocationTypes: locationTypes,
	}
}

func validApiHomeResponse() *client.FedexHomeApiResponse {
	return &client.FedexHomeApiResponse{
		EstimatedPickup:   "2026-04-10T09:00:00Z",
		EstimatedDelivery: "2026-04-12T18:00:00Z",
		NetAmount:       1000,
		CurrencyCode:    "USD",
	}
}

func validApiPickupResponse() *client.FedexPickupApiResponse {
	return &client.FedexPickupApiResponse{
		EstimatedPickup:   "2026-04-09T10:00:00Z",
		EstimatedDelivery: "2026-04-11T20:00:00Z",
		NetAmount:       899,
		CurrencyCode:     "PLN",
		Droppoints: []*client.FedexDropPoint{
			{
				LocationId:      "FEDEX-L-001",
				LocationName:  "FedEx Location",
				City:        "Warszawa",
				PostalCode:  "00-001",
				CountryCode: "PL",
				LatLong:    "52.2297,21.0122",
				Address:    "ul. Marszałkowska 1",
				LocationType: "station",
				HoursOperation: []*client.OperationWindow{
					{DaysOfOperation: "MONDAY", WindowOpen: "08:00", WindowClose: "20:00"},
					{DaysOfOperation: "TUESDAY", WindowOpen: "08:00", WindowClose: "20:00"},
					{DaysOfOperation: "WEDNESDAY", WindowOpen: "08:00", WindowClose: "20:00"},
					{DaysOfOperation: "THURSDAY", WindowOpen: "08:00", WindowClose: "20:00"},
					{DaysOfOperation: "FRIDAY", WindowOpen: "08:00", WindowClose: "20:00"},
					{DaysOfOperation: "SATURDAY", WindowOpen: "09:00", WindowClose: "16:00"},
					{DaysOfOperation: "SUNDAY", WindowOpen: "", WindowClose: ""},
				},
				CurrentlyOpen: true,
			},
		},
	}
}

func validGetQuotesHomeResponse() *domain.GetQuotesResponse {
	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct: "fedex",
			Price:          1000,
			Currency:       "USD",
			DeliveryTimeSlots: []*domain.DeliveryTimeSlot{
				{
					Start: time.Date(2026, 4, 10, 9, 0, 0, 0, time.UTC),
					End:   time.Date(2026, 4, 12, 18, 0, 0, 0, time.UTC),
				},
			},
			DeliveryType: domain.DELIVERY_TYPE_HOME_DELIVERY,
		},
	}
}

func validGetQuotesPickupResponse() *domain.GetQuotesResponse {
	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct: "fedex",
			Price:          899,
			Currency:       "PLN",
			DeliveryTimeSlots: []*domain.DeliveryTimeSlot{
				{
					Start: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
					End:   time.Date(2026, 4, 11, 20, 0, 0, 0, time.UTC),
				},
			},
			DeliveryType: domain.DELIVERY_TYPE_PICKUP,
			PickupPoints: []*domain.PickupPoint{
				{
					PickupPointId: "FEDEX-L-001",
					Name:          "FedEx Location",
					Address: &domain.Address{
						Address:    "ul. Marszałkowska 1",
						PostalCode: "00-001",
						City:       "Warszawa",
						Country:    "PL",
						Latitude:   "52.2297",
						Longitude:  "21.0122",
					},
					LocationType:  domain.LOCATION_TYPE_SERVICE_POINT,
					IsOperational: true,
					OpeningHours: []*domain.OpeningHours{
						{
							DayOfWeek: "MONDAY",
							Opens:     "08:00",
							Closes:    "20:00",
						},
						{
							DayOfWeek: "TUESDAY",
							Opens:     "08:00",
							Closes:    "20:00",
						},
						{
							DayOfWeek: "WEDNESDAY",
							Opens:     "08:00",
							Closes:    "20:00",
						},
						{
							DayOfWeek: "THURSDAY",
							Opens:     "08:00",
							Closes:    "20:00",
						},
						{
							DayOfWeek: "FRIDAY",
							Opens:     "08:00",
							Closes:    "20:00",
						},
						{
							DayOfWeek: "SATURDAY",
							Opens:     "09:00",
							Closes:    "16:00",
						},
						{
							DayOfWeek: "SUNDAY",
							Opens:     "",
							Closes:    "",
						},
					},
				},
			},
		},
	}
}

func invalidTimeslotResponse() *client.FedexHomeApiResponse {
	return &client.FedexHomeApiResponse{
		EstimatedPickup:   "invalid-date",
		EstimatedDelivery: "2026-04-12T18:00:00Z",
		NetAmount:       1000,
		CurrencyCode:    "USD",
	}
}

func invalidPickupTimeslotResponse() *client.FedexPickupApiResponse {
	return &client.FedexPickupApiResponse{
		EstimatedPickup:   "bad-date",
		EstimatedDelivery: "2026-04-11T20:00:00Z",
		NetAmount:       899,
		CurrencyCode:     "PLN",
		Droppoints: []*client.FedexDropPoint{
			{
				LocationId:      "FEDEX-L-001",
				LocationName:  "FedEx Location",
				City:        "Warszawa",
				PostalCode:  "00-001",
				CountryCode: "PL",
				LatLong:    "52.2297,21.0122",
				Address:    "ul. Marszałkowska 1",
				LocationType: "station",
				HoursOperation: []*client.OperationWindow{},
				CurrentlyOpen: true,
			},
		},
	}
}

func emptyDroppointsResponse() *client.FedexPickupApiResponse {
	return &client.FedexPickupApiResponse{
		EstimatedPickup:   "2026-04-09T10:00:00Z",
		EstimatedDelivery: "2026-04-11T20:00:00Z",
		NetAmount:       899,
		CurrencyCode:     "PLN",
		Droppoints:     []*client.FedexDropPoint{},
	}
}

func unknownLocationTypeResponse() *client.FedexPickupApiResponse {
	return &client.FedexPickupApiResponse{
		EstimatedPickup:   "2026-04-09T10:00:00Z",
		EstimatedDelivery: "2026-04-11T20:00:00Z",
		NetAmount:       899,
		CurrencyCode:     "PLN",
		Droppoints: []*client.FedexDropPoint{
			{
				LocationId:      "FEDEX-L-001",
				LocationName:  "FedEx Location",
				City:        "Warszawa",
				PostalCode:  "00-001",
				CountryCode: "PL",
				LatLong:    "52.2297,21.0122",
				Address:    "ul. Marszałkowska 1",
				LocationType: "unknown_type",
				HoursOperation: []*client.OperationWindow{},
				CurrentlyOpen: true,
			},
		},
	}
}