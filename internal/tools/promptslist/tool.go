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

	// Tool arguments.
	toolArgCollectionID = "collection_id"
	toolArgWorkspaceID  = "workspace_id"
	toolArgCurrentPage  = "current_page"
	toolArgPageSize     = "page_size"
	toolArgSearch       = "search"

	// Portkey API query parameters.
	apiParamCollectionID = "collection_id"
	apiParamWorkspaceID  = "workspace_id"
	apiParamCurrentPage  = "current_page"
	apiParamPageSize     = "page_size"
	apiParamSearch       = "search"

	errTextInternalError = "internal error while processing request"
)

var (
	ErrInvalidPageSize    = fmt.Errorf("%s must be a positive integer", toolArgPageSize)
	ErrInvalidCurrentPage = fmt.Errorf("%s must be a positive integer", toolArgCurrentPage)
)

type toolArgs struct {
	collectionID string
	workspaceID  string
	currentPage  *int
	pageSize     *int
	search       string
}

func NewTool(portkeyCfg config.Portkey, toolCfg config.BaseTool) tools.Tuple {
	description := "List all prompts in your Portkey account after applying the provided arguments. This tool allows " +
		"you to retrieve prompt metadata like prompt ID, prompt slug, name, collection, model, and status. " +
		"You can filter by various parameters and paginate results. There is the ability to search by " +
		"approximate name and slug matches."

	if toolCfg.Description != "" {
		description = toolCfg.Description
	}

	listPromptsTool := mcp.NewTool(
		toolName,
		mcp.WithDescription(description),
		mcp.WithString(toolArgCollectionID,
			mcp.Description("Optional. Filter prompts by collection ID."),
		),
		mcp.WithString(toolArgWorkspaceID,
			mcp.Description("Optional. Filter prompts by workspace ID."),
		),
		mcp.WithNumber(toolArgCurrentPage,
			mcp.Description("Optional. Page number for pagination. Starts at 1."),
		),
		mcp.WithNumber(toolArgPageSize,
			mcp.Description("Optional. Number of results per page."),
		),
		mcp.WithString(toolArgSearch,
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

		args, err := getToolArguments(request)
		if err != nil {
			lgr.Info("failed to get user-provided tool arguments from mcp request", "error", err)

			return mcp.NewToolResultErrorFromErr("invalid input", err), nil
		}

		apiURL := createURL(portkey, args)

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

func getToolArguments(request mcp.CallToolRequest) (toolArgs, error) {
	//nolint:exhaustruct
	args := toolArgs{
		collectionID: mcp.ParseString(request, toolArgCollectionID, ""),
		workspaceID:  mcp.ParseString(request, toolArgWorkspaceID, ""),
		search:       mcp.ParseString(request, toolArgSearch, ""),
	}

	// Handle optional integer arguments.
	currentPage := mcp.ParseInt(request, toolArgCurrentPage, 0)
	if currentPage > 0 {
		args.currentPage = &currentPage
	} else if currentPage < 0 {
		return toolArgs{}, ErrInvalidCurrentPage
	}

	pageSize := mcp.ParseInt(request, toolArgPageSize, 0)
	if pageSize > 0 {
		args.pageSize = &pageSize
	} else if pageSize < 0 {
		return toolArgs{}, ErrInvalidPageSize
	}

	return args, nil
}

func createURL(portkey config.Portkey, args toolArgs) string {
	baseURL := portkey.BaseURL + "/prompts"

	// Add query parameters.
	values := url.Values{}
	if args.collectionID != "" {
		values.Add(apiParamCollectionID, args.collectionID)
	}

	if args.workspaceID != "" {
		values.Add(apiParamWorkspaceID, args.workspaceID)
	}

	if args.currentPage != nil {
		values.Add(apiParamCurrentPage, strconv.Itoa(*args.currentPage))
	}

	if args.pageSize != nil {
		values.Add(apiParamPageSize, strconv.Itoa(*args.pageSize))
	}

	if args.search != "" {
		values.Add(apiParamSearch, args.search)
	}

	if len(values) > 0 {
		return baseURL + "?" + values.Encode()
	}

	return baseURL
}
