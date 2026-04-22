package aggregator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider"
)

type Service struct {
	providers []provider.Provider
	timeout   time.Duration
}

func NewService(providers []provider.Provider, timeout time.Duration) *Service {
	return &Service{
		providers: providers,
		timeout:   timeout,
	}
}

func (s *Service) FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) *domain.GetOptionsResponse {
	resultsChan := make(chan *domain.GetQuotesResponse, len(s.providers))

	fmt.Println(s.timeout)
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	for i, p := range s.providers {
		go func(p provider.Provider) {
			option, err := p.GetQuotes(ctx, req)
			if err != nil {
				log.Printf("provider %s responded with error: %s", p.Name(), err)
				return
			}

			option.Options.OptionId = int32(i)

			resultsChan <- option
		}(p)
	}

	var finalOptions []*domain.Option

	for i := 0; i < len(s.providers); i++ {
		select {
		case opt := <-resultsChan:
			finalOptions = append(finalOptions, opt.Options)
		case <-ctx.Done():
			log.Println("aggergator timeout reached, returning partial results")
			return &domain.GetOptionsResponse{
				Options: finalOptions,
			}
		}
	}

	return &domain.GetOptionsResponse{
		Options: finalOptions,
	}
}
