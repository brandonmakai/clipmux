package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	instance *Logger
	once     sync.Once
)

// TODO: Add daily rotation of logs
func GetLogger(path string) *Logger {
	once.Do(func() {
		logger_path := fmt.Sprintf("%s/logs/%s.log", path, time.Now().Format("2006-01-02"))
		if err := os.MkdirAll(filepath.Dir(logger_path), 0755); err != nil {
			panic(err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		instance = newLogger(f)
		fmt.Println("Logger instance created.")
	})
	return instance
}
