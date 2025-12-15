package views

import tea "github.com/charmbracelet/bubbletea"

type View int

const (
	HomeView View = iota
	ScanView
	ConsoleView
)

type SwitchViewMsg struct{ View View }

func SwitchViewCmd(v View) tea.Cmd {
	return func() tea.Msg { return SwitchViewMsg{View: v} }
}
