package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Header      lipgloss.Style
	TableHeader lipgloss.Style
	TableRow    lipgloss.Style
	TableRowSel lipgloss.Style
	StatusBar   lipgloss.Style
}

var Default = Theme{
	Header: lipgloss.NewStyle().
		Bold(true),

	TableHeader: lipgloss.NewStyle().
		Bold(true).
		Underline(true),

	TableRow: lipgloss.NewStyle(),

	TableRowSel: lipgloss.NewStyle().
		Reverse(true),

	StatusBar: lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("7")).
		Padding(0, 1),
}
