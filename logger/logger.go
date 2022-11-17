package logger

import (
	. "alperb/dirgod/arguments"
	"fmt"
)

type Logger struct {
	IS_DEBUG   bool
	IS_VERBOSE bool
}

func NewLogger(args Arguments) *Logger {
	return &Logger{args.DEBUG_MODE, args.VERBOSE_MODE}
}

func (l *Logger) Debug(msg string) {
	if l.IS_DEBUG || l.IS_VERBOSE {
		fmt.Println("[DEBUG] " + msg)
	}
}

func (l *Logger) Log(msg string) {
	if l.IS_VERBOSE {
		fmt.Println("[INFO] " + msg)
	}
}
