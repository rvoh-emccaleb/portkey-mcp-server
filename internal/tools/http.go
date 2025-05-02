package tools

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/middleware"
)

func MakePortkeyAPIRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := middleware.GetHTTPClient(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make http request to portkey: %w", err)
	}

	return resp, nil
}

func HandleHTTPError(resp *http.Response, respBody []byte, lgr *slog.Logger) *mcp.CallToolResult {
	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		lgr.Error("unauthorized access to portkey service",
			"status_code", resp.StatusCode,
			"response", string(respBody),
		)

		return mcp.NewToolResultError("unauthorized access to portkey service")

	case resp.StatusCode == http.StatusForbidden:
		lgr.Error("forbidden access to portkey service",
			"status_code", resp.StatusCode,
			"response", string(respBody),
		)

		return mcp.NewToolResultError("forbidden access to portkey service")

	case resp.StatusCode == http.StatusNotFound:
		lgr.Info("resource not found",
			"status_code", resp.StatusCode,
			"response", string(respBody),
		)

		return mcp.NewToolResultError("requested resource not found")

	case resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError:
		lgr.Info("invalid request to portkey service",
			"status_code", resp.StatusCode,
			"response", string(respBody),
		)

		return mcp.NewToolResultError("invalid request to portkey service")

	default:
		lgr.Error("portkey service error",
			"status_code", resp.StatusCode,
			"response", string(respBody),
		)

		return mcp.NewToolResultError("portkey service reported an error")
	}
}
