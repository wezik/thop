package log

type NoopLogger struct {
}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (l *NoopLogger) Debug(msg string, args ...any) {}

func (l *NoopLogger) Echo(msg string) {}

func (l *NoopLogger) Error(err error) {}

func (l *NoopLogger) Info(msg string, args ...any) {}

func (l *NoopLogger) Warn(msg string, args ...any) {}
