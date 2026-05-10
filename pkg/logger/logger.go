package logger

import (
	"fmt"
	"os"
	"time"
)

func stamp() string {
	return time.Now().Format(time.DateTime)
}

func log(w *os.File, level string, msg string) {
	fmt.Fprintf(w, "[%s] %s: %s\n", stamp(), level, msg)
}

func Info(msg string) {
	log(os.Stdout, "INFO", msg)
}

func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

func Warn(msg string) {
	log(os.Stderr, "WARN", msg)
}

func Warnf(format string, args ...any) {
	Warn(fmt.Sprintf(format, args...))
}

func Error(msg string) {
	log(os.Stderr, "ERROR", msg)
}

func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}
