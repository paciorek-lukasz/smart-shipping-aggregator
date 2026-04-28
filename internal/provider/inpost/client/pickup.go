package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

const (
	getQuotesPickupEndpointInpost = "/api/v1/quotes/pickup"
)

func (c *Client) GetQuotesPickup(ctx context.Context, req *InpostPickupApiRequest) (*InpostPickupApiResponse, error) {
	var resp []byte

	bb, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fullUrl, err := url.JoinPath(c.baseUrl, getQuotesPickupEndpointInpost)
	if err != nil {
		return nil, fmt.Errorf("failed to build request url: %w", err)
	}

	apiReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, bytes.NewReader(bb))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	err = c.cb.Execute(func() error {
		inpostRes, err := c.simulateApiCall(ctx, apiReq, domain.DELIVERY_TYPE_PICKUP)
		if err != nil {
			return err
		}

		resp = inpostRes
		return nil
	})

	var res InpostPickupApiResponse
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &res, err
}
