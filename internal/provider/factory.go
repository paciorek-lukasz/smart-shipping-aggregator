package provider

import (
	"time"

	"github.com/dzwiedz90/smart-shipping-aggregator/internal/config"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dhl"
	dhlclient "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dhl/client"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dpd"
	dpdclient "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/dpd/client"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/fedex"
	fedexclient "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/fedex/client"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/gls"
	glsclient "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/gls/client"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/inpost"
	inpostclient "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/inpost/client"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/ups"
	upsclient "github.com/dzwiedz90/smart-shipping-aggregator/internal/provider/ups/client"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/resilience"
)

func InitProviders(cfg config.EnvConfig) []Provider {
	var providers []Provider

	if cfg.EnableDhl {
		dhlCB := resilience.NewCircuitBreaker(5, 30*time.Second)

		dhlClient := dhlclient.New(dhlCB, cfg.DhlBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		dhlService := dhl.NewService(dhlClient, cfg.DhlApiKey, 10, 10)

		providers = append(providers, dhlService)

	}

	if cfg.EnableDpd {
		dpdCB := resilience.NewCircuitBreaker(5, 30*time.Second)

		dpdClient := dpdclient.New(dpdCB, cfg.DpdBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		dpdService := dpd.NewService(dpdClient, cfg.DpdApiKey, 10, 10)

		providers = append(providers, dpdService)

	}

	if cfg.EnableFedex {
		fedexCB := resilience.NewCircuitBreaker(5, 30*time.Second)

		fedexClient := fedexclient.New(fedexCB, cfg.FedexBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		fedexService := fedex.NewService(fedexClient, cfg.FedexApiKey, 10, 10)

		providers = append(providers, fedexService)

	}

	if cfg.EnableGls {
		glsCB := resilience.NewCircuitBreaker(5, 30*time.Second)

		glsClient := glsclient.New(glsCB, cfg.GlsBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		glsService := gls.NewService(glsClient, cfg.GlsApiKey, 10, 10)

		providers = append(providers, glsService)

	}

	if cfg.EnableInpost {
		inpostCB := resilience.NewCircuitBreaker(5, 30*time.Second)

		inpostClient := inpostclient.New(inpostCB, cfg.InpostBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		inpostService := inpost.NewService(inpostClient, cfg.InpostApiKey, 10, 10)

		providers = append(providers, inpostService)

	}

	if cfg.EnableUps {
		upsCB := resilience.NewCircuitBreaker(5, 30*time.Second)

		upsClient := upsclient.New(upsCB, cfg.UpsBaseUrl, 0.1, 100*time.Millisecond, 500*time.Millisecond)

		upsService := ups.NewService(upsClient, cfg.UpsApiKey, 10, 10)

		providers = append(providers, upsService)

	}

	return providers
}
