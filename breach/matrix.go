package breach

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MatrixModel struct {
	data   [][]Symbol
	x      int
	y      int
	axe    Axe
	keyMap KeyMap
	style  MatrixStyle
}

func (m MatrixModel) GetData() [][]Symbol { return m.data }

func (m MatrixModel) GetSymbol() Symbol { return m.data[m.y][m.x] }

func (m MatrixModel) SetSymbol(s Symbol) { m.data[m.y][m.x] = s }

func (m MatrixModel) GetCoordonates() (int, int) { return m.x, m.y }

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
	x := true
	y := false
	if m.axe == Y {
		x = false
		y = true
	}
	m.keyMap.Left.SetEnabled(x)
	m.keyMap.Right.SetEnabled(x)
	m.keyMap.Down.SetEnabled(y)
	m.keyMap.Up.SetEnabled(y)
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
	data := m.GetData()
	x, y := m.GetCoordonates()
	for i, symbols := range data {
		for j, sym := range symbols {
			msg := sym.String()
			if sym == XXX {
				msg = "  "
			}
			switch true {
			case j == x && i == y:
				if sym == XXX {
					msg = "__"
				}
				s.WriteString(m.style.CurrentSymbol.Render(msg))
			case j == x && m.axe == Y:
				s.WriteString(m.style.CurrentAxe.Render(msg))
			case i == y && m.axe == X:
				s.WriteString(m.style.CurrentAxe.Render(msg))
			default:
				s.WriteString(m.style.InactiveSymbol.Render(msg))
			}
			s.WriteString(RootStyle.Render(" "))
		}
		if i < len(data)-1 {
			newLine(&s)
		}
	}
	return SpaceBox("code Matrix", s.String(), lipgloss.Center)
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
		keyMap: DefaultKeyMap(),
		style: MatrixStyle{
			CurrentSymbol:  RootStyle.Foreground(NeonPink).Bold(true),
			InactiveSymbol: RootStyle.Foreground(Indigo),
			CurrentAxe:     RootStyle.Foreground(NeonCyan),
		},
	}
	matrix.setKeymap()
	return matrix
}
