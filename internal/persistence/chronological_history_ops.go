package persistence

import (
	"fmt"
	"time"
)

func newChronologicalHistory(capacity int, maxBytes int, maxItemBytes int) *ChronologicalHistory {
	return &ChronologicalHistory{
		buf:          make([]Item, capacity),
		capacity:     capacity,
		maxBytes:     maxBytes,
		maxItemBytes: maxItemBytes,
	}
}

func (ch *ChronologicalHistory) NewestIndex() int {
	idx := (ch.head + ch.count - 1) % ch.capacity
	return idx
}

// TODO: (Issue #3) Optimize Clipboard Contains Function
func (ch *ChronologicalHistory) Contains(text string) bool {
	for _, item := range ch.buf {
		// Check for non-nil data to avoid matching empty slots
		if item.Data != nil && string(item.Data) == text {
			return true
		}
	}
	return false
}

func (ch *ChronologicalHistory) Append(data []byte) {
	if len(data) > ch.maxItemBytes {
		data = data[:ch.maxItemBytes]
	}

	ch.mu.Lock()
	defer ch.mu.Unlock()

	for ch.count == ch.capacity || (ch.maxBytes > 0 && ch.currBytes+len(data) > ch.maxBytes) {
		ch.evictOldest()
	}

	idx := (ch.head + ch.count) % ch.capacity
	cp := append([]byte(nil), data...)

	ch.buf[idx] = Item{Data: cp, CreatedAt: time.Now()}
	ch.count++
	ch.currBytes += len(data)
}

func (ch *ChronologicalHistory) evictOldest() {
	old := &ch.buf[ch.head]
	ch.currBytes -= len(old.Data)
	*old = Item{}
	ch.head = (ch.head + 1) % ch.capacity
	ch.count--
}

func (ch *ChronologicalHistory) Newest() (Item, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if ch.count == 0 {
		return Item{}, false
	}
	idx := (ch.head + ch.count - 1) % ch.capacity
	it := ch.buf[idx]
	return Item{
		Data:      append([]byte(nil), it.Data...),
		CreatedAt: it.CreatedAt,
	}, true
}

func (ch *ChronologicalHistory) GetPos(idx int) (Item, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if ch.count == 0 {
		return Item{}, false
	}

	it := ch.buf[idx]
	return Item{
		Data:      append([]byte(nil), it.Data...),
		CreatedAt: it.CreatedAt,
	}, true
}

func (ch *ChronologicalHistory) List() {
	for idx, item := range ch.buf {
		fmt.Printf("Item: %v, Index: %v ", string(item.Data), idx)
	}
	fmt.Println()
}
