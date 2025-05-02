package config

// BaseTool holds base configuration for all tools.
type BaseTool struct {
	// Description will override the default description of the tool.
	Description string `envconfig:"DESCRIPTION" required:"false"`

	// Enabled will disable the tool if set to false.
	Enabled bool `default:"true" envconfig:"ENABLED" required:"false"`
}

// Validate validates the BaseTool configuration.
func (t *BaseTool) Validate(envPrefix string) error {
	_ = envPrefix + "_DESCRIPTION"

	return nil
}
