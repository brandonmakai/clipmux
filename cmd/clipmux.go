package main

import (
	"github.com/brandonmakai/clipmux/internal"
	"github.com/brandonmakai/clipmux/internal/bootstrap"
)

// TODO: Change this to be a param in CLI
const capacity, maxItemBytes = 10, 1500
const maxBytes = capacity * maxItemBytes

func main() {
	logger, history := bootstrap.BootStrap(capacity, maxItemBytes, maxBytes)

	clipRobot := internal.RobotGo{}
	// Change this to be customizeable path
	cm := internal.NewClipboardManager(clipRobot, history, logger)

	// TODO: Consider placing in goroutine (will need channel to handle error and select {})
	if err := cm.Run(); err != nil {
		panic(err)
	}
}
