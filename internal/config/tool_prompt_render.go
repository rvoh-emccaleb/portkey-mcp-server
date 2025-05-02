package config

type PromptRenderTool struct {
	BaseTool
}

func (t *PromptRenderTool) Validate(envPrefix string) error {
	return t.BaseTool.Validate(envPrefix)
}
