package config

type PromptRenderTool struct {
	// Description will override the default description of the tool.
	Description string `envconfig:"DESCRIPTION" required:"false"`
}

func (t *PromptRenderTool) Validate(envPrefix string) error {
	_ = envPrefix + "_DESCRIPTION"

	return nil
}
