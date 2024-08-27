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
	Points int
}

func OnSequenceStatusMsg(id int, status SequenceStatus, points int) tea.Cmd {
	return func() tea.Msg {
		return SequenceStatusMsg{Id: id, Status: status, Points: points}
	}
}

type Sequence struct {
	Id          int
	data        []Symbol
	x           int
	status      SequenceStatus
	description string
	points      int
}

func (s Sequence) GetPosition() int          { return s.x }
func (s Sequence) GetData() []Symbol         { return s.data }
func (s Sequence) GetStatus() SequenceStatus { return s.status }
func (s Sequence) GetPoints() int            { return s.points }
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
		return OnSequenceStatusMsg(s.Id, SequenceSuccess, s.points)
	}
	return nil
}

func (s Sequence) Init() tea.Cmd { return nil }

func (s Sequence) Update(msg tea.Msg) (Sequence, tea.Cmd) {
	switch msg := msg.(type) {
	case BufferTooSmallMsg:
		if msg.Id == s.Id {
			s.status = SequenceFailed
			return s, OnSequenceStatusMsg(s.Id, SequenceFailed, 0)
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
	if s.status == SequenceRunning {
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
		res.WriteString(s.description)
		return res.String()
	}
	style := defaultStyle.SuccessSequence
	if s.status == SequenceFailed {
		style = defaultStyle.FailedSequence
	}
	for _, sym := range s.data {
		res.WriteString(style.Render(sym.String() + " "))
	}
	res.WriteString(style.Render(s.description))
	return res.String()
}

func NewSequence(cfg SequenceConfig, id int) Sequence {
	return Sequence{
		Id:          id,
		data:        newSymbols(cfg.Size),
		x:           0,
		status:      SequenceRunning,
		description: cfg.Description,
		points:      cfg.Points,
	}
}

func NewSequences(cfg []SequenceConfig) []Sequence {
	res := make([]Sequence, len(cfg))
	for i, seq := range cfg {
		res[i] = NewSequence(seq, i)
	}
	return res
}
