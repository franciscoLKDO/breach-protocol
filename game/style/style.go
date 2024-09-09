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

// Styles defines the possible customizations for styles in the file picker.
type Styles struct {
	ValidatedSymbol  lipgloss.Style
	CurrentSymbol    lipgloss.Style
	FailedSequence   lipgloss.Style
	SuccessSequence  lipgloss.Style
	InactiveSymbol   lipgloss.Style
	MatrixCurrentAxe lipgloss.Style
}

// DefaultStyles defines the default styling for the file picker.
func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

// DefaultStylesWithRenderer defines the default styling for the file picker,
// with a given Lip Gloss renderer.
func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		ValidatedSymbol: r.NewStyle().Inherit(RootStyle).Foreground(BrightGold),
		CurrentSymbol:   r.NewStyle().Inherit(RootStyle).Foreground(NeonPink).Bold(true),
		InactiveSymbol:  r.NewStyle().Inherit(RootStyle).Foreground(Indigo),
	}
}

var RootStyle = lipgloss.NewStyle().Background(DarkGray)
