package breach

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SequenceStatus int

const (
	SequenceFailed SequenceStatus = iota
	SequenceSuccess
	SequenceRunning
)

type SequenceStatusMsg struct {
	Id     int
	Status SequenceStatus
}

func OnSequenceStatusMsg(id int, status SequenceStatus) tea.Cmd {
	return func() tea.Msg {
		return SequenceStatusMsg{Id: id, Status: status}
	}
}

type Sequence struct {
	Id     int
	data   []Symbol
	x      int
	status SequenceStatus
}

func (s Sequence) GetPosition() int          { return s.x }
func (s Sequence) GetData() []Symbol         { return s.data }
func (s Sequence) GetStatus() SequenceStatus { return s.status }
func (s Sequence) IsDone() bool              { return s.status < SequenceRunning }
func (s Sequence) Last() int                 { return len(s.data) - s.x }

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
		s.status = SequenceSuccess
		return OnSequenceStatusMsg(s.Id, SequenceSuccess)
	}
	return nil
}

func (s Sequence) Init() tea.Cmd { return nil }

func (s Sequence) Update(msg tea.Msg) (Sequence, tea.Cmd) {
	switch msg := msg.(type) {
	case BufferTooSmallMsg:
		if msg.Id == s.Id {
			s.status = SequenceFailed
			return s, OnSequenceStatusMsg(s.Id, SequenceFailed)
		}
	case SymbolMsg:
		if msg.selected {
			return s, s.VerifySymbol(msg.symbol)
		}
	}
	return s, nil
}

func (s Sequence) View() string {
	var res strings.Builder
	if s.status == SequenceFailed {
		for _, sym := range s.data {
			res.WriteString(defaultStyle.FailedSymbol.Render(sym.String() + " "))
		}
	} else {
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
	}
	return res.String()
}

func NewSequence(size int, id int) Sequence {
	return Sequence{
		Id:     id,
		data:   newSymbols(size),
		x:      0,
		status: SequenceRunning,
	}
}

func NewSequences(sizes []int) []Sequence {
	res := make([]Sequence, len(sizes))
	for i, size := range sizes {
		res[i] = NewSequence(size, i)
	}
	return res
}
