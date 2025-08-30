package logger

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *Logger
	once     sync.Once
)

func GetLogger(path string, debug bool) *Logger {
	once.Do(func() {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		instance = newLogger(f, debug)
		instance.Info("Logger instance created.")
	})
	return instance
}
