package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Ruohao1/penta/internal/ui/components"
)

type HomeModel struct {
	Menu components.MenuModel
}

func NewHomeModel() HomeModel {
	scanKey := key.NewBinding(
		key.WithKeys("1", "s"),
		key.WithHelp("1/s", "start scan"),
	)
	consoleKey := key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "open console"),
	)

	scanItem := components.NewMenuItem(
		"Start New Scan",
		"Launch scan workflow",
		scanKey,
		nil,
	)

	consoleItem := components.NewMenuItem(
		"Open Console",
		"Interactive commands & logs",
		consoleKey,
		SwitchViewCmd(ConsoleView),
	)

	menu := components.NewMenu([]components.MenuItem{
		scanItem,
		consoleItem,
	}, 5)

	return HomeModel{Menu: menu}
}

func (vm HomeModel) Init() tea.Cmd { return nil }

func (vm HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	var cmd tea.Cmd
	vm.Menu, cmd = vm.Menu.Update(msg)
	return vm, cmd
}

func (vm HomeModel) View() string {
	banner := components.NewBanner().Render(components.RenderContext{})
	body := lipgloss.NewStyle().MarginTop(1).Render(vm.Menu.View())
	return lipgloss.JoinVertical(lipgloss.Center, banner, "", body, "")
}
