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

const marginBottom = 5
const AppName = "Breach Protocol"

type BreachModel struct {
	matrix    MatrixModel
	sequences []Sequence
	buffer    Buffer

	timer timer.Model

	Width    int
	Height   int
	KeyMap   KeyMap
	quitting bool
	err      error
}

func (b *BreachModel) SetSize(msg tea.WindowSizeMsg) {
	b.Height = msg.Height - marginBottom
	b.Width = msg.Width
}

// Init initializes the file picker model.
func (b BreachModel) Init() tea.Cmd {
	b.buffer.SetCurrentSymbol(b.matrix.GetSymbol())
	return tea.Batch(tea.SetWindowTitle(AppName), b.timer.Init())
}

func (b BreachModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.QuitMsg:
		b.quitting = true
		return b, tea.Quit
	case timer.TickMsg:
		var cmd tea.Cmd
		b.timer, cmd = b.timer.Update(msg)
		return b, cmd
	case timer.TimeoutMsg:
		b.quitting = true
		return b, tea.Quit
	case tea.WindowSizeMsg:
		b.SetSize(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.KeyMap.Quit):
			return b, tea.Quit
		default:
			var cmd tea.Cmd
			b.matrix, cmd = b.matrix.Update(msg)
			return b, cmd
		}
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
	if b.quitting {
		if b.err != nil {
			return fmt.Sprint(b.err)
		}
		return "See ya Choum!"
	}
	var s strings.Builder
	s.WriteString(b.headerView())
	s.WriteString("\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, b.timerView(), b.buffer.View()))
	s.WriteString("\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, b.matrix.View(), b.sequencesView()))
	return s.String()
}

func (b BreachModel) headerView() string {
	border := lipgloss.DoubleBorder()
	border.Right = "╠"
	title := lipgloss.NewStyle().Background(lipgloss.Color("#fff700")).Foreground(lipgloss.Color("#030fff")).Bold(true).BorderStyle(border).Padding(0, 1).Render("Breach Protocol")
	line := strings.Repeat("─", max(0, b.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

// Create a box with title and content
func newBox(title string, content string, align lipgloss.Position) string {
	var s strings.Builder
	// Set border
	border := lipgloss.NormalBorder()
	border.Top = "═"
	border.Right = "║"
	border.TopRight = "╗"
	border.TopLeft = "╭"
	border.BottomLeft = "├"
	border.BottomRight = "║"
	// Set title box
	titleBox := lipgloss.NewStyle().Background(lipgloss.Color("#fff700")).BorderStyle(border).Padding(0, 10, 0, 0).Render(title)
	// Set content
	contentStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Align(align).Padding(0, 0).UnsetBorderTop()

	s.WriteString(titleBox)
	s.WriteString("\n")
	s.WriteString(contentStyle.Width(lipgloss.Width(titleBox) - contentStyle.GetHorizontalFrameSize()).Render(content))
	return lipgloss.NewStyle().Padding(1, 2, 1, 2).Render(s.String())
}

func (b BreachModel) sequencesView() string {
	var s strings.Builder
	for i, seq := range b.sequences {
		s.WriteString(seq.View())
		if i < len(b.sequences)-1 {
			s.WriteString("\n")
		}
	}
	return newBox("Sequence required to upload", s.String(), lipgloss.Left)
}

func (b BreachModel) timerView() string {
	var s strings.Builder
	s.WriteString(lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("#fff700")).
		Padding(0, 1).Render(fmt.Sprintf("BREACH TIME REMAINING: %s ", b.timer.View())))
	return s.String()
}

func NewBreachModel() *BreachModel {
	return &BreachModel{
		matrix: NewMatrix(5),
		buffer: NewBuffer(10),

		sequences: NewSequences([]int{3, 5, 6}),
		timer:     timer.NewWithInterval(30*time.Second, time.Millisecond),
		KeyMap:    DefaultKeyMap(),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
