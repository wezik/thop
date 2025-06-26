package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"thop/internal/problem"
)

var log *slog.Logger

func Init(logFile string) error {
	var handler slog.Handler
	if logFile == "" {
		handler = slog.NewTextHandler(io.Discard, nil)
	} else {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		handler = slog.NewTextHandler(file, nil)
	}
	log = slog.New(handler)
	return nil
}

func Cmd(msg string) {
	log.Info("Executing in shell: " + msg)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(err error) {
	switch err := err.(type) {

	case problem.Problem:
		log.Error(string(err.Key) + ": " + err.Message)
		fmt.Println(err.Key+":", err.Message)

	default:
		log.Error(err.Error())
		fmt.Println("Unknown error:", err.Error())

	}
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Message(msg string) {
	log.Info("Brodcasting message: " + msg)
	fmt.Println(msg)
}
