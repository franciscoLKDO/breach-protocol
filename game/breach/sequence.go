package breach

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/style"
)

const seqMax = 10

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
	Id          int
	data        []Symbol
	x           int
	status      SequenceStatus
	description string
	points      int
	style       SequenceStyle
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
	alignDesc := style.RootStyle.Render(strings.Repeat("   ", seqMax-len(s.data)))
	if s.status == SequenceRunning {
		for i, sym := range s.data {
			if i < s.x {
				res.WriteString(s.style.ValidatedSymbol.Render(sym.String()))
			} else if i == s.x {
				res.WriteString(s.style.CurrentSymbol.Render(sym.String()))
			} else {
				res.WriteString(s.style.NextSymbol.Render(sym.String()))
			}
			res.WriteString(style.RootStyle.Render(" "))
		}
		res.WriteString(alignDesc + style.RootStyle.Render(s.description))
	} else {
		style := s.style.Success
		if s.status == SequenceFailed {
			style = s.style.Failed
		}
		for _, sym := range s.data {
			res.WriteString(style.Render(sym.String() + " "))
		}
		res.WriteString(alignDesc + style.Render(s.description))
	}
	return res.String()
}

type SequenceStyle struct {
	CurrentSymbol   lipgloss.Style
	ValidatedSymbol lipgloss.Style
	NextSymbol      lipgloss.Style
	Failed          lipgloss.Style
	Success         lipgloss.Style
}

func NewSequence(cfg SequenceConfig, id int) Sequence {
	return Sequence{
		Id:          id,
		data:        newSymbols(cfg.Size),
		x:           0,
		status:      SequenceRunning,
		description: cfg.Description,
		style: SequenceStyle{
			CurrentSymbol:   style.RootStyle.Foreground(style.NeonPink).Bold(true),
			ValidatedSymbol: style.RootStyle.Foreground(style.LimeGreen),
			NextSymbol:      style.RootStyle.Foreground(style.NeonCyan),
			Failed:          style.RootStyle.Foreground(style.DarkRed).Bold(true),
			Success:         style.RootStyle.Foreground(style.VividGreen).Bold(true),
		},
	}
}

func NewSequences(cfg []SequenceConfig) []Sequence {
	res := make([]Sequence, len(cfg))
	for i, seq := range cfg {
		res[i] = NewSequence(seq, i)
	}
	return res
}
