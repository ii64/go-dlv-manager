package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func Logger() zerolog.Logger {
	return log.Logger
}

func Close() error {
	return nil
}
