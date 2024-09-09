package breach

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/end"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
	"golang.org/x/exp/slices"

	tea "github.com/charmbracelet/bubbletea"
)

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

const SequencesDone end.EndReason = "All sequences are completed"
const TimerDone end.EndReason = "Timer is ended"

// BreachModel is the base model for the game
type BreachModel struct {
	matrix    MatrixModel
	buffer    Buffer
	sequences []Sequence
	timer     timer.Model

	Width  int
	Height int
	keyMap keymap.KeyMap
}

// SetSize resize the window.
func (b *BreachModel) SetSize(msg tea.WindowSizeMsg) {
	b.Height = msg.Height - marginBottom
	b.Width = msg.Width
}

// checkBufferSize read the buffer free size. The game is over when size <= 0
// It raise a buffer too small if sequence size is higher than buffer size.
func (b BreachModel) checkBufferSize(size BufferSizeMsg) (tea.Model, tea.Cmd) {
	if size <= 0 {
		return b.isOver(BufferIsFull)
	}
	cmds := []tea.Cmd{}
	for _, seq := range b.sequences {
		if !seq.IsDone() && seq.Last() > int(size) {
			cmds = append(cmds, OnBufferTooSmallMsg(seq.Id))
		}
	}
	return b, tea.Batch(cmds...)
}

// isOver verify if the game is over or not. To continue, the player should have fullfilled at least one sequence in the round
func (b BreachModel) isOver(reason end.EndReason) (tea.Model, tea.Cmd) {
	isOver := true
	if idx := slices.IndexFunc(b.sequences, func(seq Sequence) bool { return seq.GetStatus() == SequenceSuccess }); idx >= 0 {
		isOver = false
	}
	return b, end.OnEndReasonMsg(reason, isOver)
}

// Init initializes the BreachModel.
func (b BreachModel) Init() tea.Cmd {
	b.buffer.SetCurrentSymbol(b.matrix.GetSymbol())
	return b.timer.Stop()
}

// Update handle messages for BreachModel.
func (b BreachModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Resize window
	case tea.WindowSizeMsg:
		b.SetSize(msg)
	// Handle timer tick update
	case timer.TickMsg, timer.StartStopMsg:
		var cmd tea.Cmd
		b.timer, cmd = b.timer.Update(msg)
		return b, cmd
	// End round on timer timeout
	case timer.TimeoutMsg:
		return b.isOver(TimerDone)
	// Check buffer size on new symbol saved to see if the game is over
	case BufferSizeMsg:
		return b.checkBufferSize(msg)
	case tea.KeyMsg:
		var cmd tea.Cmd
		b.matrix, cmd = b.matrix.Update(msg)
		return b, cmd
	// Check sequences status on symbols saved
	case SequenceStatusMsg:
		if msg.Status == SequenceSuccess || msg.Status == SequenceFailed {
			// Check if any sequence is not Done
			if idx := slices.IndexFunc(b.sequences, func(seq Sequence) bool { return !seq.IsDone() }); idx < 0 {
				reason := SequencesDone
				// If all sequences are done, update end reason if all sequences are completed or not
				if idx := slices.IndexFunc(b.sequences, func(seq Sequence) bool { return seq.GetStatus() == SequenceFailed }); idx >= 0 {
					reason = NotEnoughSpace
				}
				return b.isOver(reason)
			}
			return b, nil
		}
	// Pass all messages not already handled to buffer and sequences
	default:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		// Update buffer and sequences
		b.buffer, cmd = b.buffer.Update(msg)
		cmds = append(cmds, cmd)
		for i, seq := range b.sequences {
			var cmd tea.Cmd
			b.sequences[i], cmd = seq.Update(msg)
			cmds = append(cmds, cmd)
		}
		// Start timer on first symbol selected
		if !b.timer.Running() {
			cmds = append(cmds, b.timer.Start())
		}
		return b, tea.Batch(cmds...)
	}
	return b, nil
}

// View update console on each update
func (b BreachModel) View() string {
	var s strings.Builder

	s.WriteString(b.center(b.timerView()))
	// Workaround to force background black
	body := lipgloss.JoinHorizontal(lipgloss.Center,
		b.matrix.View(),
		lipgloss.Place(lipgloss.Width(b.sequencesView()), lipgloss.Height(b.matrix.View()), lipgloss.Left, lipgloss.Center, b.sequencesView(), lipgloss.WithWhitespaceBackground(style.DarkGray)),
	)
	s.WriteString(b.center(body))
	s.WriteString(b.center(b.buffer.View()))

	return style.RootStyle.Render(s.String())
}

func (b BreachModel) center(content string) string {
	return lipgloss.Place(b.Width, lipgloss.Height(content), lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceBackground(style.DarkGray))
}

// sequencesView return the sequences view
func (b BreachModel) sequencesView() string {
	var s strings.Builder
	for i, seq := range b.sequences {
		s.WriteString(seq.View())
		if i < len(b.sequences)-1 {
			tools.NewLine(&s)
		}
	}
	return style.SpaceBox("Sequences to upload", s.String(), lipgloss.Left)
}

// timerView return the timer view
func (b BreachModel) timerView() string {
	var s strings.Builder
	time := style.RootStyle.Foreground(style.NeonMagenta).Render(fmt.Sprintf("%.4s", b.timer.View()))
	s.WriteString(style.RootStyle.
		Border(lipgloss.NormalBorder()).BorderBackground(style.DarkGray).
		Foreground(style.MetallicGold).
		Padding(0, 1).Render("Breach Time Remaining: " + time))
	return s.String()
}

// NewBreachModel return a breach model instance
func NewBreachModel(cfg Config) BreachModel {
	return BreachModel{
		matrix:    NewMatrix(cfg.Matrix),
		buffer:    NewBuffer(cfg.Buffer),
		sequences: NewSequences(cfg.Sequences),

		timer:  timer.NewWithInterval(30*time.Second, time.Second),
		keyMap: keymap.DefaultKeyMap(),
	}
}
