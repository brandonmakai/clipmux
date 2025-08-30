package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func newLogger(file *os.File, toStdout bool) *Logger {
	return &Logger{file: file, toStdout: toStdout}
}

func (l *Logger) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")
	if today == l.currentDay && l.file != nil {
		return nil
	}

	filename := l.file.Name()
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)

	newFile := today + ext
	newFile = filepath.Join(dir, newFile)
	f, err := os.OpenFile(newFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.file = f
	l.currentDay = today
	return nil
}

func (l *Logger) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.rotateIfNeeded(); err != nil {
		panic(err)
	}

	line := fmt.Sprintf("%s [INFO] %s\n", time.Now().Format(time.RFC3339), msg)
	if _, err := fmt.Fprint(l.file, line); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
	}

	if l.toStdout {
		fmt.Println(line)
	}
}

func (l *Logger) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.rotateIfNeeded(); err != nil {
		panic(err)
	}

	line := fmt.Sprintf("%s [FATAL] %s\n", time.Now().Format(time.RFC3339), msg)
	if _, err := fmt.Fprint(l.file, line); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
	}

	if l.toStdout {
		fmt.Print(line)
	}

	os.Exit(1)
}
