package style

import "github.com/charmbracelet/lipgloss"

const (
	BrightGold     = lipgloss.Color("#FFD700")
	DarkGray       = lipgloss.Color("#0C0C0C")
	DarkRed        = lipgloss.Color("#8B0000")
	Indigo         = lipgloss.Color("#4B0082")
	LimeGreen      = lipgloss.Color("#00FF00")
	MetallicGold   = lipgloss.Color("#F4A300")
	MetallicSilver = lipgloss.Color("#F0F0F0")
	NeonCyan       = lipgloss.Color("#00FFFF")
	NeonMagenta    = lipgloss.Color("#FF00FF")
	NeonPink       = lipgloss.Color("#FF007F")
	NeonPurple     = lipgloss.Color("#D400FF")
	VividGreen     = lipgloss.Color("#00A300")
)

var RootStyle = lipgloss.NewStyle().Background(DarkGray)
var BoldStyle = lipgloss.NewStyle().Inherit(RootStyle).Bold(true)
