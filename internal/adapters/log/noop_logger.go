package log

type NoopLogger struct {
}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (l *NoopLogger) Debug(msg string, args ...any) {
	// do nothing
}

func (l *NoopLogger) Echo(msg string) {
	// do nothing
}

func (l *NoopLogger) Error(err error) {
	// do nothing
}

func (l *NoopLogger) Info(msg string, args ...any) {
	// do nothing
}

func (l *NoopLogger) Warn(msg string, args ...any) {
	// do nothing
}
