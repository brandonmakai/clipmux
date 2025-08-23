package persistence

import (
	"sync"
	"time"
)

type ClipboardHistory interface {
	NewestIndex() int
	Contains(text string) bool
	Append(data []byte)
	Newest() (Item, bool)
	GetPos(idx int) (Item, bool)
	List()
}

type ChronologicalHistory struct {
	mu           sync.RWMutex
	buf          []Item
	head         int
	count        int
	capacity     int
	maxBytes     int 
	maxItemBytes int 
	currBytes    int
}

type RecentFirstHistory struct {
	mu           sync.RWMutex
	buf          []Item
	head         int 
	count        int
	capacity     int
	maxBytes     int 
	maxItemBytes int 
	currBytes    int
}

type Item struct {
	Data      []byte
	CreatedAt time.Time
}
