package persistence

import (
	"time"
	"fmt"
)

func newRecentFirstHistory(capacity int, maxBytes int, maxItemBytes int) *RecentFirstHistory {
	return &RecentFirstHistory{
		buf:          make([]Item, capacity),
		capacity:     capacity,
		maxBytes:     maxBytes,
		maxItemBytes: maxItemBytes,
	}
}

// logicalToPhysical maps a logical history index (0 = newest, count-1 = oldest)
// into the actual physical index inside the circular buffer.
func (rh *RecentFirstHistory) logicalToPhysical(i int) int {
    return (rh.head + rh.count - 1 - i + rh.capacity) % rh.capacity
}

func (rh *RecentFirstHistory) NewestIndex() int {
	return rh.logicalToPhysical(0)
}

// TODO: Consider altering this to a small scale scan (last 50 items)?
func (rh *RecentFirstHistory) Contains(text string) bool {
	for _, item := range rh.buf {
		if string(item.Data) == text {
			return true
		}
	}
	return false
}

func (rh *RecentFirstHistory) Append(data []byte) {
	if len(data) > rh.maxItemBytes {
		data = data[:rh.maxItemBytes]
	}

	rh.mu.Lock()
	defer rh.mu.Unlock()

	for rh.count == rh.capacity || (rh.maxBytes > 0 && rh.currBytes+len(data) > rh.maxBytes) {
		rh.evictOldest()
	}

	idx := (rh.head + rh.count) % rh.capacity
	cp := append([]byte(nil), data...)

	rh.buf[idx] = Item{Data: cp, CreatedAt: time.Now()}
	rh.count++
	rh.currBytes += len(data)
}

func (rh *RecentFirstHistory) evictOldest() {
	old := &rh.buf[rh.head]
	rh.currBytes -= len(old.Data)
	*old = Item{}
	rh.head = (rh.head + 1) % rh.capacity
	rh.count--
}

func (rh *RecentFirstHistory) Newest() (Item, bool) {
	fmt.Println("RecentFirstHistory::Newest Called.")
	rh.mu.RLock()
	defer rh.mu.RUnlock()

	if rh.count == 0 {
		return Item{}, false
	}
	it := rh.buf[rh.logicalToPhysical(0)]
	return Item{
		Data:      append([]byte(nil), it.Data...),
		CreatedAt: it.CreatedAt,
	}, true
}

func (rh *RecentFirstHistory) GetPos(idx int) (Item, bool) {
	fmt.Println("RecentFirstHistory::GetPos Called.")
	rh.mu.RLock()
	defer rh.mu.RUnlock()

	if rh.count == 0 {
		return Item{}, false
	}

	rh.List()
	it := rh.buf[rh.logicalToPhysical(idx)]
	return Item{
		Data:      append([]byte(nil), it.Data...),
		CreatedAt: it.CreatedAt,
	}, true
}

func (rh *RecentFirstHistory) List() {
	fmt.Println("RecentFirstHistory::List Called.")
	for idx, item := range rh.buf {
		fmt.Printf("Item: %v, Index: %v ", string(item.Data), idx)
	}
	fmt.Println()
}
