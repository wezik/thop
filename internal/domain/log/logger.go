package log

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(err error)
	Echo(msg string)
	Debug(msg string, args ...any)
}
