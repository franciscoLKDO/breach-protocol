package breach

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/message"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
	"golang.org/x/exp/slices"

	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = Model{}

type BufferTooSmallMsg struct {
	Id int
}

func OnBufferTooSmallMsg(seqId int) tea.Cmd {
	return func() tea.Msg {
		return BufferTooSmallMsg{seqId}
	}
}

const marginBottom = 5
const AppName = "Breach Protocol"

const (
	SequencesDone = "All sequences are completed"
	TimerDone     = "Timer is ended"
)

// Model is the mini game breach-protocol model
type Model struct {
	id        int
	matrix    MatrixModel
	buffer    Buffer
	sequences []Sequence
	timer     timer.Model

	Width  int
	Height int
	keyMap keymap.KeyMap
}

// SetSize resize the window.
func (m *Model) SetSize(msg tea.WindowSizeMsg) {
	m.Height = msg.Height - marginBottom
	m.Width = msg.Width
}

// checkBufferSize read the buffer free size. The game is over when size <= 0
// It raise a buffer too small if sequence size is higher than buffer size.
func (m Model) checkBufferSize(size BufferSizeMsg) (tea.Model, tea.Cmd) {
	if size <= 0 {
		return m.isOver(BufferIsFull)
	}
	cmds := []tea.Cmd{}
	for _, seq := range m.sequences {
		if !seq.IsDone() && seq.Last() > int(size) {
			cmds = append(cmds, OnBufferTooSmallMsg(seq.Id))
		}
	}
	return m, tea.Batch(cmds...)
}

// isOver verify if the game is over or not. To continue, the player should have fullfilled at least one sequence in the round
func (m Model) isOver(reason string) (tea.Model, tea.Cmd) {
	status := message.Failed
	if idx := slices.IndexFunc(m.sequences, func(seq Sequence) bool { return seq.GetStatus() == SequenceSuccess }); idx >= 0 {
		status = message.Success
	}
	return m, message.OnEndViewMsg(message.EndModelMsg{Id: m.id, Status: status, Msg: reason})
}

// Init initializes the BreachModel.
func (m Model) Init() tea.Cmd {
	m.buffer.SetCurrentSymbol(m.matrix.GetSymbol())
	return m.timer.Init()
}

// Update handle messages for BreachModel.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Resize window
	case tea.WindowSizeMsg:
		m.SetSize(msg)
	// Handle timer tick update
	case timer.TickMsg, timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	// End round on timer timeout
	case timer.TimeoutMsg:
		return m.isOver(TimerDone)
	// Check buffer size on new symbol saved to see if the game is over
	case BufferSizeMsg:
		return m.checkBufferSize(msg)
	case tea.KeyMsg:
		var cmd tea.Cmd
		m.matrix, cmd = m.matrix.Update(msg)
		return m, cmd
	// Check sequences status on symbols saved
	case SequenceStatusMsg:
		if msg.Status == SequenceSuccess || msg.Status == SequenceFailed {
			// Check if any sequence is not Done
			if idx := slices.IndexFunc(m.sequences, func(seq Sequence) bool { return !seq.IsDone() }); idx < 0 {
				reason := SequencesDone
				// If all sequences are done, update end reason if all sequences are completed or not
				if idx := slices.IndexFunc(m.sequences, func(seq Sequence) bool { return seq.GetStatus() == SequenceFailed }); idx >= 0 {
					reason = NotEnoughSpace
				}
				return m.isOver(reason)
			}
			return m, nil
		}
	// Pass all messages not already handled to buffer and sequences
	default:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		// Update buffer and sequences
		m.buffer, cmd = m.buffer.Update(msg)
		cmds = append(cmds, cmd)
		for i, seq := range m.sequences {
			var cmd tea.Cmd
			m.sequences[i], cmd = seq.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

// View update console on each update
func (m Model) View() string {
	var s strings.Builder

	s.WriteString(m.timerView())
	// Workaround to force background black
	body := lipgloss.JoinHorizontal(lipgloss.Center,
		m.matrix.View(),
		lipgloss.Place(lipgloss.Width(m.sequencesView()), lipgloss.Height(m.matrix.View()), lipgloss.Left, lipgloss.Center, m.sequencesView(), lipgloss.WithWhitespaceBackground(style.DarkGray)),
	)
	s.WriteString(body)
	s.WriteString(m.buffer.View())

	return style.RootStyle.Render(s.String())
}

// sequencesView return the sequences view
func (m Model) sequencesView() string {
	var s strings.Builder
	for i, seq := range m.sequences {
		s.WriteString(seq.View())
		if i < len(m.sequences)-1 {
			tools.NewLine(&s)
		}
	}
	return style.SpaceBox("Sequences to upload", s.String(), lipgloss.Left)
}

// timerView return the timer view
func (m Model) timerView() string {
	var s strings.Builder
	time := style.RootStyle.Foreground(style.NeonMagenta).Render(fmt.Sprintf("%.4s", m.timer.View()))
	s.WriteString(style.RootStyle.
		Border(lipgloss.NormalBorder()).BorderBackground(style.DarkGray).
		Foreground(style.MetallicGold).
		Padding(0, 1).Render("Breach Time Remaining: " + time))
	return s.String()
}

// NewModel return a breach model instance
func NewModel(cfg Config) tea.Model {
	return Model{
		matrix:    NewMatrix(cfg.Matrix),
		buffer:    NewBuffer(cfg.Buffer),
		sequences: NewSequences(cfg.Sequences),

		timer:  timer.NewWithInterval(cfg.Timer*time.Second, time.Second),
		keyMap: keymap.DefaultKeyMap(),
	}
}
