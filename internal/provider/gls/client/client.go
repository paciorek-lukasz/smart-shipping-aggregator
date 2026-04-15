package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/resilience"
)

type Client struct {
	baseUrl     string
	name        string
	cb          *resilience.CircuitBreaker
	failureRate float64
	minDelay    time.Duration
	maxDelay    time.Duration
}

func New(cb *resilience.CircuitBreaker, baseURL string, failureRate float64, minDelay, maxDelay time.Duration) *Client {
	return &Client{
		name:        "GLS",
		cb:          cb,
		baseUrl:     baseURL,
		failureRate: failureRate,
		minDelay:    minDelay,
		maxDelay:    maxDelay,
	}
}

func (c *Client) simulateApiCall(ctx context.Context, req *http.Request, deliveryType domain.DeliveryType) ([]byte, error) {
	delay := c.minDelay + time.Duration(rand.Intn(int(c.maxDelay-c.minDelay)))

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return []byte{}, ctx.Err()
	}

	if rand.Float64() < c.failureRate {
		return []byte{}, errors.New("internal server error (500)")
	}

	var respBytes []byte
	var err error
	if deliveryType == domain.DELIVERY_TYPE_HOME_DELIVERY {
		respBytes, err = json.Marshal(&GlsHomeApiResponse{
			WindowFrom: "2026-04-12T09:00:00Z",
			WindowTo:   "2026-04-14T17:00:00Z",
			RateCents:  1399,
			Currency:   "PLN",
		})
	} else {
		respBytes, err = json.Marshal(&GlsPickupApiResponse{
			WindowFrom: "2026-04-11T10:00:00Z",
			WindowTo:   "2026-04-13T18:00:00Z",
			RateCents:  899,
			Currency:   "PLN",
			Depots: []*DepotData{
				{
					DepotCode:  "GLS-PO-001",
					DepotName:  "GLS Parcel Shop Wrocław",
					City:       "Wrocław",
					PostalCode: "50-001",
					IsoCountry: "PL",
					GeoLat:     "51.1079",
					GeoLon:     "17.0385",
					Street:     "ul. Krupnicza 3",
					Facility:   "parcelshop",
					Timetable: []*TimetableDay{
						{DayName: "MON", OpenHr: "08:00", CloseHr: "20:00"},
						{DayName: "TUE", OpenHr: "08:00", CloseHr: "20:00"},
						{DayName: "WED", OpenHr: "08:00", CloseHr: "20:00"},
						{DayName: "THU", OpenHr: "08:00", CloseHr: "20:00"},
						{DayName: "FRI", OpenHr: "08:00", CloseHr: "20:00"},
						{DayName: "SAT", OpenHr: "09:00", CloseHr: "14:00"},
						{DayName: "SUN", OpenHr: "", CloseHr: ""},
					},
					Operational: true,
				},
			},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
