package config

import "fmt"

const (
	envPrefixTools        = "TOOLS"
	envPrefixPromptCreate = "PROMPT_CREATE"
	envPrefixPromptRender = "PROMPT_RENDER"
	envPrefixPromptsList  = "PROMPTS_LIST"
)

type Tools struct {
	PromptCreate BaseTool `envconfig:"PROMPT_CREATE"`
	PromptRender BaseTool `envconfig:"PROMPT_RENDER"`
	PromptsList  BaseTool `envconfig:"PROMPTS_LIST"`
}

func (t *Tools) Validate() error {
	envPrefixPromptCreateTool := fmt.Sprintf("%s_%s", envPrefixTools, envPrefixPromptCreate)

	err := t.PromptCreate.Validate(envPrefixPromptCreateTool)
	if err != nil {
		return fmt.Errorf("error validating prompt create tool: %w", err)
	}

	envPrefixPromptRenderTool := fmt.Sprintf("%s_%s", envPrefixTools, envPrefixPromptRender)

	err = t.PromptRender.Validate(envPrefixPromptRenderTool)
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
