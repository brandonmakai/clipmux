package persistence

import (
	"fmt"
	"sync"
)

var (
	instance *ClipboardHistory
	once     sync.Once
)

func GetHistory(capacity int, maxItemBytes int, maxBytes int) *ClipboardHistory {
	once.Do(func() {
		instance = newHistory(capacity, maxBytes, maxItemBytes)
		fmt.Println("Store instance created.")
	})
	return instance
}
