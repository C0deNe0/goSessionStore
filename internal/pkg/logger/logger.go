package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(env string) {
	var handler slog.Handler

	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	}

	Logger = slog.New(handler)

	slog.SetDefault(Logger)
}
