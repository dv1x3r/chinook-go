package logger

import (
	"log/slog"
	"os"
)

type HandlerType int

const (
	_ HandlerType = iota
	TextHandler
	JsonHandler
)

var (
	level  slog.LevelVar
	logger *slog.Logger
)

func Create(handlerType HandlerType) {
	switch handlerType {
	case TextHandler:
		logger = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: &level},
			),
		)
	case JsonHandler:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: &level},
			),
		)
	}
}

func SetLevel(l slog.Level) {
	level.Set(l)
}

func Get() *slog.Logger {
	return logger
}
