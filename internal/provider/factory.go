package provider

import (
	"time"

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

type EnvConfig struct {
	EnableDhl  bool   `envconfig:"ENABLE_DHL"`
	DhlBaseUrl string `envconfig:"DHL_BASE_URL"`
	DhlApiKey  string `envconfig:"DHL_API_KEY"`

	EnableDpd  bool   `envconfig:"ENABLE_DPD"`
	DpdBaseUrl string `envconfig:"DPD_BASE_URL"`
	DpdApiKey  string `envconfig:"DPD_API_KEY"`

	EnableFedex  bool   `envconfig:"ENABLE_FEDEX"`
	FedexBaseUrl string `envconfig:"FEDEX_BASE_URL"`
	FedexApiKey  string `envconfig:"FEDEX_API_KEY"`

	EnableGls  bool   `envconfig:"ENABLE_GLS"`
	GlsBaseUrl string `envconfig:"GLS_BASE_URL"`
	GlsApiKey  string `envconfig:"GLS_API_KEY"`

	EnableInpost  bool   `envconfig:"ENABLE_INPOST"`
	InpostBaseUrl string `envconfig:"INPOST_BASE_URL"`
	InpostApiKey  string `envconfig:"INPOST_API_KEY"`

	EnableUps  bool   `envconfig:"ENABLE_UPS"`
	UpsBaseUrl string `envconfig:"UPS_BASE_URL"`
	UpsApiKey  string `envconfig:"UPS_API_KEY"`

	Timeout time.Duration `envconfig:"TIMEOUT"`
}

func InitProviders(cfg EnvConfig) []Provider {
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
