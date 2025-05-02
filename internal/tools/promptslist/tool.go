package promptslist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/middleware"
)

const (
	toolName = "prompts_list"

	// Tool parameters.
	toolParamCollectionID = "collection_id"
	toolParamWorkspaceID  = "workspace_id"
	toolParamCurrentPage  = "current_page"
	toolParamPageSize     = "page_size"
	toolParamSearch       = "search"

	// Portkey API parameters.
	apiParamCollectionID = "collection_id"
	apiParamWorkspaceID  = "workspace_id"
	apiParamCurrentPage  = "current_page"
	apiParamPageSize     = "page_size"
	apiParamSearch       = "search"

	errTextInternalError = "internal error while processing request"
)

var (
	ErrInvalidPageSize    = fmt.Errorf("%s must be a positive integer", toolParamPageSize)
	ErrInvalidCurrentPage = fmt.Errorf("%s must be a positive integer", toolParamCurrentPage)
)

type toolParams struct {
	collectionID string
	workspaceID  string
	currentPage  *int
	pageSize     *int
	search       string
}

func NewTool(portkeyCfg config.Portkey, toolCfg config.BaseTool) tools.Tuple {
	description := "List all prompts in your Portkey account after applying the provided parameters. This tool allows " +
		"you to retrieve prompt metadata like prompt ID, prompt slug, name, collection, model, and status. " +
		"You can filter by various parameters and paginate results. There is the ability to search by " +
		"approximate name and slug matches."

	if toolCfg.Description != "" {
		description = toolCfg.Description
	}

	listPromptsTool := mcp.NewTool(
		toolName,
		mcp.WithDescription(description),
		mcp.WithString(toolParamCollectionID,
			mcp.Description("Optional. Filter prompts by collection ID."),
		),
		mcp.WithString(toolParamWorkspaceID,
			mcp.Description("Optional. Filter prompts by workspace ID."),
		),
		mcp.WithNumber(toolParamCurrentPage,
			mcp.Description("Optional. Page number for pagination. Starts at 1."),
		),
		mcp.WithNumber(toolParamPageSize,
			mcp.Description("Optional. Number of results per page."),
		),
		mcp.WithString(toolParamSearch,
			mcp.Description("Optional. Search term to filter prompts by name or slug."),
		),
	)

	return tools.Tuple{
		Tool:    &listPromptsTool,
		Handler: promptsListHandler(portkeyCfg),
		Enabled: toolCfg.Enabled,
	}
}

// promptsListHandler calls the Portkey Prompts List API and returns the result.
func promptsListHandler(portkey config.Portkey) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lgr := middleware.GetLogger(ctx)

		params, err := getToolParams(request)
		if err != nil {
			lgr.Info("failed to get user-provided tool params from mcp request", "error", err)

			return mcp.NewToolResultErrorFromErr("invalid input", err), nil
		}

		apiURL := createURL(portkey, params)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
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

func getToolParams(request mcp.CallToolRequest) (toolParams, error) {
	//nolint:exhaustruct
	params := toolParams{
		collectionID: mcp.ParseString(request, toolParamCollectionID, ""),
		workspaceID:  mcp.ParseString(request, toolParamWorkspaceID, ""),
		search:       mcp.ParseString(request, toolParamSearch, ""),
	}

	// Handle optional integer parameters
	currentPage := mcp.ParseInt(request, toolParamCurrentPage, 0)
	if currentPage > 0 {
		params.currentPage = &currentPage
	} else if currentPage < 0 {
		return toolParams{}, ErrInvalidCurrentPage
	}

	pageSize := mcp.ParseInt(request, toolParamPageSize, 0)
	if pageSize > 0 {
		params.pageSize = &pageSize
	} else if pageSize < 0 {
		return toolParams{}, ErrInvalidPageSize
	}

	return params, nil
}

func createURL(portkey config.Portkey, params toolParams) string {
	baseURL := portkey.BaseURL + "/prompts"

	// Add query parameters
	values := url.Values{}
	if params.collectionID != "" {
		values.Add(apiParamCollectionID, params.collectionID)
	}

	if params.workspaceID != "" {
		values.Add(apiParamWorkspaceID, params.workspaceID)
	}

	if params.currentPage != nil {
		values.Add(apiParamCurrentPage, strconv.Itoa(*params.currentPage))
	}

	if params.pageSize != nil {
		values.Add(apiParamPageSize, strconv.Itoa(*params.pageSize))
	}

	if params.search != "" {
		values.Add(apiParamSearch, params.search)
	}

	if len(values) > 0 {
		return baseURL + "?" + values.Encode()
	}

	return baseURL
}
