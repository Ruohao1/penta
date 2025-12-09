package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ListItem is a single entry in the list.
type ListItem struct {
	Title string
	Desc  string
}

// List is a reusable, scrollable, selectable list.
//
// It does NOT know about Bubble Tea. You drive it from your model:
//   - call MoveUp/MoveDown on key presses
//   - call View() in your View().
type List struct {
	Items    []ListItem
	Selected int // index into Items
	Height   int // max number of rows to display (0 = no limit)

	// Styles
	NormalTitle   lipgloss.Style
	SelectedTitle lipgloss.Style
	DescStyle     lipgloss.Style
	Container     lipgloss.Style
}

// NewList constructs a list with sensible default styles.
func NewList(items []ListItem, height int) List {
	normal := lipgloss.NewStyle()
	selected := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	desc := lipgloss.NewStyle().Faint(true)
	container := lipgloss.NewStyle()

	return List{
		Items:         items,
		Selected:      0,
		Height:        height,
		NormalTitle:   normal,
		SelectedTitle: selected,
		DescStyle:     desc,
		Container:     container,
	}
}

// MoveUp moves the selection up by one item.
func (l *List) MoveUp() {
	if len(l.Items) == 0 {
		return
	}
	if l.Selected > 0 {
		l.Selected--
	}
}

// MoveDown moves the selection down by one item.
func (l *List) MoveDown() {
	if len(l.Items) == 0 {
		return
	}
	if l.Selected < len(l.Items)-1 {
		l.Selected++
	}
}

// View renders the list as a string, applying scrolling if Height > 0.
func (l List) View() string {
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

	return l.Container.Render(strings.Join(lines, "\n"))
}

// visibleRange computes the [start,end) indices to render given Height and Selected.
func (l List) visibleRange() (start, end int) {
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
