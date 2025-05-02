package config

import "fmt"

const (
	envPrefixTools        = "TOOLS"
	envPrefixPromptRender = "PROMPT_RENDER"
	envPrefixPromptsList  = "PROMPTS_LIST"
)

type Tools struct {
	PromptRender BaseTool `envconfig:"PROMPT_RENDER"`
	PromptsList  BaseTool `envconfig:"PROMPTS_LIST"`
}

func (t *Tools) Validate() error {
	envPrefixPromptRenderTool := fmt.Sprintf("%s_%s", envPrefixTools, envPrefixPromptRender)

	err := t.PromptRender.Validate(envPrefixPromptRenderTool)
	if err != nil {
		return fmt.Errorf("error validating prompt render tool: %w", err)
	}

	envPrefixPromptsListTool := fmt.Sprintf("%s_%s", envPrefixTools, envPrefixPromptsList)

	err = t.PromptsList.Validate(envPrefixPromptsListTool)
	if err != nil {
		return fmt.Errorf("error validating prompts list tool: %w", err)
	}

	return nil
}
