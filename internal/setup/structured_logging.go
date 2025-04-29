package setup

import (
	"log/slog"
	"os"
	"strings"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
)

const (
	AppName = "Portkey MCP Server"
)

func StructuredLogging(logLevel config.LogLevel, appVersion string) {
	var level slog.Level

	switch strings.ToLower(string(logLevel)) {
	case strings.ToLower(slog.LevelDebug.String()):
		level = slog.LevelDebug
	case strings.ToLower(slog.LevelInfo.String()):
		level = slog.LevelInfo
	case strings.ToLower(slog.LevelWarn.String()):
		level = slog.LevelWarn
	case strings.ToLower(slog.LevelError.String()):
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	//nolint:exhaustruct
	jsonHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(jsonHandler).With(
		"app_name", AppName,
		"app_version", appVersion,
	)

	slog.SetDefault(logger)
}
