package config

import "fmt"

const (
	envPrefixTools        = "TOOLS"
	envPrefixPromptRender = "PROMPT_RENDER"
)

type Tools struct {
	PromptRender PromptRenderTool `envconfig:"PROMPT_RENDER"`
}

func (t *Tools) Validate() error {
	envPrefixPromptRenderTool := fmt.Sprintf("%s_%s", envPrefixTools, envPrefixPromptRender)

	err := t.PromptRender.Validate(envPrefixPromptRenderTool)
	if err != nil {
		return fmt.Errorf("error validating prompt render tool: %w", err)
	}

	return nil
}
