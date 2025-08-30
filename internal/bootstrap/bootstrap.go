package bootstrap

import (
	"os"
	"path/filepath"
	"time"

	"github.com/brandonmakai/clipmux/internal/config"
	"github.com/brandonmakai/clipmux/internal/logger"
	"github.com/brandonmakai/clipmux/persistence"
)

const clipMux = ".clipmux"
const cfgFile = "config.toml"

// GetAppDir returns the root application directory for clipmux.
func GetAppDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user home directory: " + err.Error())
	}
	return filepath.Join(home, clipMux)
}

func BootStrap() (*logger.Logger, persistence.ClipboardHistory, *config.Config) {
	appDir := GetAppDir()
	if err := os.MkdirAll(appDir, 0755); err != nil {
		panic("failed to create application directory: " + err.Error())
	}

	cfg := config.GetConfig(filepath.Join(appDir, cfgFile))

	log := initLogger(cfg)
	history := initHistory(cfg, log)

	return log, history, cfg
}

// initLogger determines the log path, ensures the directory exists, and returns a new logger.
func initLogger(cfg *config.Config) *logger.Logger {
	logDir := cfg.LoggerDir
	if logDir == "" {
		// If no custom log directory is set in config, use the default.
		logDir = filepath.Join(GetAppDir(), "logs")
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")

	debug := cfg.Debug
	return logger.GetLogger(logFile, debug)
}

// initHistory determines capacity and maxItemBytes, provides defaults, and returns new history
func initHistory(cfg *config.Config, log *logger.Logger) persistence.ClipboardHistory {
	capacity := cfg.Capacity
	if capacity == 0 {
		capacity = 10
	}

	itemBytes := cfg.MaxItemBytes
	if itemBytes == 0 {
		itemBytes = 2048
	}

	newestFirst := cfg.NewestFirst
	return persistence.GetHistory(newestFirst, capacity, itemBytes, log)
}
