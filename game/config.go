package game

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/breach-protocol/game/breach"
	"github.com/franciscolkdo/breach-protocol/game/end"
	"github.com/franciscolkdo/breach-protocol/game/story"
	"github.com/franciscolkdo/breach-protocol/game/style"
)

type ModelConfig interface{}

type Config struct {
	Models []ModelConfig
}

func (c Config) LoadModel(idx int) (tea.Model, error) {
	if idx < 0 || idx > len(c.Models) {
		return nil, fmt.Errorf("index out of range")
	}
	mc := c.Models[idx]
	switch cfg := mc.(type) {
	case breach.Config:
		return breach.NewModel(cfg), nil
	case story.Config:
		return story.NewModel(cfg), nil
	case end.Config:
		return end.NewModel(cfg), nil
	default:
		return nil, fmt.Errorf("Model not found for config: %s", cfg)
	}
}

var DefaultConfig = Config{
	Models: []ModelConfig{
		story.Config{
			Text: style.RootStyle.Bold(true).Render("Knock Knock NEO,\nFollow the White Rabbit"),
		},
		breach.Config{
			Matrix: 7,
			Buffer: 10,
			Timer:  40 * time.Second,
			Sequences: []breach.SequenceConfig{
				{
					Size:        3,
					Description: "Avoid firewall detection",
					Points:      10,
				},
				{
					Size:        5,
					Description: "Decrypt encrypted files",
					Points:      30,
				},
				{
					Size:        10,
					Description: "Burn Netrunner ice",
					Points:      70,
				},
			},
		},
		breach.Config{
			Matrix: 3,
			Buffer: 5,
			Timer:  40 * time.Second,
			Sequences: []breach.SequenceConfig{
				{
					Size:        3,
					Description: "Lock out intruders",
					Points:      30,
				},
			},
		},
		breach.Config{
			Matrix: 5,
			Buffer: 8,
			Timer:  40 * time.Second,
			Sequences: []breach.SequenceConfig{
				{
					Size:        5,
					Description: "Find Mikoshi source code",
					Points:      30,
				},
				{
					Size:        6,
					Description: "Escape Arasaka Netrunners",
					Points:      50,
				},
			},
		},
	},
}
