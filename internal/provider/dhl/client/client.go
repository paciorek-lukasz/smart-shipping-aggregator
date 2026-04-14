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
		name: "DHL",
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
		respBytes, err = json.Marshal(&DhlHomeApiResponse{
			Earliest: "2026-04-10T09:00:00Z",
			Latest:   "2026-04-12T18:00:00Z",
			Price:    1299,
			Currency: "PLN",
		})
	} else {
		respBytes, err = json.Marshal(&DhlPickupApiResponse{
			Earliest: "2026-04-09T10:00:00Z",
			Latest:   "2026-04-11T20:00:00Z",
			Price:    899,
			Currency: "PLN",
			Locations: []*Location{
				{
					Id:          "DHL-PO-001",
					Name:        "DHL Punkt Partnerski",
					City:        "Warszawa",
					PostalCode:  "00-001",
					Country:     "PL",
					Latitude:    "52.2297",
					Longitude:   "21.0122",
					AddressLine: "ul. Marszałkowska 1",
					Type:        "postoffice",
					OpenTimes: []*OpenTimes{
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
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
