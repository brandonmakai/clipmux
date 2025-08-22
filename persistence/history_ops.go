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
	return ch.count
}

func (ch *ClipboardHistory) Append(data []byte) {
	fmt.Println("Append Triggered")
	fmt.Printf("Before String data: %s, Bytes data: %v\n", string(data), data)
	if ch.capacity > 0 && len(data) > ch.maxItemBytes {
		data = data[:ch.maxItemBytes]
	}
	fmt.Printf("After String data: %s, Bytes data: %v\n", string(data), data)

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
	fmt.Println("Index: ", idx)

	cp := make([]byte, len(data))
	copy(cp, data)

	// If there are still bytes - zero them (shouldn't occur due to evict function)
	if ch.count > ch.capacity && len(ch.buf[idx].Data) != 0 {
		ch.currBytes -= len(ch.buf[idx].Data)
	}

	fmt.Printf("Appended data: %v", cp)

	ch.buf[idx] = Item{Data: cp, CreatedAt: time.Now()}
	ch.count++
	ch.currBytes += len(data)
}

func (ch *ClipboardHistory) evictOldest() {
	old := &ch.buf[ch.head]
	ch.currBytes -= len(old.Data)
	*old = Item{}
	ch.count--
}

func (ch *ClipboardHistory) GetNewest() (Item, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if ch.count == 0 {
		return Item{}, false
	}

	// Subtract by 1 due to append statement incrementing the count
	idx := (ch.head + ch.count - 1) % ch.capacity
	fmt.Println("GetNewest Index: ", idx)
	it := ch.buf[idx]

	cp := make([]byte, len(it.Data))
	copy(cp, it.Data)
	it.Data = cp
	
	if it.Data == nil {
		fmt.Println("Newest data is nil!")
	}
	fmt.Println("Newest value: ", string(it.Data))

	return it, true
}

func (ch *ClipboardHistory) GetPos(idx int) (Item, bool) {
	fmt.Println("GetPos triggered.")
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if ch.count == 0 {
		return Item{}, false
	}

	it := ch.buf[idx]
	if it.Data == nil {
		return Item{}, false
	}

	cp := make([]byte, len(it.Data))
	copy(cp, it.Data)
	it.Data = cp

	return it, true
}
