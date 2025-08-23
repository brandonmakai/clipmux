package persistence

import (
	"fmt"
	"sync"
)

var (
	instance ClipboardHistory
	once     sync.Once
)

func GetHistory(chronological bool, capacity int, maxItemBytes int, maxBytes int) ClipboardHistory {
	once.Do(func() {
		if chronological {
			instance = newChronologicalHistory(capacity, maxBytes, maxItemBytes)
		} else {
			instance = newRecentFirstHistory(capacity, maxBytes, maxItemBytes)
		}
		fmt.Println("Store instance created.")
	})
	return instance
}
