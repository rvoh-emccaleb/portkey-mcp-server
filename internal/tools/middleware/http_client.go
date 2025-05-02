package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const clientKey = ctxKey("http_client")

const defaultClientTimeout = 30 * time.Second

func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:       defaultClientTimeout,
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
	}
}

// WithHTTPClient adds the provided HTTP client to the context.
// This allows handlers to retrieve and use a shared HTTP client for all requests.
func WithHTTPClient(client *http.Client) Middleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx = context.WithValue(ctx, clientKey, client)

			return next(ctx, request)
		}
	}
}

// GetHTTPClient retrieves the HTTP client from the context.
// If no client is found in the context, a default HTTP client is created and returned.
func GetHTTPClient(ctx context.Context) *http.Client {
	client, ok := ctx.Value(clientKey).(*http.Client)
	if !ok {
		return DefaultHTTPClient()
	}

	return client
}
