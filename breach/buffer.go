package breach

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const BufferIsFull EndReason = "Buffer is full"
const NotEnoughSpace EndReason = "Not enough space to complete sequence"

type BufferSizeMsg int

func OnBufferSizeMsg(size int) tea.Cmd {
	return func() tea.Msg {
		return BufferSizeMsg(size)
	}
}

type Buffer struct {
	data   []Symbol
	x      int
	isFull bool
	style  BufferStyle
}

func (b Buffer) Last() int { return len(b.data) - b.x }

func (b *Buffer) SetCurrentSymbol(sym Symbol) {
	b.data[b.x] = sym
}

func (b Buffer) UseBufferBlock() (Buffer, tea.Cmd) {
	if b.x <= len(b.data) {
		b.x++
	} else {
		b.isFull = true
	}
	return b, OnBufferSizeMsg(b.Last())
}

func (b Buffer) Init() tea.Cmd { return nil }

func (b Buffer) Update(msg tea.Msg) (Buffer, tea.Cmd) {
	switch msg := msg.(type) {
	case SymbolMsg:
		if b.isFull {
			return b, nil
		}
		b.SetCurrentSymbol(msg.symbol)
		if msg.selected {
			return b.UseBufferBlock()
		}
	}
	return b, nil
}

func (b Buffer) View() string {
	var buf strings.Builder
	for i, sym := range b.data {
		msg := sym.String()
		style := b.style.Selected
		if i == b.x {
			style = b.style.Current
		} else if i > b.x {
			msg = "  "
		}
		buf.WriteString(b.style.Selected.Render("["))
		buf.WriteString(style.Render(msg))
		buf.WriteString(b.style.Selected.Render("]"))
	}

	return SpaceBox(RootStyle.Foreground(BrightGold).Padding(0, 1).Render("Buffer"), RootStyle.Padding(0, 0).Render(buf.String()), lipgloss.Center)
}

type BufferStyle struct {
	Current  lipgloss.Style
	Selected lipgloss.Style
}

func NewBuffer(size int) Buffer {
	return Buffer{
		data:   make([]Symbol, size),
		x:      0,
		isFull: false,
		style: BufferStyle{
			Current:  RootStyle.Foreground(NeonPink).Bold(true),
			Selected: RootStyle.Foreground(LimeGreen),
		},
	}
}
