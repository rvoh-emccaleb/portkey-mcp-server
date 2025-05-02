package promptcreate

// Request represents the request body for the Portkey Prompt Create API.
type Request struct {
	// Required arguments
	Name         string         `json:"name"`
	CollectionID string         `json:"collection_id"`
	String       string         `json:"string"`
	Parameters   map[string]any `json:"parameters"`

	// Optional arguments
	Functions          []map[string]any `json:"functions,omitempty"`
	Tools              []map[string]any `json:"tools,omitempty"`
	ToolChoice         map[string]any   `json:"tool_choice,omitempty"`
	Model              string           `json:"model,omitempty"`
	VirtualKey         string           `json:"virtual_key,omitempty"`
	VersionDescription string           `json:"version_description,omitempty"`
	TemplateMetadata   map[string]any   `json:"template_metadata,omitempty"`
}
