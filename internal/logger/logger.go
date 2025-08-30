package logger

import (
	"os"
	"sync"
)

type Logger struct {
	file       *os.File
	mu         sync.Mutex
	currentDay string
	toStdout   bool
}
