package config

import (
	"errors"
	"fmt"
)

var ErrInvalidLogLevel = errors.New("invalid log level")

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Decode implements the envconfig.Decoder interface.
func (l *LogLevel) Decode(value string) error {
	val := LogLevel(value)
	switch val {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		*l = val

		return nil
	default:
		return fmt.Errorf("%w: %q, must be one of: %s, %s, %s, %s",
			ErrInvalidLogLevel, value, LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
	}
}
