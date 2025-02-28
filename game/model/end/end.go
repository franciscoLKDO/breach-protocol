//go:generate stringer -type=EndGameMsg -linecomment
package end

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
)

var _ tea.Model = Model{}

const title = "Game Over!"

type EndGameMsg int

const (
	Quit EndGameMsg = iota
	Restart
)

func OnEndGameMsg(msg EndGameMsg) tea.Cmd {
	return func() tea.Msg {
		return EndGameMsg(msg)
	}
}

type Model struct {
	msg           string
	keyMap        keymap.KeyMap
	options       []EndGameMsg
	currentOption int
	style         EndGameStyle
}

func (m *Model) setCurrentOption(x int) {
	m.currentOption += x
	if m.currentOption < 0 {
		m.currentOption = len(m.options) - 1
	}
	if m.currentOption >= len(m.options) {
		m.currentOption = 0
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Right):
			m.setCurrentOption(1)
		case key.Matches(msg, m.keyMap.Left):
			m.setCurrentOption(-1)
		case key.Matches(msg, m.keyMap.Select):
			return m, OnEndGameMsg(m.options[m.currentOption])
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	s.WriteString(m.msg)
	tools.NewLine(&s)
	var opt []string
	for i := 0; i < len(m.options); i++ {
		style := m.style.Inactive
		if i == m.currentOption {
			style = m.style.Active
		}
		opt = append(opt, style.Border(lipgloss.NormalBorder()).Render(m.options[i].String()))
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, opt...))
	return style.SpaceBox(title, s.String(), lipgloss.Center)
}

type EndGameStyle struct {
	Inactive lipgloss.Style
	Active   lipgloss.Style
}

func NewModel(cfg Config) tea.Model {
	return Model{
		msg:           cfg.Msg,
		keyMap:        keymap.DefaultKeyMap(),
		currentOption: 0,
		options:       []EndGameMsg{Restart, Quit},
		style: EndGameStyle{
			Inactive: style.RootStyle.Foreground(style.Indigo),
			Active:   style.RootStyle.Foreground(style.NeonPink).Bold(true),
		},
	}
}
