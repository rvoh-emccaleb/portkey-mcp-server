package promptrender

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/middleware"
)

const (
	toolName = "prompt_render"

	toolParamPromptID  = "prompt_id"
	toolParamPromptTag = "prompt_tag"
	toolParamVariables = "variables"

	errTextInternalError = "internal error while processing request"
)

var (
	ErrPromptIDRequired       = errors.New("prompt_id is required")
	ErrVariablesMustBeStrings = errors.New("all variable values must be strings")
)

type toolParams struct {
	promptID  string
	promptTag string
	variables map[string]string
}

func NewTool(portkeyCfg config.Portkey, toolCfg config.PromptRenderTool) tools.Tuple {
	description := "Render a Portkey prompt template by prompt slug and return the raw payload. This is a way to obtain " +
		"a prompt with optional variables substituted in. You can select specific versions of a prompt, or use the currently " +
		"published version."

	if toolCfg.Description != "" {
		description = toolCfg.Description
	}

	promptRenderTool := mcp.NewTool(
		toolName,
		mcp.WithDescription(description),
		mcp.WithString(toolParamPromptID,
			mcp.Required(),
			mcp.Description("The ID of the Portkey prompt to render. Specifically, this is the 'slug' of the prompt, if you "+
				"have used search tools to find this prompt."),
		),
		mcp.WithString(toolParamPromptTag,
			mcp.Description("Specific prompt version or label (e.g. '12', 'latest'). If omitted the published version is used."),
		),
		mcp.WithObject(toolParamVariables,
			mcp.Description("Variables object to substitute into the prompt template. The object should be a JSON object with "+
				"key-value pairs of string variable names to string values."),
		),
	)

	return tools.Tuple{
		Tool:    &promptRenderTool,
		Handler: promptRenderHandler(portkeyCfg),
		Enabled: toolCfg.Enabled,
	}
}

// promptRenderHandler calls the Portkey Prompt Render API and returns the result.
// Note: For validation errors (e.g. missing required fields), specific error messages are returned.
// For internal/system errors, generic error messages are returned while details are logged.
func promptRenderHandler(portkey config.Portkey) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lgr := middleware.GetLogger(ctx)

		params, err := getToolParams(request)
		if err != nil {
			lgr.Info("failed to get user-provided tool params from mcp request", "error", err)

			return mcp.NewToolResultErrorFromErr("invalid input", err), nil
		}

		url := createURL(portkey, params)

		body, err := createReqBody(params)
		if err != nil {
			lgr.Error("failed to create request body", "error", err)

			return mcp.NewToolResultError(errTextInternalError), nil
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			lgr.Error("failed to create http request", "error", err)

			return mcp.NewToolResultError(errTextInternalError), nil
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Portkey-Api-Key", string(portkey.APIKey))

		resp, err := tools.MakePortkeyAPIRequest(ctx, httpReq)
		if err != nil {
			lgr.Error("failed to call portkey api", "error", err)

			return mcp.NewToolResultError("failed to communicate with portkey service"), nil
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lgr.Error("failed to read response body", "error", err)

			return mcp.NewToolResultError("failed to process portkey response"), nil
		}

		if resp.StatusCode != http.StatusOK {
			return tools.HandleHTTPError(resp, respBody, lgr), nil
		}

		var portkeyResp Response
		if err := json.Unmarshal(respBody, &portkeyResp); err != nil {
			lgr.Error("invalid response format received from portkey service", "error", err)

			return mcp.NewToolResultError("received invalid response from portkey service"), nil
		}

		if !portkeyResp.Success {
			lgr.Error("portkey api returned success:false", "response", string(respBody))

			return mcp.NewToolResultError("portkey service reported failure"), nil
		}

		return mcp.NewToolResultText(string(respBody)), nil
	}
}

func getToolParams(request mcp.CallToolRequest) (toolParams, error) {
	promptID := mcp.ParseString(request, toolParamPromptID, "")
	if promptID == "" {
		return toolParams{}, ErrPromptIDRequired
	}

	promptTag := mcp.ParseString(request, toolParamPromptTag, "")

	var variables map[string]string

	rawVariables := mcp.ParseStringMap(request, toolParamVariables, nil)
	if rawVariables != nil {
		variables = make(map[string]string, len(rawVariables))

		for key, value := range rawVariables {
			str, ok := value.(string)
			if !ok {
				return toolParams{}, fmt.Errorf("%w: key %q has type %T", ErrVariablesMustBeStrings, key, value)
			}

			variables[key] = str
		}
	}

	return toolParams{
		promptID:  promptID,
		promptTag: promptTag,
		variables: variables,
	}, nil
}

func createURL(portkey config.Portkey, params toolParams) string {
	endpointID := params.promptID
	if params.promptTag != "" {
		endpointID = fmt.Sprintf("%s@%s", params.promptID, params.promptTag)
	}

	return fmt.Sprintf("%s/prompts/%s/render", portkey.BaseURL, endpointID)
}

func createReqBody(params toolParams) ([]byte, error) {
	//nolint:exhaustruct
	req := Request{
		Variables: params.variables,
	}
	if req.Variables == nil {
		// The API expects variables key even if empty.
		req.Variables = map[string]string{}
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return data, nil
}
