package  bootstrap 

import (
	"os"
	"path/filepath"

	"github.com/brandonmakai/clipmux/internal/logger"
	"github.com/brandonmakai/clipmux/persistence"
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

func BootStrap(capacity int, maxItemBytes int, maxBytes int) (*logger.Logger, *persistence.ClipboardHistory) {
	if err := os.MkdirAll(GetPath(), 0755); err != nil {
		panic(err)
	}

	logger := logger.GetLogger(GetPath())
	history := persistence.GetHistory(capacity, maxItemBytes, maxBytes)

	return logger, history
}
