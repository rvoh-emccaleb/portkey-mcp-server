package setup

import (
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/server"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/middleware"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/promptrender"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/promptslist"
)

func MCPTools(cfg config.App, mcpServer *server.MCPServer, downstreamTools ...tools.Tuple) error {
	httpClient, err := cfg.Portkey.Client.FromConfig()
	if err != nil {
		return fmt.Errorf("failed to create http client from config: %w", err)
	}

	middlewares := []middleware.Middleware{
		middleware.WithToolCallLogging,
		middleware.WithHTTPClient(httpClient),
	}

	allTools := []tools.Tuple{
		promptrender.NewTool(cfg.Portkey, cfg.Tools.PromptRender),
		promptslist.NewTool(cfg.Portkey, cfg.Tools.PromptsList),
	}

	allTools = append(allTools, downstreamTools...)

	enabledTools := getEnabledTools(allTools)

	for _, t := range enabledTools {
		slog.Info(fmt.Sprintf("registering %s tool", t.Tool.Name))

		mcpServer.AddTool(
			*t.Tool,
			addMiddleware(t.Handler, middlewares...),
		)
	}

	return nil
}

// getEnabledTools filters out disabled tools from the provided list.
func getEnabledTools(allTools []tools.Tuple) []tools.Tuple {
	enabledTools := make([]tools.Tuple, 0, len(allTools))

	for _, tool := range allTools {
		if tool.Enabled {
			enabledTools = append(enabledTools, tool)
		}
	}

	return enabledTools
}

// addMiddleware applies the middleware in reverse order, to ensure that a request flows through middleware in the order
// they are passed in to this function. In other words, this is necessary because of the way that wrapping works.
func addMiddleware(handler server.ToolHandlerFunc, middlewares ...middleware.Middleware) server.ToolHandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
