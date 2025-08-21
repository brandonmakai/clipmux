package internal

import (
	"github.com/go-vgo/robotgo"
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
	err := robotgo.PasteStr(text)
	return err
}
