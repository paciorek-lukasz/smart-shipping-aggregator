package aggregator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/domain"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider"
)

var validate = validator.New()

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

func (s *Service) FetchQuotes(ctx context.Context, req *domain.GetQuotesRequest) (*domain.GetOptionsResponse, map[string]string) {
	errs := validateRequest(req)
	if errs != nil {
		return nil, errs
	}

	resultsChan := make(chan *domain.GetQuotesResponse, len(s.providers))

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
			}, nil
		}
	}

	return &domain.GetOptionsResponse{
		Options: finalOptions,
	}, nil
}

func validateRequest(req *domain.GetQuotesRequest) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	errorReport := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			var message string

			switch ve.Tag() {
			case "required":
				message = "field cannot be empty"
			case "min":
				message = fmt.Sprintf("field is too short, min is %s characters", ve.Param())
			case "max":
				message = fmt.Sprintf("field is too long, max is %s characters", ve.Param())
			case "email":
				message = "invalid email address"
			case "e164":
				message = "invalid phone number"
			default:
				message = fmt.Sprintf("failed validation rule: %s", ve.Tag())
			}

			errorReport[ve.Field()] = message
		}
	}

	return errorReport
}
