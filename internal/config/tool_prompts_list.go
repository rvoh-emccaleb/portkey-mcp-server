package config

type PromptsListTool struct {
	BaseTool
}

func (t *PromptsListTool) Validate(envPrefix string) error {
	return t.BaseTool.Validate(envPrefix)
}
