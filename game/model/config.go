package model

import (
	"encoding/json"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/breach-protocol/game/model/breach"
	"github.com/franciscolkdo/breach-protocol/game/model/end"
	"github.com/franciscolkdo/breach-protocol/game/model/story"
)

type model string

const (
	breachModel model = "breach"
	storyModel  model = "story"
	endModel    model = "end"
)

var ErrLoadingConfig = errors.New("error on loading config")

type Config struct {
	Type   model           `json:"type"`
	Config json.RawMessage `json:"config"`
}

func (m Config) Load() (tea.Model, error) {
	switch m.Type {
	case breachModel:
		var cfg breach.Config
		if err := json.Unmarshal(m.Config, &cfg); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrLoadingConfig, err)
		}
		return breach.NewModel(cfg), nil
	case storyModel:
		var cfg story.Config
		if err := json.Unmarshal(m.Config, &cfg); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrLoadingConfig, err)
		}
		return story.NewModel(cfg), nil
	case endModel:
		var cfg end.Config
		if err := json.Unmarshal(m.Config, &cfg); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrLoadingConfig, err)
		}
		return end.NewModel(cfg), nil
	default:
		return nil, fmt.Errorf("model not found for config: %s", m.Type)
	}
}
