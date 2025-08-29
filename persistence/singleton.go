package persistence

import (
	"sync"

	"github.com/brandonmakai/clipmux/internal/logger"
)

var (
	instance ClipboardHistory
	once     sync.Once
)

func GetHistory(newestFirst bool, capacity int, maxItemBytes int, logger *logger.Logger) ClipboardHistory {
	maxBytes := capacity * maxItemBytes
	once.Do(func() {
		if newestFirst {
			instance = newRecentFirstHistory(capacity, maxBytes, maxItemBytes)
		} else {
			instance = newChronologicalHistory(capacity, maxBytes, maxItemBytes)
		}
		logger.Info("Store instance created.")
	})
	return instance
}
