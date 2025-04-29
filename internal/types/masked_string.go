package types

import (
	"encoding/json"
	"fmt"
)

// MaskedString is a type that fully masks a string when String is called.
//
//nolint:recvcheck
type MaskedString string

// Decode implements the envconfig.Decoder interface.
func (ms *MaskedString) Decode(value string) error {
	*ms = MaskedString(value)

	return nil
}

// String fully masks the underlying string value.
func (ms MaskedString) String() string {
	return "****"
}

// MarshalJSON implements the json.Marshaler interface to mask the value in JSON output.
// This ensures values are masked when using structured logging like slog with JSON handlers.
func (ms MaskedString) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal("****")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal masked string to json: %w", err)
	}

	return data, nil
}
