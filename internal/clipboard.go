package internal

import (
	"github.com/go-vgo/robotgo"
	"time"
)

type ReadPaster interface {
	Reader
	Paster
}

type Reader interface {
	Read() (string, error)
}

type Paster interface {
	Paste(text string) error
}

type RobotGo struct{}

func (r RobotGo) Read() (string, error) {
	text, err := robotgo.ReadAll()
	return text, err
}

func (p RobotGo) Paste(text string) error {
    // 1. Set the system clipboard
    robotgo.WriteAll(text)

    // 2. Optional: tiny delay to ensure clipboard is updated
    time.Sleep(50 * time.Millisecond)

    // 3. Simulate paste key combo
    // macOS: use "cmd", Windows/Linux: use "ctrl"
    robotgo.KeyTap("v", "cmd") // change "cmd" -> "ctrl" if needed

    return nil
}

func (p RobotGo) OldPaste(text string) error {
	err := robotgo.PasteStr(text)
	return err
}
