package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(level zerolog.Level) *zerolog.Logger {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Println(err)
	}
	consoleWriter := zerolog.ConsoleWriter{
		Out:          os.Stderr,
		TimeFormat:   time.RFC3339,
		NoColor:      true,
		TimeLocation: loc,
	}

	logger := zerolog.New(consoleWriter).
		Level(level).
		With().
		Caller().
		Timestamp().
		Logger()

	return &logger
}
