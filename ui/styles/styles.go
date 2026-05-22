package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Apple Dark Theme Colors (No explicit Background for transparency)
	FgColor      = lipgloss.Color("#f5f5f7") // Pale Apple Gray
	AccentBlue   = lipgloss.Color("#2997ff") // High-Luminance Link Blue
	BorderColor  = lipgloss.Color("#424245") // Utility Dark Gray
	Secondary    = lipgloss.Color("#86868b") // Mid Border Gray
	Black        = lipgloss.Color("#000000")

	// Window Style
	WindowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			Foreground(FgColor)

	// Modal Style
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentBlue).
			Padding(1, 4)

	Title = lipgloss.NewStyle().
		Foreground(FgColor).
		Bold(true).
		MarginBottom(1).
		Padding(0, 1).
		Background(AccentBlue)

	Button = lipgloss.NewStyle().
		Foreground(Black).
		Background(AccentBlue).
		Padding(0, 3).
		MarginRight(2).
		Bold(true).
		Render

	ButtonInactive = lipgloss.NewStyle().
			Foreground(FgColor).
			Background(BorderColor).
			Padding(0, 3).
			MarginRight(2).
			Render

	LogStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Height(8)

	ShortcutHint = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)
)
