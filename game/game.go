package game

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/message"
	"github.com/franciscolkdo/breach-protocol/game/model"
	"github.com/franciscolkdo/breach-protocol/game/model/end"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
)

type View int

const (
	Story View = iota
	Breach
	EndRound
)

// const marginBottom = 5
const AppName = "Breach Protocol"
const footerName = "Bartmoss Team"

type Model struct {
	models     []model.Config
	currentIdx int
	current    tea.Model

	keyMap   keymap.KeyMap
	ready    bool
	viewport viewport.Model
	askQuit  bool
}

// Init initializes the BreachModel.
func (m Model) Init() tea.Cmd {
	return tea.Sequence(tea.SetWindowTitle(AppName), m.current.Init())
}

func (m *Model) LoadModel() tea.Cmd {
	var err error
	if m.currentIdx > len(m.models)-1 {
		m.current = end.NewModel(end.Config{Msg: "Félicitations tu as réussi!"})
	} else {
		m.current, err = m.models[m.currentIdx].Load()
		if err != nil {
			return tea.Quit
		}
	}
	return m.current.Init()
}

// Update handle messages for BreachModel.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// Resize window
	case tea.WindowSizeMsg:
		titleHeight := lipgloss.Height(m.titleView(""))
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-(3*titleHeight))
			m.viewport.YPosition = titleHeight
			m.ready = true
			m.viewport.YPosition = titleHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - (3 * titleHeight)
		}
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	// Quit msg
	case tea.QuitMsg:
		m.askQuit = true
		cmds = append(cmds, tea.Quit)
	// Handle key strokes and send them to current model
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Quit) {
			cmds = append(cmds, tea.Quit)
		} else {
			m.current, cmd = m.current.Update(msg)
			cmds = append(cmds, cmd)
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	// EndModelMsg return the state of current model, show end game if failed or next one on success
	case message.EndModelMsg:
		if msg.Status == message.Failed {
			m.current = end.NewModel(end.Config{Msg: msg.Msg})
			cmds = append(cmds, m.current.Init())
		} else {
			m.currentIdx++
			cmds = append(cmds, m.LoadModel())
		}
	// EndGame return Restart or Quit, set currentIdx=0 on restart
	case end.EndGameMsg:
		if msg == end.Quit {
			return m, tea.Quit
		} else {
			m.currentIdx = 0
			cmds = append(cmds, m.LoadModel())
		}
	// Pass all messages not already handled (internal msg for current model)
	default:
		m.current, cmd = m.current.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport.SetContent(m.center(m.current.View()))

	return m, tea.Batch(cmds...)
}

// View update console on each update
func (m Model) View() string {
	var s strings.Builder
	// Set Header
	tools.NewLine(&s)
	s.WriteString(m.titleView(AppName))

	// Set current Model view
	tools.NewLine(&s)
	s.WriteString(m.viewport.View())

	// Set Footer
	tools.NewLine(&s)
	s.WriteString(m.titleView(footerName))
	tools.NewLine(&s)
	return style.RootStyle.Render(s.String())
}

func (m Model) center(content string) string {
	return lipgloss.Place(m.viewport.Width, lipgloss.Height(content), lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceBackground(style.DarkGray))
}

// titleView return the header or footer views of breach protocol
func (m Model) titleView(content string) string {
	border := lipgloss.DoubleBorder()
	border.Right = "╠"
	border.Left = "╣"
	title := gameStyle.Title.BorderForeground(style.MetallicGold).Bold(true).BorderStyle(border).Padding(0, 2).Render(content)
	line := gameStyle.Title.Render(strings.Repeat("═", max(0, (m.viewport.Width/2)-(lipgloss.Width(title)/2))))

	// Workaround to force background black after a border rendering
	afterline := lipgloss.Place(m.viewport.Width, lipgloss.Height(title), lipgloss.Left, lipgloss.Center, line, lipgloss.WithWhitespaceBackground(style.DarkGray))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, title, afterline)
}

type GameStyle struct {
	Title lipgloss.Style
}

var gameStyle = GameStyle{
	Title: style.RootStyle.Foreground(style.MetallicGold),
}

// NewGame return a game model instance
func NewGame(models []model.Config) Model {
	g := Model{
		models:     models,
		ready:      false,
		askQuit:    false,
		currentIdx: 0,
		keyMap:     keymap.DefaultKeyMap(),
	}
	_ = g.LoadModel()
	return g
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
