package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ctxKey string

const loggerKey = ctxKey("logger")

// GetLogger retrieves the request-scoped logger from the context.
// If no logger is found, returns the default logger.
func GetLogger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}

// WithHTTPRequestLogging adds request-scoped logging context for HTTP (SSE transport) requests.
func WithHTTPRequestLogging(ctx context.Context, r *http.Request) context.Context {
	reqLogger := GetLogger(ctx).With(
		"http_req_method", r.Method,
		"http_req_path", r.URL.Path,
		"http_req_remote_addr", r.RemoteAddr,
	)

	reqLogger.Debug("processing http request")

	return context.WithValue(ctx, loggerKey, reqLogger)
}

// WithToolCallLogging adds request-scoped logging context for MCP tool call requests.
func WithToolCallLogging(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reqLogger := GetLogger(ctx).With(
			"tool_name", req.Params.Name,
		)

		ctx = context.WithValue(ctx, loggerKey, reqLogger)

		reqLogger.Debug("processing tool call request")

		result, err := next(ctx, req)

		switch {
		case err != nil:
			reqLogger.Error("tool call request failed", "error", err)
		case result != nil && result.IsError:
			reqLogger.Info("tool call request completed with error response", "is_error", result.IsError)
		default:
			reqLogger.Debug("tool call request completed successfully")
		}

		return result, err
	}
}
