package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func init() {
	// Console writer for stdout (human-readable for systemd/journalctl)
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// JSON for file logging (structured, searchable)
	file, err := os.OpenFile("/tmp/saio.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		// Fall back to stdout only
		Logger = log.Output(console)
		return
	}

	multi := zerolog.MultiLevelWriter(console, file)
	Logger = log.Output(multi)
}

func Info(msg string, args ...interface{}) {
	Logger.Info().Msgf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Logger.Debug().Msgf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Logger.Warn().Msgf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Logger.Error().Msgf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	Logger.Fatal().Msgf(msg, args...)
}

func SetLevel(level string) {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)
}
