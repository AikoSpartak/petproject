package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {

	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(l))
}
