package log

import (
	"io"
	"log/slog"
	"os"
	"thop/internal/domain/log"
	"thop/internal/adapters/platform"
)

type LogFile string

type Slog struct {
	logger *slog.Logger
}

func NewFileLogger(
	logFile LogFile,
	platform platform.Platform,
) (logger log.Logger, err error) {
	var handler slog.Handler

	if logFile == "" {
		handler = slog.NewTextHandler(io.Discard, nil)
	} else {
		file, err := platform.OpenFile(string(logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		handler = slog.NewTextHandler(file, nil)
	}
	return &Slog{
		logger: slog.New(handler),
	}, nil
}

func (s *Slog) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s *Slog) Warn(msg string, args ...any) {
	s.logger.Warn(msg, args...)
}

func (s *Slog) Error(err error) {
	s.logger.Error(err.Error())
}

func (s *Slog) Echo(msg string) {
	s.logger.Info("Echo: " + msg)
}

func (s *Slog) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}
