package breach

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SequenceDoneMsg struct {
	Id   int
	Done bool
}

func OnSequenceDoneMsg(id int, done bool) tea.Cmd {
	return func() tea.Msg {
		return SequenceDoneMsg{Id: id, Done: done}
	}
}

type Sequence struct {
	Id     int
	data   []Symbol
	x      int
	isDone bool
}

func (s Sequence) GetPosition() int  { return s.x }
func (s Sequence) GetData() []Symbol { return s.data }
func (s Sequence) IsDone() bool      { return s.isDone }
func (s Sequence) Last() int         { return len(s.data) - s.x }

func (s *Sequence) VerifySymbol(sym Symbol) tea.Cmd {
	if s.x >= len(s.data) {
		return nil
	}
	if s.data[s.x] == sym {
		s.x++
	} else {
		s.x = 0
	}
	if s.x >= len(s.data) {
		s.isDone = true
		return OnSequenceDoneMsg(0, true) //TODO CONTINUE HERE
	}
	return nil
}

func (s Sequence) Init() tea.Cmd { return nil }

func (s Sequence) Update(msg tea.Msg) (Sequence, tea.Cmd) {
	switch msg := msg.(type) {
	case SymbolMsg:
		if msg.selected {
			return s, s.VerifySymbol(msg.symbol)
		}
	}
	return s, nil
}

func (s Sequence) View() string {
	var res strings.Builder
	for i, sym := range s.data {
		if i < s.x {
			res.WriteString(defaultStyle.ValidatedSymbol.Render(sym.String()))
		} else if i == s.x {
			res.WriteString(defaultStyle.CurrentSymbol.Render(sym.String()))
		} else {
			res.WriteString(defaultStyle.InactiveSymbol.Render(sym.String()))
		}
		res.WriteString(" ")
	}
	res.WriteString(fmt.Sprintf("last: %d", s.Last())) // debug
	return res.String()
}

func NewSequence(size int, id int) Sequence {
	return Sequence{
		Id:     id,
		data:   newSymbols(size),
		x:      0,
		isDone: false,
	}
}

func NewSequences(sizes []int) []Sequence {
	res := make([]Sequence, len(sizes))
	for i, size := range sizes {
		res[i] = NewSequence(size, i)
	}
	return res
}
