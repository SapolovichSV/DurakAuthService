package logger

import (
	"log/slog"
	"os"

	"github.com/SapolovichSV/durak/auth/internal/config"
)

type Logger struct {
	*slog.Logger
}

func New(config config.Config) Logger {
	handlerOptions := slog.HandlerOptions{
		Level: slog.Level(config.LogLevel),
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))
	return Logger{
		Logger: logger,
	}
}
