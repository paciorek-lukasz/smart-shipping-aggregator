package dhlclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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

func (c *Client) simulateApiCall(ctx context.Context, req *http.Request) ([]byte, error) {
	delay := c.minDelay + time.Duration(rand.Intn(int(c.maxDelay-c.minDelay)))

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return []byte{}, ctx.Err()
	}

	if rand.Float64() < c.failureRate {
		return []byte{}, errors.New("internal server error (500)")
	}

	respBytes, err := json.Marshal(&DhlHomeApiResponse{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse response mock: %w", err)
	}

	return respBytes, nil
}
