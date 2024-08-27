package breach

import "github.com/charmbracelet/lipgloss"

// Styles defines the possible customizations for styles in the file picker.
type Styles struct {
	ValidatedSymbol lipgloss.Style
	CurrentSymbol   lipgloss.Style
	FailedSymbol    lipgloss.Style
	InactiveSymbol  lipgloss.Style
	CurrentAxe      lipgloss.Style
}

// DefaultStyles defines the default styling for the file picker.
func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

// DefaultStylesWithRenderer defines the default styling for the file picker,
// with a given Lip Gloss renderer.
func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		ValidatedSymbol: r.NewStyle().Foreground(lipgloss.Color("#05832d")),
		CurrentSymbol:   r.NewStyle().Foreground(lipgloss.Color("#00ff00")).Bold(true),
		FailedSymbol:    r.NewStyle().Background(lipgloss.Color("#780606")).Bold(true),
		InactiveSymbol:  r.NewStyle().Foreground(lipgloss.Color("#008bc1")),
		CurrentAxe:      r.NewStyle().Foreground(lipgloss.Color("#666700")),
	}
}

var defaultStyle Styles = DefaultStyles()
