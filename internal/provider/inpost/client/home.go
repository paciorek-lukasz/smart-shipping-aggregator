package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
)

const (
	getQuotesHomeEndpointInpost = "/api/v1/quotes/home"
)

func (c *Client) GetQuotesHome(ctx context.Context, req *InpostHomeApiRequest) (*InpostHomeApiResponse, error) {
	var resp []byte

	bb, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fullUrl, err := url.JoinPath(c.baseUrl, getQuotesHomeEndpointInpost)
	if err != nil {
		return nil, fmt.Errorf("failed to build request url: %w", err)
	}

	apiReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, bytes.NewReader(bb))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	err = c.cb.Execute(func() error {
		inpostRes, err := c.simulateApiCall(ctx, apiReq, domain.DELIVERY_TYPE_HOME_DELIVERY)
		if err != nil {
			return err
		}

		resp = inpostRes
		return nil
	})

	var res InpostHomeApiResponse
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &res, err
}
