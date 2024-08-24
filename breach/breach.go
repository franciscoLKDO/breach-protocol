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

const marginBottom = 5
const AppName = "Breach Protocol"

type BreachModel struct {
	matrix    *Matrix
	sequences []*Sequence
	buffer    *Buffer
	selected  [][]int
	axe       Axe
	timer     timer.Model

	Width    int
	Height   int
	KeyMap   KeyMap
	quitting bool
	err      error
}

func (b *BreachModel) setKeymapX() {
	b.KeyMap.Left.SetEnabled(true)
	b.KeyMap.Right.SetEnabled(true)
	b.KeyMap.Down.SetEnabled(false)
	b.KeyMap.Up.SetEnabled(false)
}

func (b *BreachModel) setKeymapY() {
	b.KeyMap.Left.SetEnabled(false)
	b.KeyMap.Right.SetEnabled(false)
	b.KeyMap.Down.SetEnabled(true)
	b.KeyMap.Up.SetEnabled(true)
}

func (b *BreachModel) rotateAxe() {
	b.axe = 1 - b.axe
	if b.axe == X {
		b.setKeymapX()
	} else {
		b.setKeymapY()
	}
}

func (b BreachModel) selectedContains(xy []int) bool {
	for _, s := range b.selected {
		if slices.Equal(s, xy) {
			return true
		}
	}
	return false
}

func (b *BreachModel) onMove(cb func(int, int)) {
	cb(b.matrix.GetCoordonates())
	b.buffer.SetCurrentSymbol(b.matrix.GetSymbol())
}

func (b *BreachModel) down() {
	b.onMove(func(_, y int) {
		b.matrix.SetY(y + 1)
	})
}

func (b *BreachModel) up() {
	b.onMove(func(_, y int) {
		b.matrix.SetY(y - 1)
	})
}

func (b *BreachModel) left() {
	b.onMove(func(x, _ int) {
		b.matrix.SetX(x - 1)
	})
}

func (b *BreachModel) right() {
	b.onMove(func(x, _ int) {
		b.matrix.SetX(x + 1)
	})
}

// applySymbol on enter cmd, selected symbol should not be actual symbol  (if actual exist)
func (b *BreachModel) applySymbol() {
	sym := b.matrix.GetSymbol()
	// Check if selected symbol is not the actual symbol
	x, y := b.matrix.GetCoordonates()
	xy := []int{x, y}
	// If coordonates are already choosen, pass
	if b.selectedContains(xy) {
		return
	}

	// if buffer is done or sequences can't be completed, the round is done -> check results
	if b.buffer.isDone || !b.buffer.CanFillSequences(b.sequences) {
		return // TODO goto results
	}
	b.selected = append(b.selected, xy)
	// Test over sequences
	for _, seq := range b.sequences {
		seq.VerifySymbol(sym)
	}
	b.buffer.AddSymbol(sym)
	b.buffer.SetCurrentSymbol(sym)
	b.rotateAxe()
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
		case key.Matches(msg, b.KeyMap.Right):
			b.right()
		case key.Matches(msg, b.KeyMap.Left):
			b.left()
		case key.Matches(msg, b.KeyMap.Up):
			b.up()
		case key.Matches(msg, b.KeyMap.Down):
			b.down()
		case key.Matches(msg, b.KeyMap.Select):
			b.applySymbol()
		}
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
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, b.timerView(), b.bufferView()))
	s.WriteString("\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, b.matrixView(), b.sequencesView()))
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

func (b BreachModel) matrixView() string {
	var s strings.Builder
	data := b.matrix.GetData()
	x, y := b.matrix.GetCoordonates()
	for i, symbols := range data {
		for j, sym := range symbols {
			switch true {
			case b.selectedContains([]int{j, i}):
				s.WriteString(defaultStyle.CurrentSymbol.Render("XX"))
			case j == x && i == y:
				s.WriteString(defaultStyle.CurrentSymbol.Render(sym.String()))
			case j == x && b.axe == Y:
				s.WriteString(defaultStyle.CurrentAxe.Render(sym.String()))
			case i == y && b.axe == X:
				s.WriteString(defaultStyle.CurrentAxe.Render(sym.String()))
			default:
				s.WriteString(defaultStyle.InactiveSymbol.Render(sym.String()))
			}
			s.WriteString(" ")
		}
		if i < len(data)-1 {
			s.WriteString("\n")
		}
	}
	return newBox("code Matrix", s.String(), lipgloss.Center)
}

func (b BreachModel) sequencesView() string {
	var s strings.Builder
	for i, seq := range b.sequences {
		for j, sym := range seq.data {
			if j < seq.x {
				s.WriteString(defaultStyle.ValidatedSymbol.Render(sym.String()))
			} else if j == seq.x {
				s.WriteString(defaultStyle.CurrentSymbol.Render(sym.String()))
			} else {
				s.WriteString(defaultStyle.InactiveSymbol.Render(sym.String()))
			}
			s.WriteString(" ")
		}
		s.WriteString(fmt.Sprintf("last: %d", seq.Last()))
		if i < len(b.sequences)-1 {
			s.WriteString("\n")
		}
	}
	return newBox("Sequence required to upload", s.String(), lipgloss.Left)
}

func (b BreachModel) bufferView() string {
	var s strings.Builder
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#fff700")).Padding(0, 1).Render("Buffer"))
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 0).Render(b.buffer.String()))
	return s.String()
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
	b := &BreachModel{
		matrix:    NewMatrix(5),
		buffer:    NewBuffer(10),
		selected:  nil,
		sequences: []*Sequence{NewSequence(3), NewSequence(5), NewSequence(6)},
		timer:     timer.NewWithInterval(30*time.Second, time.Millisecond),
		KeyMap:    DefaultKeyMap(),
	}
	b.setKeymapX()
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
