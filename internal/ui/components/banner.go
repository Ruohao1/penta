package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var _ Component = (*Banner)(nil)

type Banner struct {
	title    string
	name     string
	subtitle string
	version  string
}

func NewBanner(version string) Banner {
	return Banner{
		asciiTitle(),
		"Penta",
		"Pentest Automation Engine",
		version,
	}
}

func (b Banner) Render(ctx RenderContext) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("45")).
		Align(lipgloss.Center).
		Render(asciiTitle())

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246")).
		Align(lipgloss.Center)

	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Faint(true).
		Align(lipgloss.Center)

	sub := subtitleStyle.Render(b.subtitle)
	meta := metaStyle.Render(fmt.Sprintf("%s · v%s", b.name, b.version))

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle,
		sub,
		meta,
	)
}

func asciiTitle() string {
	return "    ____             __       \n" +
		"   / __ \\___  ____  / /_____ _\n" +
		"  / /_/ / _ \\/ __ \\/ __/ __ `/\n" +
		" / ____/  __/ / / / /_/ /_/ / \n" +
		"/_/    \\___/_/ /_/\\__/\\__,_/  \n"
}
