package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARNING
	ERROR
)

type Logger struct {
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warning: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
		error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.info.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05"), message)
}

func (l *Logger) Warning(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.warning.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05"), message)
}

func (l *Logger) Error(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.error.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05"), message)
}
