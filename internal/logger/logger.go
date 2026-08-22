package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {

	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(l).With(
		slog.Group("metadata",
			slog.String("app", os.Getenv("APP")),
			slog.String("host", os.Getenv("DOMAIN"))),
		slog.String("replica", os.Getenv("REPLICA")),
	)

	slog.SetDefault(logger)
}
