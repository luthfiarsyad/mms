package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// L is the global logger instance accessible across packages after Init() call.
var L zerolog.Logger

// Init initializes a zerolog.Logger with the given log level.
// You typically call this once in main.go after config.Load().
func Init(level string) {
	// default: info
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	// configure global time format and pretty console output for dev mode
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	L = zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.SetGlobalLevel(lvl)
}

// Get returns the initialized logger instance (in case you prefer function access).
func Get() zerolog.Logger {
	return L
}
