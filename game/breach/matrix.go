package breach

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
)

type Axe int

const (
	X Axe = iota
	Y
)

const matrixTitle = "Code Matrix"

type MatrixModel struct {
	data   [][]Symbol
	x      int
	y      int
	axe    Axe
	keyMap keymap.KeyMap
	style  MatrixStyle
}

func (m MatrixModel) GetSymbol() Symbol { return m.data[m.y][m.x] }

func (m MatrixModel) SetSymbol(s Symbol) { m.data[m.y][m.x] = s }

func (m *MatrixModel) setX(x int) {
	m.x += x
	if m.x < 0 {
		m.x = len(m.data[m.y]) - 1
	} else if m.x >= len(m.data[m.y]) {
		m.x = 0
	}
}

func (m *MatrixModel) setY(y int) {
	m.y += y
	if m.y < 0 {
		m.y = len(m.data) - 1
	} else if m.y >= len(m.data) {
		m.y = 0
	}
}

func (m *MatrixModel) setKeymap() {
	m.keyMap.Left.SetEnabled(m.axe == X)
	m.keyMap.Right.SetEnabled(m.axe == X)
	m.keyMap.Down.SetEnabled(m.axe == Y)
	m.keyMap.Up.SetEnabled(m.axe == Y)
}

func (m *MatrixModel) rotateAxe() {
	m.axe = 1 - m.axe
	m.setKeymap()
}

// applySymbol on enter cmd, selected symbol should not be actual symbol  (if actual exist)
func (m MatrixModel) applySymbol() (MatrixModel, tea.Cmd) {
	// If coordonates are already choosen, pass
	sym := m.GetSymbol()
	if sym == XXX {
		return m, nil
	}
	m.SetSymbol(XXX)
	m.rotateAxe()
	return m, OnSymbol(sym, true)
}

func (m MatrixModel) Init() tea.Cmd { return nil }

func (m MatrixModel) Update(msg tea.Msg) (MatrixModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Right):
			m.setX(1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Left):
			m.setX(-1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Up):
			m.setY(-1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Down):
			m.setY(+1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Select):
			return m.applySymbol()
		}
	}
	return m, nil
}

func (m MatrixModel) View() string {
	var s strings.Builder
	for i, symbols := range m.data {
		for j, sym := range symbols {
			msg := sym.String()
			if sym == XXX {
				msg = "  "
			}
			switch true {
			case j == m.x && i == m.y:
				if sym == XXX {
					msg = "__"
				}
				s.WriteString(m.style.CurrentSymbol.Render(msg))
			case j == m.x && m.axe == Y:
				s.WriteString(m.style.CurrentAxe.Render(msg))
			case i == m.y && m.axe == X:
				s.WriteString(m.style.CurrentAxe.Render(msg))
			default:
				s.WriteString(m.style.InactiveSymbol.Render(msg))
			}
			s.WriteString(style.RootStyle.Render(" "))
		}
		if i < len(m.data)-1 {
			tools.NewLine(&s)
		}
	}
	return style.SpaceBox(matrixTitle, s.String(), lipgloss.Center)
}

type MatrixStyle struct {
	CurrentSymbol  lipgloss.Style
	InactiveSymbol lipgloss.Style
	CurrentAxe     lipgloss.Style
}

func NewMatrix(size int) MatrixModel {
	m := make([][]Symbol, size)
	for i := range m {
		m[i] = newSymbols(size)
	}
	matrix := MatrixModel{
		data: m,
		x:    0,
		y:    0,

		axe:    X,
		keyMap: keymap.DefaultKeyMap(),
		style: MatrixStyle{
			CurrentSymbol:  style.RootStyle.Foreground(style.NeonPink).Bold(true),
			InactiveSymbol: style.RootStyle.Foreground(style.Indigo),
			CurrentAxe:     style.RootStyle.Foreground(style.NeonCyan),
		},
	}
	matrix.setKeymap()
	return matrix
}
