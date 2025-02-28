package story

import (
	"io"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/message"
)

var _ tea.Model = Model{}

// Model is a model to show text story
type Model struct {
	text    io.RuneReader
	output  []rune
	isended bool

	keyMap keymap.KeyMap
	ticker *time.Ticker
}

type tickMsg struct{}

func (m Model) OnTick() tea.Cmd {
	return func() tea.Msg {
		<-m.ticker.C
		return tickMsg{}
	}
}

// Init initializes the StoryModel.
func (m Model) Init() tea.Cmd {
	return m.OnTick()
}

func (m Model) setEndReadingState() tea.Model {
	m.ticker.Stop()
	m.isended = true
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		for { // Don't count white space in writing effet
			r, _, err := m.text.ReadRune()
			if err != nil {
				if err == io.EOF {
					return m.setEndReadingState(), nil
				}
			}
			m.output = append(m.output, r)
			if unicode.IsLetter(r) {
				break
			}
		}
		return m, m.OnTick()
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Select) {
			if !m.isended {
				for { // Read all the remaining text
					r, _, err := m.text.ReadRune()
					if err != nil {
						if err == io.EOF {
							return m.setEndReadingState(), nil
						}
					}
					m.output = append(m.output, r)
				}
			}
			return m, message.OnEndViewMsg(message.EndModelMsg{Status: message.Success})
		}
	}
	return m, nil
}

func (m Model) View() string {
	return string(m.output)
}

// NewModel return a breach model instance
func NewModel(cfg Config) tea.Model {
	return Model{
		text:    cfg.newReader(),
		isended: false,
		keyMap:  keymap.DefaultKeyMap(),
		ticker:  time.NewTicker(30 * time.Millisecond),
	}
}
