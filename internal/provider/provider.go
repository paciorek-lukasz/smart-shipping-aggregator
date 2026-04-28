package provider

import (
	"context"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
)

type Provider interface {
	GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetQuotesResponse, error)
	Name() string
}
