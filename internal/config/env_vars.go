package config

import "fmt"

type EnvVars struct {
	LogLevel     LogLevel `default:"info" envconfig:"LOG_LEVEL"`
	Portkey      Portkey
	Tools        Tools
	Transport    TransportType `default:"stdio"           envconfig:"TRANSPORT"`
	TransportSSE SSETransport  `envconfig:"TRANSPORT_SSE"`
}

func (cfg *EnvVars) Validate() error {
	if err := cfg.Portkey.Validate(); err != nil {
		return fmt.Errorf("error validating portkey config: %w", err)
	}

	if err := cfg.Tools.Validate(); err != nil {
		return fmt.Errorf("error validating tools config: %w", err)
	}

	return nil
}
