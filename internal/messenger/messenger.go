package messenger

import (
	"fmt"
	"log/slog"
)

type Messenger struct{
	Logger *slog.Logger
}

func (m *Messenger) Info(msg string) {
	fmt.Println(msg)
	m.Logger.Info("Message: " + msg)
}
