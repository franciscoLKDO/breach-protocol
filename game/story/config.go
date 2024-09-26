package story

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
)

type StoryType string

const (
	Text StoryType = "text"
	Chat StoryType = "chat"
)

type Replica struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Config struct {
	Type StoryType `json:"type"`
	Text string    `json:"text"`
	Chat []Replica `json:"chat"`
}

// newReader return a string reader with rendered text
func (c Config) newReader() *strings.Reader {
	var s strings.Builder
	if c.Type == Chat {
		for _, user := range c.Chat {
			s.WriteString(style.BoldStyle.Render(user.Name+": ") + style.RootStyle.AlignHorizontal(lipgloss.Left).Width(100).Render(user.Text))
			tools.NewLine(&s)
			tools.NewLine(&s)
		}
	} else {
		s.WriteString(style.RootStyle.Render(c.Text))
	}

	return strings.NewReader(s.String())
}
