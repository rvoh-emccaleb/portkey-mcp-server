package config_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

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
		BaseURL: "https://api.portkey.example.com",
		APIKey:  types.MaskedString(apiKey),
	}

	// Test string representation via fmt - should be masked
	strOutput := fmt.Sprintf("%v", cfg)

	if strings.Contains(strOutput, apiKey) {
		t.Error("APIKey not masked in string output")
	}

	if !strings.Contains(strOutput, "****") {
		t.Error("Expected masked characters **** not found in string output")
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)

	if !strings.Contains(jsonStr, apiKey) {
		t.Error("Expected APIKey to be present in JSON output")
	}
}
