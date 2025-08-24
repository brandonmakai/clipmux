package bootstrap

import (
	"os"
	"time"
	"fmt"
	"path/filepath"

	"github.com/brandonmakai/clipmux/internal/logger"
	"github.com/brandonmakai/clipmux/persistence"
)

const clipMux string = "clipmux"

func GetPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	clipmuxPath := filepath.Join(home, clipMux)

	return clipmuxPath
}

func BootStrap(capacity int, maxItemBytes int, maxBytes int, chronologicalHistory bool, loggerPath string, debug bool) (*logger.Logger, persistence.ClipboardHistory) {
	if err := os.MkdirAll(GetPath(), 0755); err != nil {
		panic(err)
	}
	
	logger := initLogger(loggerPath, debug)
	history := persistence.GetHistory(chronologicalHistory, capacity, maxItemBytes, maxBytes)

	return logger, history
}

func initLogger(path string, debug bool) *logger.Logger {
	if path == "" { 
		path = fmt.Sprintf("%s/logs/%s.log", GetPath(), time.Now().Format("2006-01-02"))
	} else { 
		path = fmt.Sprintf("%s%s.log", path, time.Now().Format("2006-01-02"))
	}
	return logger.GetLogger(path, debug) 

}
