package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	DateFormat string `json:"DateFormat"`
	BaseRoute  string `json:"BaseRoute"`
	Port       string `json:"Port"`
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads config from config.json (singleton pattern)
func LoadConfig(filePath string) (*Config, error) {
	var err error
	once.Do(func() {
		file, e := os.Open(filePath)
		if e != nil {
			err = e
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		cfg = &Config{}
		err = decoder.Decode(cfg)
	})
	return cfg, err
}
