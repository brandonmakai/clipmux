package internal

import (
	"github.com/go-vgo/robotgo"

	"time"
	"runtime"
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

	 // TODO: (Issue #1) Refactor Clipboard To Use OS-level Paste 
	 pasteKeys := []string{}
	 switch runtime.GOOS {
 	 case "darwin":
		pasteKeys = append(pasteKeys, "v", "cmd")
	 default: 
	 	pasteKeys = append(pasteKeys, "v", "ctrl")
	 }
	 robotgo.KeyTap(pasteKeys[0], pasteKeys[1])
}

