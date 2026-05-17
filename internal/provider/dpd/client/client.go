package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/resilience"
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
		name:        "DPD",
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
		respBytes, err = json.Marshal(&DpdHomeApiResponse{
			Earliest: "2026-04-11T08:00:00Z",
			Latest:   "2026-04-13T17:00:00Z",
			Price:    1099,
			Currency: "PLN",
		})
	} else {
		respBytes, err = json.Marshal(&DpdPickupApiResponse{
			Earliest: "2026-04-10T09:00:00Z",
			Latest:   "2026-04-12T19:00:00Z",
			Price:    799,
			Currency: "PLN",
			Locations: []*DpdLocation{
				{
					Id:          "DPD-PS-001",
					Name:        "DPD Parcel Shop",
					City:        "Kraków",
					PostalCode:  "30-001",
					Country:     "PL",
					Latitude:    "50.0647",
					Longitude:   "19.9450",
					AddressLine: "ul. Floriańska 10",
					Type:        "packageshop",
					OpenTimes: []*DpdOpenTimes{
						{DayOfWeek: "MONDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "TUESDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "WEDNESDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "THURSDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "FRIDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "SATURDAY", Opens: "08:00", Closes: "18:00"},
						{DayOfWeek: "SUNDAY", Opens: "", Closes: ""},
					},
					IsAvailable: true,
				},
				{
					Id:          "DPD-PS-002",
					Name:        "DPD Parcel Shop",
					City:        "Kraków",
					PostalCode:  "30-002",
					Country:     "PL",
					Latitude:    "50.0647",
					Longitude:   "19.9450",
					AddressLine: "ul. Floriańska 12",
					Type:        "packageshop",
					OpenTimes: []*DpdOpenTimes{
						{DayOfWeek: "MONDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "TUESDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "WEDNESDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "THURSDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "FRIDAY", Opens: "07:00", Closes: "21:00"},
						{DayOfWeek: "SATURDAY", Opens: "08:00", Closes: "18:00"},
						{DayOfWeek: "SUNDAY", Opens: "", Closes: ""},
					},
					IsAvailable: true,
				},
			},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
