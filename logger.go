package flex

import (
	"fmt"
)

type logMessage struct {
	Message string `json:"message"`
	Level   string `json:"level"`
}

// Logger ...
type Logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
	Fatal(message string)
}

type logger struct {
}

func newLogger() Logger {
	l := &logger{}
	return l
}

func (l *logger) log(level string, message string) {
	if message == "" {
		return
	}
	lm := logMessage{
		Message: message,
		Level:   level,
	}
	bytes, err := json.Marshal(lm)
	if err != nil {
	}
	fmt.Println(string(bytes))
}

// Info ...
func (l *logger) Info(message string) {
	l.log("info", message)
}

// Warn ...
func (l *logger) Warn(message string) {
	l.log("warning", message)
}

// Error ...
func (l *logger) Error(message string) {
	l.log("error", message)
}

// Fatal ...
func (l *logger) Fatal(message string) {
	l.log("fatal", message)
}
