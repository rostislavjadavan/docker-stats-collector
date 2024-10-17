package lib

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"os"
)

func SetupLogger(logLevel string) {
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
