package model

import (
	"encoding/json"
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

type Config struct {
	Type   model           `json:"type"`
	Config json.RawMessage `json:"config"`
}

func (m Config) Load() (tea.Model, error) {
	switch m.Type {
	case breachModel:
		return newModel(breach.NewModel, m.Config)
	case storyModel:
		return newModel(story.NewModel, m.Config)
	case endModel:
		return newModel(end.NewModel, m.Config)
	default:
		return nil, fmt.Errorf("model not found for config: %s", m.Type)
	}
}

func newModel[T any](cb func(T) tea.Model, config json.RawMessage) (tea.Model, error) {
	var cfg T
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, fmt.Errorf("error on loading config: %w", err)
	}
	return cb(cfg), nil
}
