package config

import "fmt"

type App struct {
	BuildTimeVars
	EnvVars
}

func (cfg *App) Validate() error {
	err := cfg.BuildTimeVars.Validate()
	if err != nil {
		return fmt.Errorf("error validating build time config: %w", err)
	}

	if err := cfg.EnvVars.Validate(); err != nil {
		return fmt.Errorf("error validating run time environment config: %w", err)
	}

	return nil
}
