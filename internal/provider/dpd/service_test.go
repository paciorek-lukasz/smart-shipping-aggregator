package dpd

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dpd/client"
	mocks "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dpd/mocks"
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
				testutilities.GotResult(validGetQuotesPickupResponse(), testutilities.IgnoreUnexportedFields()),
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
			name: "success pickup, empty locations",
			req:  validGetQuotesRequest(domain.DELIVERY_TYPE_PICKUP, nil),
			prepareApiClient: func(recorder *mocks.MockapiClientMockRecorder) {
				recorder.GetQuotesPickup(gomock.Any(), gomock.Any()).Return(emptyLocationsResponse(), nil)
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

func validApiHomeResponse() *client.DpdHomeApiResponse {
	return &client.DpdHomeApiResponse{
		Earliest: "2026-04-10T09:00:00Z",
		Latest:   "2026-04-12T18:00:00Z",
		Price:    1000,
		Currency: "USD",
	}
}

func validApiPickupResponse() *client.DpdPickupApiResponse {
	return &client.DpdPickupApiResponse{
		Earliest:  "2026-04-09T10:00:00Z",
		Latest:    "2026-04-11T20:00:00Z",
		Price:     899,
		Currency:  "PLN",
		Locations: []*client.DpdLocation{
			{
				Id:         "DPD-L-001",
				Name:       "DPD Location",
				City:       "Warszawa",
				PostalCode: "00-001",
				Country:    "PL",
				Latitude:   "52.2297",
				Longitude:  "21.0122",
				AddressLine: "ul. Marszałkowska 1",
				Type:       "packageshop",
				OpenTimes: []*client.DpdOpenTimes{
					{DayOfWeek: "MONDAY", Opens: "08:00", Closes: "20:00"},
					{DayOfWeek: "TUESDAY", Opens: "08:00", Closes: "20:00"},
					{DayOfWeek: "WEDNESDAY", Opens: "08:00", Closes: "20:00"},
					{DayOfWeek: "THURSDAY", Opens: "08:00", Closes: "20:00"},
					{DayOfWeek: "FRIDAY", Opens: "08:00", Closes: "20:00"},
					{DayOfWeek: "SATURDAY", Opens: "09:00", Closes: "16:00"},
					{DayOfWeek: "SUNDAY", Opens: "", Closes: ""},
				},
				IsAvailable: true,
			},
		},
	}
}

func validGetQuotesHomeResponse() *domain.GetQuotesResponse {
	return &domain.GetQuotesResponse{
		Options: &domain.Option{
			CarrierProduct: "dpd",
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
			CarrierProduct: "dpd",
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
					PickupPointId: "DPD-L-001",
					Name:          "DPD Location",
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

func invalidTimeslotResponse() *client.DpdHomeApiResponse {
	return &client.DpdHomeApiResponse{
		Earliest: "invalid-date",
		Latest:   "2026-04-12T18:00:00Z",
		Price:    1000,
		Currency: "USD",
	}
}

func invalidPickupTimeslotResponse() *client.DpdPickupApiResponse {
	return &client.DpdPickupApiResponse{
		Earliest:  "bad-date",
		Latest:    "2026-04-11T20:00:00Z",
		Price:     899,
		Currency:  "PLN",
		Locations: []*client.DpdLocation{
			{
				Id:         "DPD-L-001",
				Name:       "DPD Location",
				City:       "Warszawa",
				PostalCode: "00-001",
				Country:    "PL",
				Latitude:   "52.2297",
				Longitude:  "21.0122",
				AddressLine: "ul. Marszałkowska 1",
				Type:       "packageshop",
				OpenTimes:  []*client.DpdOpenTimes{},
				IsAvailable: true,
			},
		},
	}
}

func emptyLocationsResponse() *client.DpdPickupApiResponse {
	return &client.DpdPickupApiResponse{
		Earliest:  "2026-04-09T10:00:00Z",
		Latest:    "2026-04-11T20:00:00Z",
		Price:     899,
		Currency:  "PLN",
		Locations: []*client.DpdLocation{},
	}
}

func unknownLocationTypeResponse() *client.DpdPickupApiResponse {
	return &client.DpdPickupApiResponse{
		Earliest:  "2026-04-09T10:00:00Z",
		Latest:    "2026-04-11T20:00:00Z",
		Price:     899,
		Currency:  "PLN",
		Locations: []*client.DpdLocation{
			{
				Id:         "DPD-L-001",
				Name:       "DPD Location",
				City:       "Warszawa",
				PostalCode: "00-001",
				Country:    "PL",
				Latitude:   "52.2297",
				Longitude:  "21.0122",
				AddressLine: "ul. Marszałkowska 1",
				Type:      "unknown_type",
				OpenTimes:  []*client.DpdOpenTimes{},
				IsAvailable: true,
			},
		},
	}
}