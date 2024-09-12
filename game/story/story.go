package story

import (
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/message"
)

// Model is a model to show text story
type Model struct {
	text   io.RuneReader
	output []rune
	keyMap keymap.KeyMap
	ticker *time.Ticker
}

type tickMsg string

func (m Model) OnTick() tea.Cmd {
	return func() tea.Msg {
		<-m.ticker.C
		return tickMsg("")
	}
}

// Init initializes the StoryModel.
func (m Model) Init() tea.Cmd {
	return m.OnTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		r, _, err := m.text.ReadRune()
		if err != nil {
			if err == io.EOF {
				m.ticker.Stop()
				return m, nil
			}
		}
		m.output = append(m.output, r)
		return m, m.OnTick()
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Select) {
			return m, message.OnEndViewMsg(message.EndModelMsg{Status: message.Success})
		}
	}
	return m, nil
}

func (m Model) View() string {
	return string(m.output)
}

// NewModel return a breach model instance
func NewModel(cfg Config) Model {
	return Model{
		text:   strings.NewReader(cfg.Text),
		keyMap: keymap.DefaultKeyMap(),
		ticker: time.NewTicker(30 * time.Millisecond),
	}
}
