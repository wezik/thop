package logger

import (
	"io"
	"log/slog"
	"os"
)

func New(logFile string) (*slog.Logger, error) {
	if logFile == "" {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return slog.New(slog.NewTextHandler(file, nil)), nil
}