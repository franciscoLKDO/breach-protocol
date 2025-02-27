package config

import (
	_ "embed"
	"encoding/json"

	"fmt"
	"os"

	"github.com/franciscolkdo/breach-protocol/game/model"
)

//go:embed config.json
var configData []byte

type Config struct {
	Models []model.Config `json:"models"`
}

// NewGameConfig
func GetConfig(path string) (Config, error) {
	var err error
	if path != "" {
		if configData, err = os.ReadFile(path); err != nil {
			return Config{}, fmt.Errorf("error on loading config: %w", err)
		}
	}
	var cfg Config
	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error on unmarshal config data: %s", err)
	}
	return cfg, nil
}
