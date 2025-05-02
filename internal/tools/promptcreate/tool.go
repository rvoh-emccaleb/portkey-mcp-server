package promptcreate

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
	toolName = "prompt_create"

	// Tool arguments.
	toolArgName               = "name"
	toolArgCollectionID       = "collection_id"
	toolArgString             = "string"
	toolArgParameters         = "parameters"
	toolArgFunctions          = "functions"
	toolArgTools              = "tools"
	toolArgToolChoice         = "tool_choice"
	toolArgModel              = "model"
	toolArgVirtualKey         = "virtual_key"
	toolArgVersionDescription = "version_description"
	toolArgTemplateMetadata   = "template_metadata"

	errTextInternalError = "internal error while processing request"
)

var (
	ErrNameRequired         = errors.New("name is required")
	ErrCollectionIDRequired = errors.New("collection_id is required")
	ErrStringRequired       = errors.New("string is required")
	ErrParametersRequired   = errors.New("parameters is required")
	ErrInvalidArrayFormat   = errors.New("invalid array format: expected array of objects")
)

type toolArgs struct {
	name               string
	collectionID       string
	promptString       string
	parameters         map[string]any
	functions          []map[string]any
	tools              []map[string]any
	toolChoice         map[string]any
	model              string
	virtualKey         string
	versionDescription string
	templateMetadata   map[string]any
}

func NewTool(portkeyCfg config.Portkey, toolCfg config.BaseTool) tools.Tuple {
	description := "Create a new prompt in your Portkey account with the provided arguments. " +
		"This tool allows you to create a prompt with a name, template string, parameters, " +
		"and other optional settings."

	if toolCfg.Description != "" {
		description = toolCfg.Description
	}

	promptCreateTool := mcp.NewTool(
		toolName,
		mcp.WithDescription(description),
		mcp.WithString(toolArgName,
			mcp.Required(),
			mcp.Description("Name of the prompt to create."),
		),
		mcp.WithString(toolArgCollectionID,
			mcp.Required(),
			mcp.Description("UUID or slug of the collection to add the prompt to."),
		),
		mcp.WithString(toolArgString,
			mcp.Required(),
			mcp.Description("Prompt template in string format. Use {{variable_name}} syntax "+
				"to define variables that can be substituted at runtime (e.g., 'Hello {{name}}, how are you?')."),
		),
		mcp.WithObject(toolArgParameters,
			mcp.Required(),
			mcp.Description("Parameters for the prompt. This defines the variable schema for the template. "+
				"Each key in this object will be available as {{key}} in the prompt template. Uses Mustache templating - "+
				"keys should be the variable names with values as expected data types "+
				"(e.g., {\"name\": \"string\", \"age\": \"number\"}). "+
				"At runtime, users will provide actual values for these variables."),
		),
		mcp.WithArray(toolArgFunctions,
			mcp.Description("Functions for the prompt."),
		),
		mcp.WithArray(toolArgTools,
			mcp.Description("Tools for the prompt."),
		),
		mcp.WithObject(toolArgToolChoice,
			mcp.Description("Tool Choice for the prompt."),
		),
		mcp.WithString(toolArgModel,
			mcp.Description("The model to use for the prompt."),
		),
		mcp.WithString(toolArgVirtualKey,
			mcp.Description("The virtual key to use for the prompt."),
		),
		mcp.WithString(toolArgVersionDescription,
			mcp.Description("The description of the prompt version."),
		),
		mcp.WithObject(toolArgTemplateMetadata,
			mcp.Description("Metadata for the prompt."),
		),
	)

	return tools.Tuple{
		Tool:    &promptCreateTool,
		Handler: promptCreateHandler(portkeyCfg),
		Enabled: toolCfg.Enabled,
	}
}

// promptCreateHandler calls the Portkey Prompt Create API and returns the result.
// Note: For validation errors (e.g. missing required fields), specific error messages are returned.
// For internal/system errors, generic error messages are returned while details are logged.
func promptCreateHandler(portkey config.Portkey) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lgr := middleware.GetLogger(ctx)

		args, err := getToolArguments(request)
		if err != nil {
			lgr.Info("failed to get user-provided tool arguments from mcp request", "error", err)

			return mcp.NewToolResultErrorFromErr("invalid input", err), nil
		}

		url := portkey.BaseURL + "/prompts"

		body, err := createReqBody(args)
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

		return mcp.NewToolResultText(string(respBody)), nil
	}
}

func getToolArguments(request mcp.CallToolRequest) (toolArgs, error) {
	name := mcp.ParseString(request, toolArgName, "")
	if name == "" {
		return toolArgs{}, ErrNameRequired
	}

	collectionID := mcp.ParseString(request, toolArgCollectionID, "")
	if collectionID == "" {
		return toolArgs{}, ErrCollectionIDRequired
	}

	promptString := mcp.ParseString(request, toolArgString, "")
	if promptString == "" {
		return toolArgs{}, ErrStringRequired
	}

	parameters := mcp.ParseStringMap(request, toolArgParameters, nil)
	if parameters == nil {
		return toolArgs{}, ErrParametersRequired
	}

	// Optional arguments
	functions, err := extractArrayOfObjects(request, toolArgFunctions)
	if err != nil {
		return toolArgs{}, err
	}

	tools, err := extractArrayOfObjects(request, toolArgTools)
	if err != nil {
		return toolArgs{}, err
	}

	toolChoice := mcp.ParseStringMap(request, toolArgToolChoice, nil)
	model := mcp.ParseString(request, toolArgModel, "")
	virtualKey := mcp.ParseString(request, toolArgVirtualKey, "")
	versionDescription := mcp.ParseString(request, toolArgVersionDescription, "")
	templateMetadata := mcp.ParseStringMap(request, toolArgTemplateMetadata, nil)

	return toolArgs{
		name:               name,
		collectionID:       collectionID,
		promptString:       promptString,
		parameters:         parameters,
		functions:          functions,
		tools:              tools,
		toolChoice:         toolChoice,
		model:              model,
		virtualKey:         virtualKey,
		versionDescription: versionDescription,
		templateMetadata:   templateMetadata,
	}, nil
}

func extractArrayOfObjects(request mcp.CallToolRequest, argName string) ([]map[string]any, error) {
	rawValue, exists := request.Params.Arguments[argName]
	if !exists || rawValue == nil {
		return nil, nil
	}

	rawArray, ok := rawValue.([]any)
	if !ok {
		return nil, fmt.Errorf("%w for argument %q", ErrInvalidArrayFormat, argName)
	}

	result := make([]map[string]any, 0, len(rawArray))

	for i, item := range rawArray {
		obj, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%w for argument %q at index %d", ErrInvalidArrayFormat, argName, i)
		}

		result = append(result, obj)
	}

	return result, nil
}

func createReqBody(args toolArgs) ([]byte, error) {
	req := Request{
		Name:               args.name,
		CollectionID:       args.collectionID,
		String:             args.promptString,
		Parameters:         args.parameters,
		Functions:          args.functions,
		Tools:              args.tools,
		ToolChoice:         args.toolChoice,
		Model:              args.model,
		VirtualKey:         args.virtualKey,
		VersionDescription: args.versionDescription,
		TemplateMetadata:   args.templateMetadata,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return data, nil
}
