package config

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	once     sync.Once
	instance *Config
	err      error
)

func GetConfig(path string) *Config {
	once.Do(func() {
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(err)
		}

		if _, err = InitConfig(path); err != nil {
			panic(err)
		}

		if instance, err = LoadConfig(path); err != nil {
			panic(err)
		}
	})
	return instance
}
