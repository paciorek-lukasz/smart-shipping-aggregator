package provider

import (
	"testing"

	"github.com/paciorek-lukasz/smart-shipping-aggregator/internal/config"
)

func TestInitProviders_AllEnabled(t *testing.T) {
	cfg := config.EnvConfig{
		EnableDhl:    true,
		EnableDpd:   true,
		EnableFedex: true,
		EnableGls:   true,
		EnableInpost: true,
		EnableUps:   true,
	}

	providers := InitProviders(cfg)

	if len(providers) != 6 {
		t.Errorf("expected 6 providers, got %d", len(providers))
	}
}

func TestInitProviders_OnlyDhl(t *testing.T) {
	cfg := config.EnvConfig{
		EnableDhl:    true,
		EnableDpd:   false,
		EnableFedex: false,
		EnableGls:   false,
		EnableInpost: false,
		EnableUps:   false,
	}

	providers := InitProviders(cfg)

	if len(providers) != 1 {
		t.Errorf("expected 1 provider, got %d", len(providers))
	}
}

func TestInitProviders_NoneEnabled(t *testing.T) {
	cfg := config.EnvConfig{
		EnableDhl:    false,
		EnableDpd:   false,
		EnableFedex: false,
		EnableGls:   false,
		EnableInpost: false,
		EnableUps:   false,
	}

	providers := InitProviders(cfg)

	if len(providers) != 0 {
		t.Errorf("expected 0 providers, got %d", len(providers))
	}
}

func TestInitProviders_MultipleEnabled(t *testing.T) {
	cfg := config.EnvConfig{
		EnableDhl:    true,
		EnableDpd:   true,
		EnableFedex: false,
		EnableGls:   false,
		EnableInpost: true,
		EnableUps:   false,
	}

	providers := InitProviders(cfg)

	if len(providers) != 3 {
		t.Errorf("expected 3 providers, got %d", len(providers))
	}
}