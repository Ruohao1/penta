package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	"github.com/Ruohao1/penta/internal/ui/components"
)

type HomeViewModel struct {
	Version string

	Menu components.List
	// Add later:
	// LastScanSummary *ScanSummary
	// Notifications   []string
}

var HomeKeyMap = components.KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move menu up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move menu down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select menu item"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	ToggleHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func NewHomeView() HomeViewModel {
	menuItems := []components.ListItem{
		{Title: "Start New Scan", Desc: "[1]"},
		{Title: "Load Targets from File", Desc: "[2]"},
		{Title: "Interactive mode", Desc: "[i]"},
		{Title: "View Last Report", Desc: "[r]"},
		{Title: "Quit", Desc: "[q]"},
	}

	menu := components.NewList(menuItems, 5)
	menu.Selected = 0

	return HomeViewModel{
		Version: "0.0.1",
		Menu:    menu,
	}
}

func (vm HomeViewModel) Render() string {
	banner := components.NewBanner(vm.Version).Render(components.RenderContext{})

	// Placeholder body; you can add menu + last scan summary here later.
	body := lipgloss.NewStyle().
		MarginTop(1).
		Render(vm.Menu.View())

	// Layout: banner at top, body in middle, cheat sheet at bottom.
	return lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		"",
		body,
		"",
	)
}
