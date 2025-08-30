package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxFiles = 3

func newLogger(file *os.File, toStdout bool) *Logger {
	return &Logger{file: file, toStdout: toStdout}
}

// enforceRetention removes outdated files in the logger directory if there is more than maxFiles entries 
func (l *Logger) enforceRetention() error {
	dir := filepath.Dir(l.file.Name())
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var names []string
	for _, e := range entries {
		entryName := e.Name()
		if !e.IsDir() && filepath.Ext(entryName) == ".log" {
			names = append(names, entryName)
		}
	}

	if len(names) <= maxFiles {
		return nil
	}
	
	sort.Strings(names)
	
	for i := 0; i <= maxFiles; i++ {
		if err := os.Remove(filepath.Join(dir, names[i])); err != nil {
			return fmt.Errorf("remove %q: %w", names[i], err)
		}
	}

	return nil
}

// rotateIfNeeded rotates logger files daily to avoid excessive disk usage 
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

	if err := l.enforceRetention(); err != nil {
		panic(err)
	}

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
