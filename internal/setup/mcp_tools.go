package setup

import (
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/server"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/middleware"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/promptrender"
)

func MCPTools(cfg config.App, mcpServer *server.MCPServer, downstreamTools ...tools.Tuple) error {
	middlewares := []middleware.Middleware{
		middleware.WithToolCallLogging,
	}

	ts := []tools.Tuple{
		promptrender.NewTool(cfg.Portkey, cfg.Tools.PromptRender),
	}

	ts = append(ts, downstreamTools...)

	for _, t := range ts {
		slog.Info(fmt.Sprintf("registering %s tool", t.Tool.Name))

		mcpServer.AddTool(
			*t.Tool,
			addMiddleware(t.Handler, middlewares...),
		)
	}

	return nil
}

// addMiddleware applies the middleware in reverse order, to ensure that a request flows through middleware in the order
// they are passed in to this function. In other words, this is necessary because of the way that wrapping works.
func addMiddleware(handler server.ToolHandlerFunc, middlewares ...middleware.Middleware) server.ToolHandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
