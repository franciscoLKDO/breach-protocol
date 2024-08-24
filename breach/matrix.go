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
}

func (m MatrixModel) GetData() [][]Symbol { return m.data }

func (m MatrixModel) GetSymbol() Symbol { return m.data[m.y][m.x] }

func (m MatrixModel) SetSymbol(s Symbol) { m.data[m.y][m.x] = s }

func (m MatrixModel) GetCoordonates() (int, int) { return m.x, m.y }

func (m *MatrixModel) SetX(x int) {
	if x >= 0 && x < len(m.data) {
		m.x = x
	}
}

func (m *MatrixModel) SetY(y int) {
	if y >= 0 && y < len(m.data) {
		m.y = y
	}
}

func (m *MatrixModel) setKeymapX() {
	m.keyMap.Left.SetEnabled(true)
	m.keyMap.Right.SetEnabled(true)
	m.keyMap.Down.SetEnabled(false)
	m.keyMap.Up.SetEnabled(false)
}

func (m *MatrixModel) setKeymapY() {
	m.keyMap.Left.SetEnabled(false)
	m.keyMap.Right.SetEnabled(false)
	m.keyMap.Down.SetEnabled(true)
	m.keyMap.Up.SetEnabled(true)
}

func (m *MatrixModel) rotateAxe() {
	m.axe = 1 - m.axe
	if m.axe == X {
		m.setKeymapX()
	} else {
		m.setKeymapY()
	}
}

// applySymbol on enter cmd, selected symbol should not be actual symbol  (if actual exist)
func (m *MatrixModel) applySymbol() tea.Cmd {
	// If coordonates are already choosen, pass
	sym := m.GetSymbol()
	if sym == XXX {
		return nil
	}
	m.SetSymbol(XXX)
	m.rotateAxe()
	return OnSymbol(sym, true)
}

func (m MatrixModel) Init() tea.Cmd { return nil }

func (m MatrixModel) Update(msg tea.Msg) (MatrixModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Right):
			m.SetX(m.x + 1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Left):
			m.SetX(m.x - 1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Up):
			m.SetY(m.y - 1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Down):
			m.SetY(m.y + 1)
			return m, OnSymbol(m.GetSymbol(), false)
		case key.Matches(msg, m.keyMap.Select):
			return m, m.applySymbol()
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
			switch true {
			case sym == XXX:
				s.WriteString(defaultStyle.CurrentSymbol.Render(sym.String()))
			case j == x && i == y:
				s.WriteString(defaultStyle.CurrentSymbol.Render(sym.String()))
			case j == x && m.axe == Y:
				s.WriteString(defaultStyle.CurrentAxe.Render(sym.String()))
			case i == y && m.axe == X:
				s.WriteString(defaultStyle.CurrentAxe.Render(sym.String()))
			default:
				s.WriteString(defaultStyle.InactiveSymbol.Render(sym.String()))
			}
			s.WriteString(" ")
		}
		if i < len(data)-1 {
			s.WriteString("\n")
		}
	}
	return newBox("code Matrix", s.String(), lipgloss.Center)
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
	}
	matrix.setKeymapX()
	return matrix
}
