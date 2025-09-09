package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"thop/internal/domain/log"
)

type EchoLogLevel int

const (
	Echo EchoLogLevel = iota
	Error
	Warn
	Info
	Debug
)

type EchoLogger struct {
	level EchoLogLevel
}

func NewEchoLogger(level EchoLogLevel) log.Logger {
	return &EchoLogger{
		level: level,
	}
}

func (e *EchoLogger) Info(msg string, args ...any) {
	if e.level >= Info {
		fmt.Println("Info: " + msg)
	}
}

func (e *EchoLogger) Warn(msg string, args ...any) {
	if e.level >= Warn {
		fmt.Println("Warn: " + msg)
	}
}

func (e *EchoLogger) Error(err error) {
	if e.level >= Error {
		fmt.Println("Error: " + err.Error())
	}
}

func (e *EchoLogger) Echo(msg string) {
	fmt.Println(msg)
}

func (e *EchoLogger) Debug(msg string, args ...any) {
	pc, _, _, ok := runtime.Caller(2) // 2 since it's called first by the group logger
	var fnName string
	if !ok {
		fnName = "???"
	} else {
		fnName = filepath.Base(runtime.FuncForPC(pc).Name())
	}

	if e.level >= Debug {
		fmt.Println("Debug [" + fnName + "]: " + msg)
	}
}
