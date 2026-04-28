package provider

import (
	"time"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/config"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/dhl"
	dhlclient "github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/dhl/client"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/dpd"
	dpdclient "github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/dpd/client"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/fedex"
	fedexclient "github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/fedex/client"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/gls"
	glsclient "github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/gls/client"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/inpost"
	inpostclient "github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/inpost/client"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/ups"
	upsclient "github.com/paciorek-lukasz/smart-shipping-aggregator/internal/provider/ups/client"
	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/resilience"
)

func InitProviders(cfg config.EnvConfig) []Provider {
	var providers []Provider

	cb := resilience.NewCircuitBreaker(cfg.CircuitBreakerTreshold, cfg.CircuitBreakerCooldown)

	if cfg.EnableDhl {
		dhlClient := dhlclient.New(cb, cfg.DhlBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		dhlService := dhl.NewService(dhlClient, cfg.DhlApiKey, 10, 10)

		providers = append(providers, dhlService)

	}

	if cfg.EnableDpd {
		dpdClient := dpdclient.New(cb, cfg.DpdBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		dpdService := dpd.NewService(dpdClient, cfg.DpdApiKey, 10, 10)

		providers = append(providers, dpdService)

	}

	if cfg.EnableFedex {
		fedexClient := fedexclient.New(cb, cfg.FedexBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		fedexService := fedex.NewService(fedexClient, cfg.FedexApiKey, 10, 10)

		providers = append(providers, fedexService)

	}

	if cfg.EnableGls {
		glsClient := glsclient.New(cb, cfg.GlsBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		glsService := gls.NewService(glsClient, cfg.GlsApiKey, 10, 10)

		providers = append(providers, glsService)

	}

	if cfg.EnableInpost {
		inpostClient := inpostclient.New(cb, cfg.InpostBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		inpostService := inpost.NewService(inpostClient, cfg.InpostApiKey, 10, 10)

		providers = append(providers, inpostService)

	}

	if cfg.EnableUps {
		upsClient := upsclient.New(cb, cfg.UpsBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		upsService := ups.NewService(upsClient, cfg.UpsApiKey, 10, 10)

		providers = append(providers, upsService)

	}

	return providers
}
