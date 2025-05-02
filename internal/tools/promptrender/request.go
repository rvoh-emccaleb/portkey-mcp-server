package promptrender

// Request represents the request body for the Portkey Prompt Render API.
type Request struct {
	// Required arguments
	Variables map[string]string `json:"variables"`

	// Optional arguments - all at root level as per API spec
	Messages          []Message       `json:"messages,omitempty"`
	Model             string          `json:"model,omitempty"`
	FrequencyPenalty  *float64        `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int  `json:"logit_bias,omitempty"`
	LogProbs          bool            `json:"logprobs,omitempty"`
	TopLogProbs       *int            `json:"top_logprobs,omitempty"`
	MaxTokens         *int            `json:"max_tokens,omitempty"`
	N                 *int            `json:"n,omitempty"`
	PresencePenalty   *float64        `json:"presence_penalty,omitempty"`
	ResponseFormat    *ResponseFormat `json:"response_format,omitempty"`
	Seed              *int64          `json:"seed,omitempty"`
	Stop              any             `json:"stop,omitempty"` // Can be string or []string
	Stream            bool            `json:"stream,omitempty"`
	StreamOptions     *StreamOptions  `json:"stream_options,omitempty"`
	Thinking          *ThinkingConfig `json:"thinking,omitempty"`
	Temperature       *float64        `json:"temperature,omitempty"`
	TopP              *float64        `json:"top_p,omitempty"`
	Tools             []Tool          `json:"tools,omitempty"`
	ToolChoice        any             `json:"tool_choice,omitempty"` // Can be string or object
	ParallelToolCalls bool            `json:"parallel_tool_calls,omitempty"`
	User              string          `json:"user,omitempty"`
	FunctionCall      any             `json:"function_call,omitempty"` // Deprecated
	Functions         []Function      `json:"functions,omitempty"`     // Deprecated
}
