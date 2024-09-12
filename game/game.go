package game

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/end"
	"github.com/franciscolkdo/breach-protocol/game/keymap"
	"github.com/franciscolkdo/breach-protocol/game/message"
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
	cfg        Config
	currentIdx int
	current    tea.Model

	keyMap  keymap.KeyMap
	askQuit bool
	Width   int
	Height  int
}

// Init initializes the BreachModel.
func (m Model) Init() tea.Cmd {
	return tea.Sequence(tea.SetWindowTitle(AppName), m.current.Init())
}

// SetSize resize the window.
func (m *Model) SetSize(msg tea.WindowSizeMsg) {
	m.Height = msg.Height - marginBottom
	m.Width = msg.Width
}

func (m *Model) LoadModel() tea.Cmd {
	var err error
	m.current, err = m.cfg.LoadModel(m.currentIdx)
	if err != nil {
		//TODO handle in case of error
		return tea.Quit
	}
	return m.current.Init()
}

// Update handle messages for BreachModel.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// Resize window
	case tea.WindowSizeMsg:
		m.SetSize(msg)
		m.current, cmd = m.current.Update(msg)
		return m, cmd
	// Quit msg
	case tea.QuitMsg:
		m.askQuit = true
		return m, tea.Quit
	// Handle key strokes and send them to current model
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Quit) {
			cmd = tea.Quit
		} else {
			m.current, cmd = m.current.Update(msg)
		}
		return m, cmd
	// EndModelMsg return the state of current model, show end game if failed or next one on success
	case message.EndModelMsg:
		if msg.Status == message.Failed {
			m.current = end.NewModel(end.Config{Msg: msg.Msg})
			return m, m.current.Init()
		} else {
			m.currentIdx++
			return m, m.LoadModel()
		}
	// EndGame return Restart or Quit, set currentIdx=0 on restart
	case end.EndGameMsg:
		if msg == end.Quit {
			return m, tea.Quit
		} else {
			m.currentIdx = 0
			return m, m.LoadModel()
		}
	// Pass all messages not already handled (internal msg for current model)
	default:
		m.current, cmd = m.current.Update(msg)
		return m, cmd
	}
}

// View update console on each update
func (m Model) View() string {
	var s strings.Builder
	// Set Header
	tools.NewLine(&s)
	s.WriteString(m.titleView(AppName))

	// Set current Model view
	tools.NewLine(&s)
	s.WriteString(m.center(m.current.View()))

	// Set Footer
	tools.NewLine(&s)
	s.WriteString(m.titleView(footerName))
	tools.NewLine(&s)
	return style.RootStyle.Render(s.String())
}

func (m Model) center(content string) string {
	return lipgloss.Place(m.Width, lipgloss.Height(content), lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceBackground(style.DarkGray))
}

// titleView return the header or footer views of breach protocol
func (m Model) titleView(content string) string {
	border := lipgloss.DoubleBorder()
	border.Right = "╠"
	border.Left = "╣"
	title := gameStyle.Title.BorderForeground(style.MetallicGold).Bold(true).BorderStyle(border).Padding(0, 2).Render(content)
	line := gameStyle.Title.Render(strings.Repeat("═", max(0, (m.Width/2)-(lipgloss.Width(title)/2))))

	// Workaround to force background black after a border rendering
	afterline := lipgloss.Place(m.Width, lipgloss.Height(title), lipgloss.Left, lipgloss.Center, line, lipgloss.WithWhitespaceBackground(style.DarkGray))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, title, afterline)
}

type GameStyle struct {
	Title lipgloss.Style
}

var gameStyle = GameStyle{
	Title: style.RootStyle.Foreground(style.MetallicGold),
}

// NewGame return a game model instance
func NewGame(cfg Config) Model {
	g := Model{
		cfg:        cfg,
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
