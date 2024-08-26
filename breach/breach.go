package breach

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	Intro View = iota
	Breach
	EndRound
)

const marginBottom = 5
const AppName = "Breach Protocol"

const SequencesDone EndReason = "All sequences are completed"
const TimerDone EndReason = "Timer is ended"

type BreachModel struct {
	matrix        MatrixModel
	buffer        Buffer
	sequences     []Sequence
	sequencesDone int
	endRound      EndRoundModel
	askQuit       bool
	currentView   View

	timer timer.Model

	Width    int
	Height   int
	KeyMap   KeyMap
	quitting bool
}

func (b *BreachModel) SetSize(msg tea.WindowSizeMsg) {
	b.Height = msg.Height - marginBottom
	b.Width = msg.Width
}

func (b *BreachModel) NewGame() {
	b.matrix = NewMatrix(5)
	b.sequences = NewSequences([]int{3, 5, 6})
	b.buffer = NewBuffer(10)
	b.timer.Timeout = 30 * time.Second
}

// Init initializes the file picker model.
func (b BreachModel) Init() tea.Cmd {
	b.buffer.SetCurrentSymbol(b.matrix.GetSymbol())
	return tea.Batch(tea.SetWindowTitle(AppName), b.timer.Init())
}

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

func (b BreachModel) checkBufferSize(size BufferSizeMsg) (tea.Model, tea.Cmd) {
	if size <= 0 {
		return b.isOver(BufferIsFull)
	}
	for _, seq := range b.sequences {
		if !seq.IsDone() && seq.Last() <= int(size) {
			return b, nil
		}
	}
	return b.isOver(NotEnoughSpace)
}

func (b BreachModel) isOver(reason EndReason) (tea.Model, tea.Cmd) {
	isOver := true
	if b.sequencesDone > 0 {
		isOver = false
	}
	return b, OnEndReasonMsg(reason, isOver)
}

func (b BreachModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg)
	case tea.QuitMsg:
		b.askQuit = true
		return b, tea.Quit
	case timer.TickMsg, timer.StartStopMsg:
		var cmd tea.Cmd
		b.timer, cmd = b.timer.Update(msg)
		return b, cmd
	case tea.KeyMsg:
		return b.updateKeysMsg(msg)
	case EndReasondMsg:
		var cmd tea.Cmd
		b.endRound, cmd = b.endRound.Update(msg)
		b.currentView = EndRound
		return b, tea.Batch(b.timer.Stop(), cmd)
	case EndGameMsg:
		if msg == Quit {
			return b, tea.Quit
		}
		b.NewGame()
		b.currentView = Breach
		return b, b.timer.Init()
	case timer.TimeoutMsg:
		return b.isOver(TimerDone)
	case BufferSizeMsg:
		return b.checkBufferSize(msg)
	case SequenceDoneMsg:
		b.sequencesDone++
		if b.sequencesDone == len(b.sequences) {
			return b.isOver(SequencesDone)
		}
		return b, nil
	default:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		b.buffer, cmd = b.buffer.Update(msg)
		cmds = append(cmds, cmd)
		for i, seq := range b.sequences {
			var cmd tea.Cmd
			b.sequences[i], cmd = seq.Update(msg)
			cmds = append(cmds, cmd)
		}
		return b, tea.Batch(cmds...)
	}
	return b, nil
}

func (b BreachModel) View() string {
	var s strings.Builder
	s.WriteString(b.headerView())
	s.WriteString("\n")

	if b.currentView == EndRound {
		s.WriteString(b.endRound.View())
		return s.String()
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, b.timerView(), b.buffer.View()))
	s.WriteString("\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, b.matrix.View(), b.sequencesView()))
	s.WriteString(b.FooterView())
	return s.String()
}

func (b BreachModel) headerView() string {
	border := lipgloss.DoubleBorder()
	border.Right = "╠"
	title := lipgloss.NewStyle().Background(lipgloss.Color("#fff700")).Foreground(lipgloss.Color("#030fff")).Bold(true).BorderStyle(border).Padding(0, 1).Render("Breach Protocol")
	line := strings.Repeat("─", max(0, b.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (b BreachModel) FooterView() string {
	border := lipgloss.DoubleBorder()
	border.Left = "╣"
	foot := lipgloss.NewStyle().Bold(true).BorderStyle(border).Padding(0, 1).Render("Bartmoss Team")
	line := strings.Repeat("─", max(0, b.Width-lipgloss.Width(foot)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, foot)
}

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

func (b BreachModel) timerView() string {
	var s strings.Builder
	s.WriteString(lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#fff700")).
		Padding(0, 1).Render(fmt.Sprintf("BREACH TIME REMAINING: %s ", b.timer.View())))
	return s.String()
}

func NewBreachModel() BreachModel {
	return BreachModel{
		matrix:        NewMatrix(5),
		buffer:        NewBuffer(10),
		sequences:     NewSequences([]int{3, 5, 6}),
		sequencesDone: 0,
		endRound:      NewEndRound(),
		askQuit:       false,
		quitting:      false,
		currentView:   Breach,

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
