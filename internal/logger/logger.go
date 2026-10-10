package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {

	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	app := os.Getenv("APP")
	if app == "" {
		app = "default"
	}

	host := os.Getenv("DOMAIN")
	if host == "" {
		host = "default"
	}

	replica := os.Getenv("REPLICA")
	if replica == "" {
		replica = "default"
	}

	logger := slog.New(l).With(
		slog.Group("metadata",
			slog.String("app", os.Getenv("APP")),
			slog.String("host", os.Getenv("DOMAIN"))),
		slog.String("replica", os.Getenv("REPLICA")),
	)

	slog.SetDefault(logger)
}
