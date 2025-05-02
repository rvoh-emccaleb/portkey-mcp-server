package promptcreate

// Response represents the full response structure from the Portkey Prompt Create API.
type Response struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	VersionID string `json:"version_id"`
	Object    string `json:"object"`
}
