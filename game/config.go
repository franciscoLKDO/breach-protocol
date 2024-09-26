package game

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/breach-protocol/game/breach"
	"github.com/franciscolkdo/breach-protocol/game/end"
	"github.com/franciscolkdo/breach-protocol/game/story"
)

type model string

const (
	breachModel model = "breach"
	storyModel  model = "story"
	endModel    model = "end"
)

type ModelConfig struct {
	Type   model           `json:"type"`
	Config json.RawMessage `json:"config"`
}

type Config struct {
	Models []ModelConfig `json:"models"`
}

func (c Config) LoadModel(idx int) (tea.Model, error) {
	if idx < 0 || idx > len(c.Models)-1 {
		return nil, fmt.Errorf("index out of range")
	}
	mc := c.Models[idx]

	switch mc.Type {
	case breachModel:
		var cfg breach.Config
		if err := json.Unmarshal(mc.Config, &cfg); err != nil {
			return nil, fmt.Errorf("error on loading config: %s", err)
		}
		return breach.NewModel(cfg), nil
	case storyModel:
		var cfg story.Config
		if err := json.Unmarshal(mc.Config, &cfg); err != nil {
			return nil, fmt.Errorf("error on loading config: %s", err)
		}
		return story.NewModel(cfg), nil
	case endModel:
		var cfg end.Config
		if err := json.Unmarshal(mc.Config, &cfg); err != nil {
			return nil, fmt.Errorf("error on loading config: %s", err)
		}
		return end.NewModel(cfg), nil
	default:
		return nil, fmt.Errorf("model not found for config: %s", mc.Type)
	}
}

// NewGameConfig
func ReadConfigFile(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error on loading config file: %s", err)
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error on unmarshal config data: %s", err)
	}
	return cfg, nil
}
