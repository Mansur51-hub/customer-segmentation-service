package logger

import (
	"github.com/rs/zerolog"
)

func SetLevel(level string) {
	lvl, err := zerolog.ParseLevel(level)

	if err != nil {
		lvl = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
