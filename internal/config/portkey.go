package config

import (
	"errors"
	"fmt"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/types"
)

var (
	ErrBaseURLRequired = errors.New("base url is required")
	ErrAPIKeyRequired  = errors.New("api key is required")
)

type Portkey struct {
	APIKey  types.MaskedString `envconfig:"API_KEY"                 json:"api_key"       required:"true"`
	BaseURL string             `default:"https://api.portkey.ai/v1" envconfig:"BASE_URL" json:"base_url" required:"true"` //nolint:lll
	Client  HTTPClient         `envconfig:"CLIENT"                  json:"client"`
}

func (cfg *Portkey) Validate() error {
	if cfg.BaseURL == "" {
		return ErrBaseURLRequired
	}

	if cfg.APIKey == "" {
		return ErrAPIKeyRequired
	}

	if err := cfg.Client.Validate(); err != nil {
		return fmt.Errorf("invalid http client: %w", err)
	}

	return nil
}
