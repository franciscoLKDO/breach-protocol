package breach

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SpaceBox(title string, content string, align lipgloss.Position) string {
	var s strings.Builder
	// Set titleBorder
	titleBorder := lipgloss.NormalBorder()
	titleBorder.Top = "═"
	titleBorder.Right = "║"
	titleBorder.TopRight = "╗"
	titleBorder.TopLeft = "╭"
	titleBorder.BottomLeft = "├"
	titleBorder.BottomRight = "║"
	// Set title box
	titleStyle := lipgloss.NewStyle().Background(lipgloss.Color("#fff700")).BorderStyle(titleBorder).Padding(0, 10, 0, 0)
	titleBox := titleStyle.Render(title)
	// Set contentBorder
	contentBorder := lipgloss.NormalBorder()
	contentBorder.Right = "║"
	contentBorder.BottomRight = "╯"
	// Set content box
	contentStyle := lipgloss.NewStyle().Border(contentBorder).Align(align).Padding(0, 0).UnsetBorderTop()
	contentBox := contentStyle.Render(content)

	// Align title and content boxes
	if lipgloss.Width(contentBox) > lipgloss.Width(titleBox) {
		titleBox = titleStyle.Width(lipgloss.Width(contentBox) - contentStyle.GetHorizontalFrameSize()).Render(title)
	}
	s.WriteString(titleBox)
	s.WriteString("\n")
	s.WriteString(contentStyle.Width(lipgloss.Width(titleBox) - contentStyle.GetHorizontalFrameSize()).Render(content))

	return lipgloss.NewStyle().Padding(1, 2, 1, 2).Render(s.String())
}
