package story

import (
	tea "github.com/charmbracelet/bubbletea"
)

// StoryModel is a model to show text story
type StoryModel struct {
	// text   bytes.Buffer
	output string
	// done   bool
}

// Init initializes the StoryModel.
func (s StoryModel) Init() tea.Cmd {
	return nil
}

func (s StoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return s, nil
}

func (s StoryModel) View() string {
	return s.output
}

// NewStoryModel return a breach model instance
func NewStoryModel() StoryModel {
	return StoryModel{}
}
