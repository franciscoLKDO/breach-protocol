package game

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/breach"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
)

type View int

const (
	Story View = iota
	Breach
	EndRound
)

const marginBottom = 5
const AppName = "Breach Protocol"
const footerName = "Bartmoss Team"

type Model struct {
	currentModel tea.Model

	keyMap  keymap.KeyMap
	askQuit bool
	Width   int
	Height  int
}

// Init initializes the BreachModel.
func (g Model) Init() tea.Cmd {
	return tea.Sequence(tea.SetWindowTitle(AppName), g.currentModel.Init())
}

// SetSize resize the window.
func (b *Model) SetSize(msg tea.WindowSizeMsg) {
	b.Height = msg.Height - marginBottom
	b.Width = msg.Width
}

// Update handle messages for BreachModel.
func (g Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Resize window
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		g.SetSize(msg)
		g.currentModel, cmd = g.currentModel.Update(msg)
		return g, cmd
	// Quit msg
	case tea.QuitMsg:
		g.askQuit = true
		return g, tea.Quit
	case tea.KeyMsg:
		var cmd tea.Cmd
		if key.Matches(msg, g.keyMap.Quit) {
			cmd = tea.Quit
		} else {
			g.currentModel, cmd = g.currentModel.Update(msg)
		}
		return g, cmd
	//TODO pass keys
	// EndReasonMsg trigger endRound view
	// case EndReasondMsg:
	// 	var cmd tea.Cmd
	// 	g.endRound, cmd = g.endRound.Update(msg)
	// 	g.currentView = EndRound
	// 	return g, tea.Batch(g.timer.Stop(), cmd)
	// // EndGameMsg quit or restart a new round
	// case EndGameMsg:
	// 	switch msg {
	// 	case Quit:
	// 		return g, tea.Quit
	// 	case Restart:
	// 		g.round = 0
	// 		g.score = 0
	// 	}
	// 	g.NewRound()
	// 	return g, nil
	// End round on timer timeout

	// Pass all messages not already handled to buffer and sequences
	default:
		var cmd tea.Cmd
		g.currentModel, cmd = g.currentModel.Update(msg)
		return g, cmd
	}
	return g, nil
}

// View update console on each update
func (g Model) View() string {
	var s strings.Builder
	// Set Header
	tools.NewLine(&s)
	s.WriteString(g.titleView(AppName))

	// Set Breach or EndRound View based on currentView
	// if g.currentView == EndRound {
	// 	s.WriteString(g.center(g.endRound.View()))
	// } else {
	// 	tools.NewLine(&s)
	// 	s.WriteString(g.center(g.timerView()))
	// 	// Workaround to force background black
	// 	body := lipgloss.JoinHorizontal(lipgloss.Center,
	// 		g.matrix.View(),
	// 		lipgloss.Place(lipgloss.Width(g.sequencesView()), lipgloss.Height(g.matrix.View()), lipgloss.Left, lipgloss.Center, g.sequencesView(), lipgloss.WithWhitespaceBackground(DarkGray)),
	// 	)
	// 	s.WriteString(g.center(body))
	// 	s.WriteString(g.center(g.buffer.View()))
	// }
	s.WriteString(g.currentModel.View())
	// Set Footer
	tools.NewLine(&s)
	s.WriteString(g.titleView(footerName))
	tools.NewLine(&s)
	return style.RootStyle.Render(s.String())
}

// titleView return the header or footer views of breach protocol
func (g Model) titleView(content string) string {
	border := lipgloss.DoubleBorder()
	border.Right = "╠"
	border.Left = "╣"
	title := gameStyle.Title.BorderForeground(style.MetallicGold).Bold(true).BorderStyle(border).Padding(0, 2).Render(content)
	line := gameStyle.Title.Render(strings.Repeat("═", max(0, (g.Width/2)-(lipgloss.Width(title)/2))))

	// Workaround to force background black after a border rendering
	afterline := lipgloss.Place(g.Width, lipgloss.Height(title), lipgloss.Left, lipgloss.Center, line, lipgloss.WithWhitespaceBackground(style.DarkGray))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, title, afterline)
}

type GameStyle struct {
	Title lipgloss.Style
}

var gameStyle = GameStyle{
	Title: style.RootStyle.Foreground(style.MetallicGold),
}

// NewGame return a game model instance
func NewGame() Model {
	return Model{
		askQuit:      false,
		currentModel: breach.NewBreachModel(breach.DefaultConfig),
		keyMap:       keymap.DefaultKeyMap(),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
