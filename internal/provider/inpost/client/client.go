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
		name:        "InPost",
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
		respBytes, err = json.Marshal(&InpostHomeApiResponse{
			DeliveryStart: "2026-04-09T10:00:00Z",
			DeliveryEnd:   "2026-04-11T18:00:00Z",
			ShipmentCost:  1199,
			IsoCurrency:   "PLN",
		})
	} else {
		respBytes, err = json.Marshal(&InpostPickupApiResponse{
			DeliveryStart: "2026-04-08T09:00:00Z",
			DeliveryEnd:   "2026-04-10T20:00:00Z",
			ShipmentCost:  799,
			IsoCurrency:   "PLN",
			Machines: []*ParcelLocker{
				{
					MachineId:  "KRA01-A",
					Name:       "Kraków Nowa Huta",
					Locality:   "Kraków",
					ZipCode:    "31-860",
					CountryIso: "PL",
					Latitude:   "50.0917",
					Longitude:  "20.0506",
					Address:    "al. Andersa 5",
					Status:     "active",
					Hours: []*Opening{
						{DayIndex: "1", OpenTime: "00:00", CloseTime: "23:59"},
						{DayIndex: "2", OpenTime: "00:00", CloseTime: "23:59"},
						{DayIndex: "3", OpenTime: "00:00", CloseTime: "23:59"},
						{DayIndex: "4", OpenTime: "00:00", CloseTime: "23:59"},
						{DayIndex: "5", OpenTime: "00:00", CloseTime: "23:59"},
						{DayIndex: "6", OpenTime: "00:00", CloseTime: "23:59"},
						{DayIndex: "7", OpenTime: "00:00", CloseTime: "23:59"},
					},
					Accepts: []string{"pack"},
				},
			},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
