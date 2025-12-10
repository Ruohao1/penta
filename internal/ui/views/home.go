package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Ruohao1/penta/internal/ui/components"
)

type HomeModel struct {
	Menu components.ListModel
}

func NewHomeModel() HomeModel {
	menuItems := []components.ListItem{
		{Title: "Start New Scan", Desc: "[1]"},
		{Title: "Load Targets from File", Desc: "[2]"},
		{Title: "Interactive mode", Desc: "[i]"},
		{Title: "View Last Report", Desc: "[r]"},
		{Title: "Quit", Desc: "[q]"},
	}

	menu := components.NewList(menuItems, 5)
	menu.Selected = 0

	return HomeModel{
		Menu: menu,
	}
}

func (vm HomeModel) Init() tea.Cmd {
	return nil
}

func (vm HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	var cmd tea.Cmd
	vm.Menu, cmd = vm.Menu.Update(msg)
	return vm, cmd
}

func (vm HomeModel) View() string {
	banner := components.NewBanner().Render(components.RenderContext{})

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
