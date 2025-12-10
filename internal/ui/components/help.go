package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GeneralHelpKeyMap struct {
	Toggle key.Binding
}

var generalHelpKeyMap = GeneralHelpKeyMap{
	Toggle: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k GeneralHelpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle}
}

func (k GeneralHelpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Toggle},
	}
}

type mergedKeyMap struct {
	maps []help.KeyMap
}

func (m mergedKeyMap) ShortHelp() []key.Binding {
	var out []key.Binding
	for _, km := range m.maps {
		out = append(out, km.ShortHelp()...)
	}
	return out
}

func (m mergedKeyMap) FullHelp() [][]key.Binding {
	var out [][]key.Binding
	for _, km := range m.maps {
		out = append(out, km.FullHelp()...)
	}
	return out
}

func MergeKeyMaps(maps ...help.KeyMap) help.KeyMap {
	return mergedKeyMap{maps: maps}
}

type HelpModel struct {
	keys       help.KeyMap
	toggle     key.Binding
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
}

func NewHelpModel(keys help.KeyMap) HelpModel {
	return HelpModel{
		keys:   MergeKeyMaps(keys, generalHelpKeyMap),
		toggle: generalHelpKeyMap.Toggle,
		help:   help.New(),
		inputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m *HelpModel) SetKeyMap(keys help.KeyMap) {
	m.keys = MergeKeyMaps(keys, generalHelpKeyMap)
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.toggle):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) View() string {
	return m.help.View(m.keys)
}
