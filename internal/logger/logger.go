package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	env := strings.ToLower(os.Getenv("APP_ENV"))

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", "devices-api").
		Logger().
		Level(level)

	// pretty logs for dev, structured for prod
	if env == "development" || env == "dev" {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	return logger
}

func parseLevel(lvl string) zerolog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
