package promptslist

import "time"

// Response represents the full response structure from the Portkey List Prompts API.
type Response struct {
	Data  []PromptData `json:"data"`
	Total int          `json:"total"`
}

// PromptData represents a single prompt entry in the Prompts List API response.
type PromptData struct {
	ID            string    `json:"id"`
	Slug          string    `json:"slug"`
	Name          string    `json:"name"`
	CollectionID  string    `json:"collection_id"`
	Model         string    `json:"model"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Object        string    `json:"object"`
}
