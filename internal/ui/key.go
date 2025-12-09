package ui

import (
	"github.com/Ruohao1/penta/internal/ui/views"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m RootModel) handleHomeKey(msg tea.KeyMsg) (RootModel, tea.Cmd) {
	homeKeyMap := views.HomeKeyMap

	switch {
	case key.Matches(msg, homeKeyMap.Up):
		m.home.Menu.MoveUp()

	case key.Matches(msg, homeKeyMap.Down):
		m.home.Menu.MoveDown()

	case key.Matches(msg, homeKeyMap.Enter):
		switch m.home.Menu.Selected {
		case 0:
			m.activeView = 1
			return m, nil
		case 1:
			m.activeView = 2
			return m, nil
		case 2:
			m.activeView = 3
			return m, nil
		case 3:
			m.activeView = 4
			return m, nil
		case 4:
			return m, tea.Quit
		}

	case key.Matches(msg, homeKeyMap.Quit):
		return m, tea.Quit
	}

	return m, nil
}
