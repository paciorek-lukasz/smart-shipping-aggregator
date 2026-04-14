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

func New(cb *resilience.CircuitBreaker) *Client {
	return &Client{
		name: "Fedex",
		cb:   cb,
	}
}

func (c *Client) Name() string { return c.name }

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
		respBytes, err = json.Marshal(&FedexHomeApiResponse{
			EstimatedPickup:   "2026-04-10T14:00:00Z",
			EstimatedDelivery: "2026-04-13T17:00:00Z",
			NetAmount:         1999,
			CurrencyCode:      "PLN",
		})
	} else {
		respBytes, err = json.Marshal(&FedexPickupApiResponse{
			EstimatedPickup:   "2026-04-09T11:00:00Z",
			EstimatedDelivery: "2026-04-11T19:00:00Z",
			NetAmount:         1299,
			CurrencyCode:      "PLN",
			Droppoints: []*FedexDropPoint{
				{
					LocationId:      "FX-LOC-001",
					LocationName:    "FedEx Station Poznań",
					City:            "Poznań",
					StateOrProvince: "WP",
					PostalCode:      "61-001",
					CountryCode:     "PL",
					LatLong:         "52.4064,16.9252",
					Address:         "ul. Wieniawskiego 5",
					LocationType:    "station",
					HoursOperation: []*OperationWindow{
						{DaysOfOperation: "Mon", WindowOpen: "07:00", WindowClose: "21:00"},
						{DaysOfOperation: "Tue", WindowOpen: "07:00", WindowClose: "21:00"},
						{DaysOfOperation: "Wed", WindowOpen: "07:00", WindowClose: "21:00"},
						{DaysOfOperation: "Thu", WindowOpen: "07:00", WindowClose: "21:00"},
						{DaysOfOperation: "Fri", WindowOpen: "07:00", WindowClose: "21:00"},
						{DaysOfOperation: "Sat", WindowOpen: "08:00", WindowClose: "16:00"},
						{DaysOfOperation: "Sun", WindowOpen: "", WindowClose: ""},
					},
					CurrentlyOpen: true,
				},
			},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
