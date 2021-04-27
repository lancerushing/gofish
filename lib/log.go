package lib

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure zerolog.
func LogSetup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func CheckError(err error) {
	if err != nil {
		log.Panic().Err(err).Send()
	}
}
