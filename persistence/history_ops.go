package persistence

import (
	"fmt"
	"time"
)

func newHistory(capacity int, maxBytes int, maxItemBytes int) *ClipboardHistory {
	return &ClipboardHistory{
		buf:          make([]Item, capacity),
		capacity:     capacity,
		maxBytes:     maxBytes,
		maxItemBytes: maxItemBytes,
	}
}

func (ch *ClipboardHistory) Append(data []byte) {
	if ch.capacity > 0 && len(data) > ch.maxItemBytes {
		data = data[:ch.maxItemBytes]
	}
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for {
		full := ch.capacity <= ch.count
		overBytes := ch.maxBytes > 0 && ch.currBytes+len(data) > ch.maxBytes

		if !full && !overBytes {
			break
		}

		ch.evictOldest()
	}

	idx := (ch.head + ch.count) % ch.capacity

	cp := make([]byte, len(data))
	copy(cp, data)

	// If there are still bytes - zero them (shouldn't occur due to evict function)
	if ch.count > ch.capacity && len(ch.buf[idx].Data) != 0 {
		ch.currBytes -= len(ch.buf[idx].Data)
	}

	ch.buf[idx] = Item{Data: cp, CreatedAt: time.Now()}
	ch.count++
	ch.currBytes += len(data)
}

func (ch *ClipboardHistory) evictOldest() {
	old := &ch.buf[ch.head]
	ch.currBytes = len(old.Data)
	*old = Item{}
	ch.count--
}

func (ch *ClipboardHistory) GetNewest() (Item, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if ch.count == 0 {
		return Item{}, false
	}

	idx := (ch.head + ch.count) % ch.capacity
	it := ch.buf[idx]

	cp := make([]byte, len(it.Data))
	cp = it.Data
	it.Data = cp

	return it, true
}
