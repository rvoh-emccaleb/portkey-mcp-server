package config_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/config"
	"github.com/rvoh-emccaleb/portkey-mcp-server/internal/types"
)

func TestPortkeyConfigMasking(t *testing.T) {
	t.Parallel()

	const (
		apiKey     = "super-secret-api-key-123"
		virtualKey = "very-secret-virtual-key-456"
	)

	cfg := config.Portkey{
		APIKey:  types.MaskedString(apiKey),
		BaseURL: "https://api.portkey.example.com",
		Client: config.HTTPClient{
			CustomCACertPath:   "/foo/bar",
			InsecureSkipVerify: false,
			Timeout:            10 * time.Second,
		},
	}

	// Test string representation via fmt - should be masked
	strOutput := fmt.Sprintf("%v", cfg)

	if strings.Contains(strOutput, apiKey) {
		t.Error("APIKey not masked in string output")
	}

	if !strings.Contains(strOutput, "****") {
		t.Error("Expected masked characters **** not found in string output")
	}

	// Test direct JSON marshaling
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)

	// With the new MarshalJSON implementation, the API key should be masked in JSON
	if strings.Contains(jsonStr, apiKey) {
		t.Error("APIKey not masked in direct JSON output")
	}

	if !strings.Contains(jsonStr, "****") {
		t.Error("Expected masked characters **** not found in JSON output")
	}

	// Test with slog JSON handler (simulating main.go)
	var logBuf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{ //nolint:exhaustruct
		Level: slog.LevelInfo,
	})
	logger := slog.New(jsonHandler)

	// This mirrors how it's logged in main.go
	logger.Info("using config", "portkey", cfg)

	logOutput := logBuf.String()

	if strings.Contains(logOutput, apiKey) {
		t.Error("APIKey not masked in slog JSON output")
	}

	if !strings.Contains(logOutput, "****") {
		t.Error("Expected masked characters **** not found in slog JSON output")
	}
}
