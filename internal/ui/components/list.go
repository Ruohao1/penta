package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListItem is a single entry in the list.
type ListItem struct {
	Title string
	Desc  string
}

// ListModel is a reusable, scrollable, selectable list.
//
// It does NOT know about Bubble Tea. You drive it from your model:
//   - call MoveUp/MoveDown on key presses
//   - call View() in your View().
type ListModel struct {
	Items    []ListItem
	Selected int // index into Items
	Height   int // max number of rows to display (0 = no limit)

	// Styles
	NormalTitle   lipgloss.Style
	SelectedTitle lipgloss.Style
	DescStyle     lipgloss.Style
	Container     lipgloss.Style

	keys ListKeyMap
}

// NewList constructs a list with sensible default styles.
func NewList(items []ListItem, height int) ListModel {
	normal := lipgloss.NewStyle()
	selected := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	desc := lipgloss.NewStyle().Faint(true)
	container := lipgloss.NewStyle()

	return ListModel{
		Items:         items,
		Selected:      0,
		Height:        height,
		NormalTitle:   normal,
		SelectedTitle: selected,
		DescStyle:     desc,
		Container:     container,
		keys:          NewListKeyMap(),
	}
}

func (l ListModel) Init() tea.Cmd {
	return nil
}

func (l ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, l.keys.Up):
			l.MoveUp()
		case key.Matches(msg, l.keys.Down):
			l.MoveDown()
		case key.Matches(msg, l.keys.Quit):
			return l, tea.Quit
		}
	}
	return l, nil
}

// MoveUp moves the selection up by one item.
func (l *ListModel) MoveUp() {
	if len(l.Items) == 0 {
		return
	}
	if l.Selected > 0 {
		l.Selected--
	}
}

// MoveDown moves the selection down by one item.
func (l *ListModel) MoveDown() {
	if len(l.Items) == 0 {
		return
	}
	if l.Selected < len(l.Items)-1 {
		l.Selected++
	}
}

// View renders the list as a string, applying scrolling if Height > 0.
func (l ListModel) View() string {
	if len(l.Items) == 0 {
		return l.Container.Render("(no items)")
	}

	start, end := l.visibleRange()

	var lines []string
	for i := start; i < end; i++ {
		item := l.Items[i]

		titleStyle := l.NormalTitle
		if i == l.Selected {
			titleStyle = l.SelectedTitle
		}

		line := titleStyle.Render(item.Title)
		if item.Desc != "" {
			line = line + " " + l.DescStyle.Render(item.Desc)
		}
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		l.Container.Render(strings.Join(lines, "\n")),
		"",
		NewHelpModel(l.keys).View())
}

// visibleRange computes the [start,end) indices to render given Height and Selected.
func (l ListModel) visibleRange() (start, end int) {
	n := len(l.Items)
	if l.Height <= 0 || l.Height >= n {
		// No scrolling needed.
		return 0, n
	}

	// Ensure selected is always visible.
	start = l.Selected - l.Height/2
	if start < 0 {
		start = 0
	}
	end = start + l.Height
	if end > n {
		end = n
		start = end - l.Height
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

type ListKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func NewListKeyMap() ListKeyMap {
	return ListKeyMap{
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
	}
}

func (k ListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k ListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Quit},
	}
}
