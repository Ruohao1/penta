package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TableColumn struct {
	Title string
	Width int
}

type TableRow struct {
	Cells []string
}

type Table struct {
	Columns  []TableColumn
	Rows     []TableRow
	Selected int
}

func (t Table) Render(
	headerStyle lipgloss.Style,
	rowStyle lipgloss.Style,
	selectedRowStyle lipgloss.Style,
) string {
	var out []string

	var headerCells []string
	for _, c := range t.Columns {
		headerCells = append(headerCells, fmt.Sprintf("%-*s", c.Width, c.Title))
	}
	out = append(out, headerStyle.Render(strings.Join(headerCells, " ")))

	for i, row := range t.Rows {
		var cells []string
		for j, c := range row.Cells {
			width := t.Columns[j].Width
			cells = append(cells, fmt.Sprintf("%-*s", width, c))
		}
		line := strings.Join(cells, " ")

		style := rowStyle
		if i == t.Selected {
			style = selectedRowStyle
		}
		out = append(out, style.Render(line))
	}

	return strings.Join(out, "\n")
}
