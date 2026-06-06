package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Core npkill-inspired palette
	BrandColor   = lipgloss.Color("#00E676") // Bright Green
	AlertColor   = lipgloss.Color("#FF3D00") // Neon Red
	WarningColor = lipgloss.Color("#FFC400") // Yellow
	SubtleColor  = lipgloss.Color("#757575") // Gray
	White        = lipgloss.Color("#FFFFFF")

	// UI Element Styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(White).
			Background(BrandColor).
			Padding(0, 1).
			MarginBottom(1)

	CategoryTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(BrandColor).
				Underline(true)

	PathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0"))

	SubtleStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)

	BrandStyle = lipgloss.NewStyle().
			Foreground(BrandColor).
			Bold(true)

	// Dynamic Size Styles
	sizeSmall  = lipgloss.NewStyle().Foreground(BrandColor)
	sizeMedium = lipgloss.NewStyle().Foreground(WarningColor)
	sizeLarge  = lipgloss.NewStyle().Foreground(AlertColor).Bold(true)
)

// FormatSize dynamically colors the megabyte output based on how heavy the folder is.
func FormatSize(mb float64) string {
	s := sizeSmall
	if mb > 500 {
		s = sizeLarge
	} else if mb > 100 {
		s = sizeMedium
	}
	return s.Render(fmt.Sprintf("%8.2f MB", mb))
}
