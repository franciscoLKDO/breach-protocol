package breach

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BufferIsFullMsg bool

func BufferIsFull() tea.Cmd {
	return func() tea.Msg {
		return BufferIsFullMsg(true)
	}
}

type Buffer struct {
	data   []Symbol
	x      int
	isDone bool
}

func (b Buffer) Last() int { return len(b.data) - b.x }

func (b *Buffer) SetCurrentSymbol(sym Symbol) {
	b.data[b.x] = sym
}

func (b *Buffer) UseBufferBlock() tea.Cmd {
	if b.x >= len(b.data)-1 {
		b.isDone = true
		return BufferIsFull()
	}
	b.x++
	return nil
}

// CanFillSequences compare length to see if buffer can fill a given sequence
func (b Buffer) CanFillSequences(sequences []*Sequence) bool {
	for _, seq := range sequences {
		if !seq.IsDone() && seq.Last() <= b.Last() {
			return true
		}
	}
	return false
}

func (b Buffer) Init() tea.Cmd { return nil }

func (b Buffer) Update(msg tea.Msg) (Buffer, tea.Cmd) {
	switch msg := msg.(type) {
	case SymbolMsg:
		b.SetCurrentSymbol(msg.symbol)
		if msg.selected {
			return b, b.UseBufferBlock()
		}
	}
	return b, nil
}

func (b Buffer) View() string {
	var buf strings.Builder
	for i, sym := range b.data {
		if i <= b.x {
			buf.WriteString(defaultStyle.InactiveSymbol.Render(fmt.Sprintf("[%s]", sym)))
		} else {
			buf.WriteString(defaultStyle.InactiveSymbol.Render("[  ]"))
		}
	}
	buf.WriteString(fmt.Sprintf("last: %d", b.Last()))
	var s strings.Builder
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#fff700")).Padding(0, 1).Render("Buffer"))
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 0).Render(buf.String()))
	return s.String()
}

func NewBuffer(size int) Buffer {
	return Buffer{
		data:   make([]Symbol, size),
		x:      0,
		isDone: false,
	}
}
