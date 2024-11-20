package logging

import (
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func Setup(env, path string) *slog.Logger {
	output := os.Stdout
	if path != "" {
		var err error
		if output, err = os.Open(path); err != nil {
			slog.Error("Cannot open log file!")
			output = os.Stdout
		}
	}

	switch env {
	case envDev:
		return slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		slog.Warn("Undefined log level, using level info")
		return slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}