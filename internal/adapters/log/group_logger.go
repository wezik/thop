package log

import (
	"thop/internal/domain/log"
)

type GroupLogger struct {
	Loggers []log.Logger
}

func NewGroupLogger(loggers []log.Logger) log.Logger {
	return &GroupLogger{
		Loggers: loggers,
	}
}

func (s *GroupLogger) Info(msg string, args ...any) {
	for _, logger := range s.Loggers {
		if logger == nil {
			continue
		}
		logger.Info(msg, args...)
	}
}

func (s *GroupLogger) Warn(msg string, args ...any) {
	for _, logger := range s.Loggers {
		if logger == nil {
			continue
		}
		logger.Warn(msg, args...)
	}
}

func (s *GroupLogger) Error(err error) {
	for _, logger := range s.Loggers {
		if logger == nil {
			continue
		}
		logger.Error(err)
	}
}

func (g *GroupLogger) Echo(msg string) {
	for _, logger := range g.Loggers {
		if logger == nil {
			continue
		}
		logger.Echo(msg)
	}
}

func (g *GroupLogger) Debug(msg string, args ...any) {
	for _, logger := range g.Loggers {
		if logger == nil {
			continue
		}
		logger.Debug(msg, args...)
	}
}
