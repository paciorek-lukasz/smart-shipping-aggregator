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
		name:        "UPS",
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
		respBytes, err = json.Marshal(&UpsHomeApiResponse{
			AvailableFrom: "2026-04-11T08:00:00Z",
			AvailableTo:   "2026-04-13T17:00:00Z",
			FeeAmount:     1499,
			FeeCurrency:   "PLN",
		})
	} else {
		respBytes, err = json.Marshal(&UpsPickupApiResponse{
			AvailableFrom: "2026-04-10T09:00:00Z",
			AvailableTo:   "2026-04-12T19:00:00Z",
			FeeAmount:     999,
			FeeCurrency:   "PLN",
			DropPoints: []*DropPoint{
				{
					DropPointId: "UPS-SHOP-001",
					DisplayName: "UPS Package Shop",
					City:        "Kraków",
					ZipCode:     "31-001",
					Country:     "PL",
					CoordLat:    "50.0647",
					CoordLng:    "19.9450",
					FullAddress: "ul. Floriańska 12",
					Category:    "package-shop",
					OperationHours: []*Schedule{
						{Weekday: "MONDAY", StartHour: "07:00", EndHour: "21:00"},
						{Weekday: "TUESDAY", StartHour: "07:00", EndHour: "21:00"},
						{Weekday: "WEDNESDAY", StartHour: "07:00", EndHour: "21:00"},
						{Weekday: "THURSDAY", StartHour: "07:00", EndHour: "21:00"},
						{Weekday: "FRIDAY", StartHour: "07:00", EndHour: "21:00"},
						{Weekday: "SATURDAY", StartHour: "08:00", EndHour: "15:00"},
						{Weekday: "SUNDAY", StartHour: "", EndHour: ""},
					},
					IsActive: true,
				},
			},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
