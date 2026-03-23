package dhl

import (
	"context"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/resilience"
)

type Client struct {
	name string
	cb   *resilience.CircuitBreaker
	// Tu byłby np. http.Client
}

func New(cb *resilience.CircuitBreaker) *Client {
	return &Client{
		name: "DHL",
		cb:   cb,
	}
}

func (c *Client) Name() string { return c.name }

func (c *Client) GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error) {
	var opt domain.Option

	err := c.cb.Execute(func() error {
		dhlReq := toDhlRequest(req)

		dhlRes, err := c.simulateApiCall(ctx, dhlReq)
		if err != nil {
			return err
		}

		opt = toDomainOption(dhlRes)
		return nil
	})

	return opt, err
}
