package main

import (
	"github.com/brandonmakai/clipmux/internal"
	"github.com/brandonmakai/clipmux/internal/bootstrap"
)

// TODO: Change this to be a param in CLI
const capacity, maxItemBytes = 10, 1500
const maxBytes = capacity * maxItemBytes
const chronologicalHistory = false
const baseIndex = 1 // what index the first clibpoard item is at
const trackExistingItems = true // if the tool should copy the item already in the clipboard on startup
const debug = false
const loggerPath = ""

func main() {
	logger, history := bootstrap.BootStrap(capacity, maxItemBytes, maxBytes, chronologicalHistory, loggerPath, debug)

	clipRobot := internal.RobotGo{}
	// Change this to be customizeable path
	cm := internal.NewClipboardManager(clipRobot, history, logger)

	if err := cm.Run(); err != nil {
		panic(err)
	}
}
