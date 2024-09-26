package message

import tea "github.com/charmbracelet/bubbletea"

type EndViewStatus int

const (
	Success EndViewStatus = iota
	Failed
	Error
)

type EndModelMsg struct {
	Id     int           // Id of sender view
	Status EndViewStatus // End status
	Msg    string        // additional data from sender
}

func OnEndViewMsg(msg EndModelMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
