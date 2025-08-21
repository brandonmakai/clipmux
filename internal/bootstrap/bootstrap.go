package  bootstrap 

import (
	"os"
	"path/filepath"
)

const CLIPMUX string = "clipmux"

func GetPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
			panic(err)
	}
	clipmuxPath := filepath.Join(home, CLIPMUX)

	return clipmuxPath
}

func BootStrap() {
	if err := os.MkdirAll(GetPath(), 0755); err != nil {
		panic(err)
	}
}
