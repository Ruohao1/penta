package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ruohao1/penta/internal/ui/components"
	"github.com/Ruohao1/penta/internal/ui/views"
)

type View int

const (
	HomeView View = iota
	ScanView
)

type RootModel struct {
	activeView View
	// ...

	help components.HelpComponent
	home views.HomeViewModel
}

func NewRootModel() RootModel {
	help := components.NewHelpComponent()
	help.SetKeys(views.HomeKeyMap)
	return RootModel{
		activeView: 0,
		help:       help,
		home:       views.NewHomeView(),
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

func (m RootModel) Init() tea.Cmd { return nil }

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

		// other message types: eventMsg, tickMsg, window size, etc.
	}

	return m, nil
}

func (m RootModel) handleKey(msg tea.KeyMsg) (RootModel, tea.Cmd) {
	key := msg.String()

	// --- global keys, independent of view ---
	switch key {
	case "ctrl+c":
		return m, tea.Quit

	case "?":
		m.help.Toggle()
		return m, nil
	}

	// --- per-view keys ---
	switch m.activeView {
	case HomeView:
		return m.handleHomeKey(msg)

	// case ScanView:
	// 	return m.handleScanKey(msg)

	// add more views as needed
	default:
		return m, nil
	}
}

func (m RootModel) View() string {
	var main string
	switch m.activeView {
	case HomeView:
		main = m.home.Render()
	default:
		return "unknown view"
	}

	helpView := m.help.View()
	if helpView == "" {
		return main
	}

	return main + "\n\n" + helpView
}
