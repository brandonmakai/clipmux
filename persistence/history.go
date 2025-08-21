package persistence

import (
	"sync"
	"time"
)

type ClipboardHistory struct {
	mu sync.RWMutex
	buf []Item 
	head int // index of head of store 
	count int 
	capacity int 
	maxBytes int // max byte size of store 
	maxItemBytes int // max per item byte size 
	currBytes int // current amount of bytes in store
}

type Item struct {
	Data []byte
	CreatedAt time.Time
}

