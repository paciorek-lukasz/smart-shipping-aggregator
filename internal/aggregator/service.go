package aggregator

import (
	"context"
	"log"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
)

type QuotesProvider interface {
	GetQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.Option, error)
	Name() string
}

type Service struct {
	providers []*QuotesProvider
	timeout   time.Duration
}

func NewService(providers []*QuotesProvider, timeout time.Duration) *Service {
	return &Service{
		providers: providers,
		timeout:   timeout,
	}
}

func (s *Service) FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) *domain.GetQuotesResponse {
	resultsChan := make(chan *domain.Option, len(s.providers))

	ctx, cancel := context.WithTimeout(ctx, s.timeout*time.Second)
	defer cancel()

	for _, provider := range s.providers {
		go func(p QuotesProvider) {
			option, err := p.GetQuotes(ctx, req)
			if err != nil {
				log.Printf("provider %s responded with error: %s", p.Name(), err)
				return
			}

			resultsChan <- option
		}(*provider)
	}

	var finalOptions []*domain.Option

	for i := 0; i < len(s.providers); i++ {
		select {
		case opt := <-resultsChan:
			finalOptions = append(finalOptions, opt)
		case <-ctx.Done():
			log.Println("aggergator timeout reached, returning partial results")
			return &domain.GetQuotesResponse{
				Options: finalOptions,
			}
		}
	}

	return &domain.GetQuotesResponse{
		Options: finalOptions,
	}
}
