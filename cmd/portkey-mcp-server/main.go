package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/setup"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/tools/middleware"
)

// Build-time variables.
var (
	appVersion string //nolint:gochecknoglobals

	ErrUnknownTransportType = errors.New("unknown transport type")
)

const (
	serverShutdownTimeout = 10 * time.Second
)

func main() {
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.Error("error loading .env file", "error", err)

		return
	}

	buildTimeVars := config.BuildTimeVars{
		AppVersion: appVersion,
	}

	cfg, err := setup.AppConfig(buildTimeVars)
	if err != nil {
		slog.Error("error setting up application config", "error", err)

		return
	}

	setup.StructuredLogging(cfg.LogLevel, cfg.AppVersion)

	slog.Info("starting up...")
	slog.Info("using config", "config", cfg) // hides sensitive values

	srv, errChan, err := startServer(cfg)
	if err != nil {
		slog.Error("error starting server", "error", err)

		return
	}

	waitForShutdown(rootCtx, errChan)
	shutdownServer(srv)

	slog.Info("goodbye!")
}

func startServer(cfg config.App) (*server.SSEServer, chan error, error) {
	mcpServer := server.NewMCPServer(
		setup.AppName,
		cfg.AppVersion,
		server.WithLogging(),
		server.WithRecovery(),
	)

	if err := setup.MCPTools(cfg, mcpServer); err != nil {
		return nil, nil, fmt.Errorf("failed to register tools: %w", err)
	}

	errChan := make(chan error, 1)

	switch cfg.Transport {
	case config.TransportStdio:
		go func() {
			slog.Info("starting stdio server")

			if err := server.ServeStdio(mcpServer); err != nil {
				errChan <- fmt.Errorf("stdio server error: %w", err)
			}
		}()

		return nil, errChan, nil

	case config.TransportSSE:
		sseServer := server.NewSSEServer(
			mcpServer,
			[]server.SSEOption{
				server.WithSSEContextFunc(middleware.WithHTTPRequestLogging),
			}...,
		)

		go func() {
			slog.Info("starting sse server", "address", cfg.TransportSSE.Address)

			if err := sseServer.Start(cfg.TransportSSE.Address); err != nil {
				errChan <- fmt.Errorf("sse server error: %w", err)
			}
		}()

		return sseServer, errChan, nil

	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrUnknownTransportType, cfg.Transport)
	}
}

func waitForShutdown(ctx context.Context, errChan chan error) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		slog.Error("server error", "error", err)
	case sig := <-osSignals:
		slog.Info("signal received from os", "signal", sig)
	case <-ctx.Done():
		slog.Info("context cancelled")
	}
}

func shutdownServer(srv *server.SSEServer) {
	// `srv` is nil if we're using the stdio transport.
	if srv == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	slog.Info("shutting down sse server...")

	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("sse server shutdown error", "error", err)

		return
	}

	slog.Info("sse server has been shut down gracefully")
}
