package main

import (
	"github.com/brandonmakai/clipmux/internal"
	"github.com/brandonmakai/clipmux/internal/bootstrap"
)

func main() {
	logger, history, cfg := bootstrap.BootStrap()

	clipRobot := internal.RobotGo{}
	cm := internal.NewClipboardManager(clipRobot, history, logger, cfg)

	if err := cm.Run(); err != nil {
		panic(err)
	}
}
