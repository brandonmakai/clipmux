package config

import (
	"github.com/BurntSushi/toml"
)

// TODO: Implement baseIndex and CaptureInitialClipboard
type Config struct {
	Capacity                int      // default = 10
	MaxItemBytes            int      // default = 2048
	AllowDuplicates         bool     // default = false
	NewestFirst             bool     // default = false
	BaseIndex               int      // default = 0
	PasteHotkeysBase        []string // default = 'ctrl + shift + h'
	CaptureInitialClipboard bool     // default = true
	Debug                   bool     // default = false
	LoggerDir               string   // default = '$HOME/./clipmux/logs'
	ConfigDir               string   // default = '$HOME/./clipmux'
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
