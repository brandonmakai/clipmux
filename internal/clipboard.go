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
	Paste(text string) 
}

type RobotGo struct{}

func (r RobotGo) Read() (string, error) {
	text, err := robotgo.ReadAll()
	return text, err
}

func (r RobotGo) Paste(text string) {
	 robotgo.WriteAll(text)
	 time.Sleep(50 * time.Millisecond)
	 // TODO: Add OS specific "ctrl" for Linux/Windows
	 robotgo.KeyTap("v", "cmd")
}

