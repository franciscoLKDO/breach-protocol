package breach

import (
	"fmt"
	"strings"
)

type Buffer struct {
	data   []Symbol
	x      int
	isDone bool
}

func (b Buffer) Last() int { return len(b.data) - b.x }

func (b *Buffer) SetCurrentSymbol(sym Symbol) {
	b.data[b.x] = sym
}

func (b *Buffer) AddSymbol(sym Symbol) {
	b.SetCurrentSymbol(sym)
	if b.x >= len(b.data)-1 {
		b.isDone = true
		return
	}
	b.x++
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

func (b Buffer) String() string {
	var s strings.Builder
	for i, sym := range b.data {
		if i <= b.x {
			s.WriteString(defaultStyle.InactiveSymbol.Render(fmt.Sprintf("[%s]", sym)))
		} else {
			s.WriteString(defaultStyle.InactiveSymbol.Render("[  ]"))
		}
	}
	s.WriteString(fmt.Sprintf("last: %d", b.Last()))
	return s.String()

}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		data:   make([]Symbol, size),
		x:      0,
		isDone: false,
	}
}
