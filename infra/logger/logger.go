package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Logger is a type alias on underlying logger implementation
type Logger = zerolog.Logger

// Config will hold logger specific configuration
type Config struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

// NewLogger will instantiate and configure the logger
func NewLogger(c Config) (Logger, error) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	lvl, err := zerolog.ParseLevel(strings.ToLower(c.Level))
	if err != nil {
		return Logger{}, fmt.Errorf("can't parse level: %w", err)
	}
	logger.Level(lvl)

	return logger, nil
}
