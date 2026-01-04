package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ruohao1/penta/internal/ui/views"
)

type RootModel struct {
	activeView views.View

	home    views.HomeModel
	scan    views.ScanModel
	console views.ConsoleModel
}

func NewRootModel() RootModel {
	return RootModel{
		home:    views.NewHomeModel(),
		scan:    views.NewScanModel(),
		console: views.NewConsoleModel(),
	}
}

func RunTUI(ctx context.Context, _ TuiOptions) error {
	m := NewRootModel()
	p := tea.NewProgram(
		m,
		tea.WithContext(ctx),
	)
	_, err := p.Run()
	return err
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// switch msg := msg.(type) {
	// case controller.SwitchViewMsg:
	// 	m.activeView = msg.View
	// 	return m, nil
	// }

	switch m.activeView {
	case views.HomeView:
		m.home, cmd = m.home.Update(msg)
	// case ScanView:
	// 	m.scan.Update(msg)
	case views.ConsoleView:
		m.console, cmd = m.console.Update(msg)
	}

	return m, cmd
}

func (m RootModel) View() string {
	switch m.activeView {

	case views.HomeView:
		return m.home.View()
	// case ScanView:
	// 	return m.scan.View()
	case views.ConsoleView:
		return m.console.View()
	default:
		return "Undefined view"
	}
}
