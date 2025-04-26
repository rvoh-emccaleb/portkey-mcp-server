package promptrender

// Response represents the full response structure from the Portkey Prompt Render API.
type Response struct {
	Success bool       `json:"success"`
	Data    RenderData `json:"data"`
}

// RenderData represents the data field in the Portkey Prompt Render API response.
type RenderData struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	// Optional parameters
	FrequencyPenalty  *float64        `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int  `json:"logit_bias,omitempty"`
	LogProbs          bool            `json:"logprobs,omitempty"`
	TopLogProbs       *int            `json:"top_logprobs,omitempty"`
	MaxTokens         *int            `json:"max_tokens,omitempty"`
	N                 *int            `json:"n,omitempty"`
	PresencePenalty   *float64        `json:"presence_penalty,omitempty"`
	ResponseFormat    *ResponseFormat `json:"response_format,omitempty"`
	Seed              *int64          `json:"seed,omitempty"`
	Stop              interface{}     `json:"stop,omitempty"` // Can be string or []string
	Stream            bool            `json:"stream,omitempty"`
	StreamOptions     *StreamOptions  `json:"stream_options,omitempty"`
	Thinking          *ThinkingConfig `json:"thinking,omitempty"`
	Temperature       *float64        `json:"temperature,omitempty"`
	TopP              *float64        `json:"top_p,omitempty"`
	Tools             []Tool          `json:"tools,omitempty"`
	ToolChoice        interface{}     `json:"tool_choice,omitempty"` // Can be string or object
	ParallelToolCalls bool            `json:"parallel_tool_calls,omitempty"`
	User              string          `json:"user,omitempty"`
	FunctionCall      interface{}     `json:"function_call,omitempty"` // Deprecated
	Functions         []Function      `json:"functions,omitempty"`     // Deprecated
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type ThinkingConfig struct {
	Type         string `json:"type"`
	BudgetTokens int    `json:"budget_tokens,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}
