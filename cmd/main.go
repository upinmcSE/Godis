package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/upinmcSE/godis/internal/server"
	"github.com/upinmcSE/godis/pkg/logger"
)

func main() {
	fmt.Println("Godis")

	LOG := logger.NewLogger(zerolog.TraceLevel)

	LOG.Trace().Msg("trace message")
	LOG.Info().Msg("info message")
	LOG.Debug().Msg("debug message")
	LOG.Error().Msg("error message")

	server.RunServer(LOG)
}
