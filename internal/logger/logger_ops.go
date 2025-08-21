package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func newLogger(file *os.File) *Logger {
	return &Logger{file: file}
}

func (l *Logger) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "%s [INFO] %s\n", time.Now().Format(time.RFC3339), msg)
}

func (l *Logger) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "%s [FATAL] %s\n", time.Now().Format(time.RFC3339), msg)
	log.Fatal(msg)
}
