package  bootstrap 

import (
	"os"
	"path/filepath"

	"github.com/brandonmakai/clipmux/internal/logger"
)

const CLIPMUX string = "clipmux"

func GetPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
			panic(err)
	}
	clipmuxPath := filepath.Join(home, CLIPMUX)

	return clipmuxPath
}

func BootStrap() *logger.Logger {
	if err := os.MkdirAll(GetPath(), 0755); err != nil {
		panic(err)
	}

	logger := logger.GetLogger(GetPath())
	return logger
}
