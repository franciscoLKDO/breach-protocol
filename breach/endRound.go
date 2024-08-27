//go:generate stringer -type=EndGameMsg -linecomment
package breach

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EndReason string

type EndReasondMsg struct {
	reason EndReason
	isOver bool
}

func OnEndReasonMsg(reason EndReason, isOver bool) tea.Cmd {
	return func() tea.Msg {
		return EndReasondMsg{reason: reason, isOver: isOver}
	}
}

type EndGameMsg int

const (
	Quit EndGameMsg = iota
	Continue
	Restart
)

func OnEndGameMsg(msg EndGameMsg) tea.Cmd {
	return func() tea.Msg {
		return EndGameMsg(msg)
	}
}

type EndRoundModel struct {
	msg           EndReason
	keyMap        KeyMap
	options       []EndGameMsg
	currentOption int
}

func (e *EndRoundModel) setCurrentOption(x int) {
	e.currentOption += x
	if e.currentOption < 0 {
		e.currentOption = len(e.options) - 1
	}
	if e.currentOption >= len(e.options) {
		e.currentOption = 0
	}
}

func (e EndRoundModel) Init() tea.Cmd { return nil }

func (e EndRoundModel) Update(msg tea.Msg) (EndRoundModel, tea.Cmd) {
	switch msg := msg.(type) {
	case EndReasondMsg:
		e.msg = msg.reason
		e.options = []EndGameMsg{Continue, Quit}
		if msg.isOver {
			e.options = []EndGameMsg{Restart, Quit}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keyMap.Right):
			e.setCurrentOption(1)
		case key.Matches(msg, e.keyMap.Left):
			e.setCurrentOption(-1)
		case key.Matches(msg, e.keyMap.Select):
			return e, OnEndGameMsg(e.options[e.currentOption])
		}
	}
	return e, nil
}

func (e EndRoundModel) View() string {
	var s strings.Builder
	s.WriteString("Continue?")
	s.WriteString("\n")
	var opt []string
	for i := 0; i < len(e.options); i++ {
		style := defaultStyle.InactiveSymbol
		if i == e.currentOption {
			style = defaultStyle.CurrentSymbol
		}
		opt = append(opt, style.Border(lipgloss.NormalBorder()).Render(e.options[i].String()))
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, opt...))
	return SpaceBox(string(e.msg), s.String(), lipgloss.Center)
}

func NewEndRound() EndRoundModel {
	return EndRoundModel{
		msg:           "",
		keyMap:        DefaultKeyMap(),
		currentOption: 0,
		options:       []EndGameMsg{Continue, Quit},
	}
}
