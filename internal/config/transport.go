package config

import (
	"errors"
	"fmt"
)

var ErrInvalidTransportType = errors.New("invalid transport type")

type TransportType string

const (
	TransportStdio TransportType = "stdio"
	TransportSSE   TransportType = "sse"
)

// Decode implements the envconfig.Decoder interface.
func (t *TransportType) Decode(value string) error {
	val := TransportType(value)
	switch val {
	case TransportStdio, TransportSSE:
		*t = val

		return nil
	default:
		return fmt.Errorf("%w: %q, must be one of: %s, %s", ErrInvalidTransportType, value, TransportStdio, TransportSSE)
	}
}
