package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Title string
	Desc  string

	Key key.Binding // optional hotkey shown in help
	Cmd tea.Cmd     // executed on Enter (or hotkey)
}

func NewMenuItem(title, desc string, k key.Binding, cmd tea.Cmd) MenuItem {
	return MenuItem{Title: title, Desc: desc, Key: k, Cmd: cmd}
}

type GeneralMenuKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func NewGeneralMenuKeyMap() GeneralMenuKeyMap {
	return GeneralMenuKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// Implements help.KeyMap so it can be rendered
func (k GeneralMenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k GeneralMenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Enter}, {k.Quit}}
}

type MenuModel struct {
	Items    []MenuItem
	Selected int
	Height   int

	NormalTitle   lipgloss.Style
	SelectedTitle lipgloss.Style
	DescStyle     lipgloss.Style
	Container     lipgloss.Style

	navKeys GeneralMenuKeyMap // used for key.Matches
	helpKM  help.KeyMap       // used only for help rendering

	hotkeys map[string]tea.Cmd
	help    help.Model
}

func NewMenu(items []MenuItem, height int) MenuModel {
	normal := lipgloss.NewStyle()
	selected := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	desc := lipgloss.NewStyle().Faint(true)
	container := lipgloss.NewStyle()

	nav := NewGeneralMenuKeyMap()

	// Build hotkey -> cmd map from MenuItem.Key
	hk := map[string]tea.Cmd{}
	for _, it := range items {
		if it.Cmd == nil {
			continue
		}
		for _, s := range it.Key.Keys() {
			if s != "" {
				hk[s] = it.Cmd
			}
		}
	}

	// Help keymap = nav keys + all item keys (+ whatever global help keys you merged inside NewHelpModel)
	itemHelp := itemKeyMap(items)         // implements help.KeyMap
	helpKM := MergeKeyMaps(nav, itemHelp) // your existing MergeKeyMaps

	return MenuModel{
		Items:         items,
		Selected:      0,
		Height:        height,
		NormalTitle:   normal,
		SelectedTitle: selected,
		DescStyle:     desc,
		Container:     container,
		navKeys:       nav,
		helpKM:        helpKM,
		hotkeys:       hk,
		help:          help.New(),
	}
}

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		if cmd, ok := m.hotkeys[msg.String()]; ok && cmd != nil {
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.navKeys.Up):
			m.MoveUp()
			return m, nil
		case key.Matches(msg, m.navKeys.Down):
			m.MoveDown()
			return m, nil
		case key.Matches(msg, m.navKeys.Enter):
			if len(m.Items) == 0 {
				return m, nil
			}
			if cmd := m.Items[m.Selected].Cmd; cmd != nil {
				return m, cmd
			}
			return m, nil
		case key.Matches(msg, m.navKeys.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *MenuModel) MoveUp() {
	if m.Selected > 0 {
		m.Selected--
	}
}

func (m *MenuModel) MoveDown() {
	if m.Selected < len(m.Items)-1 {
		m.Selected++
	}
}

func (m MenuModel) View() string {
	if len(m.Items) == 0 {
		return m.Container.Render("(no items)")
	}

	start, end := m.visibleRange()

	var lines []string
	for i := start; i < end; i++ {
		item := m.Items[i]

		titleStyle := m.NormalTitle
		if i == m.Selected {
			titleStyle = m.SelectedTitle
		}

		titleCol := titleStyle.Width(24).PaddingRight(2).Render(item.Title)
		descCol := m.DescStyle.Render(item.Desc)
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, titleCol, descCol))
	}

	body := m.Container.Render(strings.Join(lines, "\n"))
	helpView := NewHelpModel(m.helpKM).View() // ok, but you can store HelpModel if you want

	return lipgloss.JoinVertical(lipgloss.Left, body, "", helpView)
}

func (m MenuModel) visibleRange() (start, end int) {
	n := len(m.Items)
	if m.Height <= 0 || m.Height >= n {
		return 0, n
	}
	start = m.Selected - m.Height/2
	if start < 0 {
		start = 0
	}
	end = start + m.Height
	if end > n {
		end = n
		start = end - m.Height
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

// itemKeyMap makes a help.KeyMap from item bindings (for help display only).
type itemKeyMap []MenuItem

func (k itemKeyMap) ShortHelp() []key.Binding {
	var out []key.Binding
	for _, it := range k {
		if it.Key.Help().Key != "" {
			out = append(out, it.Key)
		}
	}
	return out
}

func (k itemKeyMap) FullHelp() [][]key.Binding {
	// put item keys on one row; tweak if you want multiple rows
	return [][]key.Binding{k.ShortHelp()}
}
