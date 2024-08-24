package breach

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines key bindings for each user action.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Quit   key.Binding
}

// DefaultKeyMap defines the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:     key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "up")),
		Down:   key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "down")),
		Left:   key.NewBinding(key.WithKeys("backspace", "left", "esc"), key.WithHelp("h", "left")),
		Right:  key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "right")),
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("a", "select")),
		Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	}
}
