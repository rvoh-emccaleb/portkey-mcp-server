package types

import (
	"encoding/json"
)

const numExposedChars = 4 // number of characters to show at the end of mostly masked strings

// MostlyMaskedString is a type that mostly masks a string when String is called.
//
//nolint:recvcheck
type MostlyMaskedString string

// Decode implements the envconfig.Decoder interface.
func (mms *MostlyMaskedString) Decode(value string) error {
	*mms = MostlyMaskedString(value)

	return nil
}

// String mostly masks the underlying string value, showing only the last few characters.
func (mms MostlyMaskedString) String() string {
	if len(mms) <= numExposedChars {
		return "****"
	}

	return "****" + string(mms)[len(mms)-numExposedChars:]
}

// MarshalJSON implements the json.Marshaler interface to mask the value in JSON output.
// This ensures values are masked when using structured logging like slog with JSON handlers.
func (mms MostlyMaskedString) MarshalJSON() ([]byte, error) {
	// Use the String method to maintain consistent masking format
	return json.Marshal(mms.String())
}
