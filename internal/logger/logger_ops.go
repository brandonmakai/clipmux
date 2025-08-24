package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func newLogger(file *os.File, toStdout bool) *Logger {
	return &Logger{file: file, toStdout: toStdout}
}

func (l *Logger) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	line := fmt.Sprintf("%s [INFO] %s\n", time.Now().Format(time.RFC3339), msg)
	fmt.Fprintf(l.file, line) 

	if l.toStdout {
		fmt.Println(line)
	}
}

func (l *Logger) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	line := fmt.Sprintf("%s [FATAL] %s\n", time.Now().Format(time.RFC3339), msg) 
	fmt.Fprintf(l.file,line) 
	log.Fatal(msg)

	if l.toStdout {
		fmt.Println(line)
	}
}
