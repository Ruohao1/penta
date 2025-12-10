package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ruohao1/penta/internal/ui/views"
)

type View int

const (
	HomeView View = iota
	ScanView
	ConsoleView
)

type RootModel struct {
	activeView View
	// ...

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
	switch m.activeView {
	case HomeView:
		m.home, cmd = m.home.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	// case ScanView:
	// 	m.scan.Update(msg)
	case ConsoleView:
		m.console.Update(msg)
	}

	return m, nil
}

func (m RootModel) View() string {
	switch m.activeView {

	case HomeView:
		return m.home.View()
	// case ScanView:
	// 	return m.scan.View()
	case ConsoleView:
		return m.console.View()
	default:
		return "Undefined view"
	}
}
