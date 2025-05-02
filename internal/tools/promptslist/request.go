package promptslist

// Request represents the request body for the Portkey Prompts List API.
type Request struct {
	// Optional parameters
	CollectionID string `json:"collection_id,omitempty"`
	WorkspaceID  string `json:"workspace_id,omitempty"`
	CurrentPage  *int   `json:"current_page,omitempty"`
	PageSize     *int   `json:"page_size,omitempty"`
	Search       string `json:"search,omitempty"`
}
