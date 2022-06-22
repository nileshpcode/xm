package log

import (
	"github.com/rs/zerolog"
	"io"
	"time"
)

// New creates new logger
func New(serviceName string, logLevel zerolog.Level, writers io.Writer) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339       // global settings
	zerolog.DurationFieldUnit = time.Millisecond // global settings
	zerolog.DurationFieldInteger = true          // global settings

	logger := zerolog.New(writers).Level(logLevel).With().Timestamp().Str("service", serviceName).Logger()
	return &logger
}
