package config

import "time"

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
