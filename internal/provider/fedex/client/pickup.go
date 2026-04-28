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
	getQuotesPickupEndpointFedex = "/rate/v1/quotes/pickup"
)

func (c *Client) GetQuotesPickup(ctx context.Context, req *FedexPickupApiRequest) (*FedexPickupApiResponse, error) {
	var resp []byte

	bb, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fullUrl, err := url.JoinPath(c.baseUrl, getQuotesPickupEndpointFedex)
	if err != nil {
		return nil, fmt.Errorf("failed to build request url: %w", err)
	}

	apiReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, bytes.NewReader(bb))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	err = c.cb.Execute(func() error {
		fedexRes, err := c.simulateApiCall(ctx, apiReq, domain.DELIVERY_TYPE_PICKUP)
		if err != nil {
			return err
		}

		resp = fedexRes
		return nil
	})

	var res FedexPickupApiResponse
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &res, err
}
