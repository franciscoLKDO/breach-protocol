package breach

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"

	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	Intro View = iota
	Breach
	EndRound
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

const SequencesDone EndReason = "All sequences are completed"
const TimerDone EndReason = "Timer is ended"

// BreachModel is the base model for the game
type BreachModel struct {
	matrix      MatrixModel
	buffer      Buffer
	sequences   []Sequence
	endRound    EndRoundModel
	askQuit     bool
	currentView View

	timer timer.Model

	Width    int
	Height   int
	KeyMap   KeyMap
	quitting bool
}

// SetSize resize the window.
func (b *BreachModel) SetSize(msg tea.WindowSizeMsg) {
	b.Height = msg.Height - marginBottom
	b.Width = msg.Width
}

// NewGame reset breach values to restart a new round/game
func (b BreachModel) NewGame() (tea.Model, tea.Cmd) {
	b.matrix = NewMatrix(5)
	b.sequences = NewSequences([]int{3, 5, 6})
	b.buffer = NewBuffer(10)
	b.currentView = Breach
	b.timer.Timeout = 30 * time.Second
	return b, nil
}

// updateKeysMsg reroute keys stokes to current view
func (b BreachModel) updateKeysMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, b.KeyMap.Quit):
		return b, tea.Quit
	default:
		var cmd tea.Cmd
		if b.currentView == EndRound {
			b.endRound, cmd = b.endRound.Update(msg)
		} else {
			b.matrix, cmd = b.matrix.Update(msg)
		}
		return b, cmd
	}
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
func (b BreachModel) isOver(reason EndReason) (tea.Model, tea.Cmd) {
	isOver := true
	if idx := slices.IndexFunc(b.sequences, func(seq Sequence) bool { return seq.GetStatus() == SequenceSuccess }); idx >= 0 {
		isOver = false
	}
	return b, OnEndReasonMsg(reason, isOver)
}

// Init initializes the BreachModel.
func (b BreachModel) Init() tea.Cmd {
	b.buffer.SetCurrentSymbol(b.matrix.GetSymbol())
	return tea.Sequence(tea.SetWindowTitle(AppName), b.timer.Stop())
}

// Update handle messages for BreachModel.
func (b BreachModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Resize window
	case tea.WindowSizeMsg:
		b.SetSize(msg)
	// Quit msg
	case tea.QuitMsg:
		b.askQuit = true
		return b, tea.Quit
	// Handle timer tick update
	case timer.TickMsg, timer.StartStopMsg:
		var cmd tea.Cmd
		b.timer, cmd = b.timer.Update(msg)
		return b, cmd
	// Handle key stokes
	case tea.KeyMsg:
		return b.updateKeysMsg(msg)
	// EndReasonMsg trigger endRound view
	case EndReasondMsg:
		var cmd tea.Cmd
		b.endRound, cmd = b.endRound.Update(msg)
		b.currentView = EndRound
		return b, tea.Batch(b.timer.Stop(), cmd)
	// EndGameMsg quit or restart a new round
	case EndGameMsg:
		if msg == Quit {
			return b, tea.Quit
		}
		return b.NewGame()
	// End round on timer timeout
	case timer.TimeoutMsg:
		return b.isOver(TimerDone)
	// Check buffer size on new symbol saved to see if the game is over
	case BufferSizeMsg:
		return b.checkBufferSize(msg)
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
	// Set Header
	s.WriteString(b.headerView())
	s.WriteString("\n")

	// Set Breach or EndRound View based on currentView
	if b.currentView == EndRound {
		s.WriteString(b.endRound.View())
	} else {
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, b.timerView(), b.buffer.View()))
		s.WriteString("\n")
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, b.matrix.View(), b.sequencesView()))
	}

	// Set Footer
	s.WriteString("\n")
	s.WriteString(b.footerView())
	s.WriteString("\n")
	return s.String()
}

// headerView return the header view of breach protocol
func (b BreachModel) headerView() string {
	border := lipgloss.DoubleBorder()
	border.Right = "╠"
	title := lipgloss.NewStyle().Background(lipgloss.Color("#fff700")).Foreground(lipgloss.Color("#030fff")).Bold(true).BorderStyle(border).Padding(0, 1).Render("Breach Protocol")
	line := strings.Repeat("─", max(0, b.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

// footerView return the footer view of breach protocol
func (b BreachModel) footerView() string {
	border := lipgloss.DoubleBorder()
	border.Left = "╣"
	foot := lipgloss.NewStyle().Bold(true).BorderStyle(border).Padding(0, 1).Render("Bartmoss Team")
	line := strings.Repeat("─", max(0, b.Width-lipgloss.Width(foot)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, foot)
}

// sequencesView return the sequences view
func (b BreachModel) sequencesView() string {
	var s strings.Builder
	for i, seq := range b.sequences {
		s.WriteString(seq.View())
		if i < len(b.sequences)-1 {
			s.WriteString("\n")
		}
	}
	return SpaceBox("Sequence required to upload", s.String(), lipgloss.Left)
}

// timerView return the timer view
func (b BreachModel) timerView() string {
	var s strings.Builder
	s.WriteString(lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#fff700")).
		Padding(0, 1).Render(fmt.Sprintf("BREACH TIME REMAINING: %s ", b.timer.View())))
	return s.String()
}

// NewBreachModel return a breach model instance
func NewBreachModel() BreachModel {
	return BreachModel{
		matrix:      NewMatrix(10), // TODO transform const to var/config
		buffer:      NewBuffer(10),
		sequences:   NewSequences([]int{3, 5, 6}),
		endRound:    NewEndRound(),
		askQuit:     false,
		quitting:    false,
		currentView: Breach,

		timer:  timer.NewWithInterval(30*time.Second, time.Second),
		KeyMap: DefaultKeyMap(),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
