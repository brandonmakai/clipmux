package persistence

import (
	"time"
	"fmt"
)

func newHistory(capacity int, maxBytes int, maxItemBytes int) *ClipboardHistory {
	return &ClipboardHistory{
		buf:          make([]Item, capacity),
		capacity:     capacity,
		maxBytes:     maxBytes,
		maxItemBytes: maxItemBytes,
	}
}

func (ch *ClipboardHistory) Newest() int {
	// Subtract to prevent off by one errors
	// Append increases count so we must decrement to get newest
	idx := (ch.head + ch.count - 1) % ch.capacity
	return idx
}

// TODO: Consider altering this to a small scale scan (last 50 items)?
func (ch *ClipboardHistory) Contains(text string) bool {
	for _, item := range ch.buf {
		if string(item.Data) == text {
			return true
		}
	}
	return false
}

func (ch *ClipboardHistory) Append(data []byte) {
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

func (ch *ClipboardHistory) evictOldest() {
	old := &ch.buf[ch.head]
	ch.currBytes -= len(old.Data)
	*old = Item{}
	ch.head = (ch.head + 1) % ch.capacity
	ch.count--
}

func (ch *ClipboardHistory) GetNewest() (Item, bool) {
	fmt.Println("GetNewest Called.")
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

func (ch *ClipboardHistory) GetPos(idx int) (Item, bool) {
	fmt.Println("GetPos Called.")
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
