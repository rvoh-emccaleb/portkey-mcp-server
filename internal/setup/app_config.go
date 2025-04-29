package setup

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
)

func AppConfig(buildTimeVars config.BuildTimeVars) (config.App, error) {
	var cfg config.App

	err := envconfig.Process("", &cfg)
	if err != nil {
		return config.App{}, fmt.Errorf("error loading environment variables into app config: %w", err)
	}

	cfg.BuildTimeVars = buildTimeVars

	err = cfg.Validate()
	if err != nil {
		return config.App{}, fmt.Errorf("error validating app config: %w", err)
	}

	return cfg, nil
}
