package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	file    *os.File
}

// NewLogger initializes the logger and writes WARNING and ERROR logs to a file
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		debug:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
		info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warning: log.New(file, "WARNING: ", log.Ldate|log.Ltime),
		error:   log.New(file, "ERROR: ", log.Ldate|log.Ltime),
		file:    file,
	}, nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.debug.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05"), message)
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

// Close closes the log file
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
