package views

import (
	"strings"

	"github.com/Ruohao1/penta/internal/ui/controller"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConsoleModel struct {
	input textinput.Model
	lines []string

	OnCommand func(cmd string) tea.Cmd

	focused bool
	width   int
	height  int
}

func NewConsoleModel() ConsoleModel {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.Focus()

	return ConsoleModel{
		input:   ti,
		focused: true,
		width:   0,
		height:  0,
	}
}

func (m ConsoleModel) Init() tea.Cmd {
	if m.focused {
		return textinput.Blink
	}
	return nil
}

func (m *ConsoleModel) SetSize(w, h int) {
	m.width, m.height = w, h
}

func (m *ConsoleModel) AppendLine(s string) {
	m.lines = append(m.lines, s)
}

func (m ConsoleModel) Update(msg tea.Msg) (ConsoleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		// ESC: go back (don’t kill the whole program)
		case "esc":
			return m, controller.SwitchViewCmd(controller.HomeView)

		// Ctrl+C: quit program (global behavior)
		case "ctrl+c":
			return m, tea.Quit

		// Enter: submit command
		case "enter":
			raw := strings.TrimSpace(m.input.Value())
			if raw == "" {
				m.input.SetValue("")
				return m, nil
			}

			m.AppendLine(m.input.Prompt + raw)
			m.input.SetValue("")

			if m.OnCommand != nil {
				if cmd := m.OnCommand(raw); cmd != nil {
					return m, cmd
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m ConsoleModel) lastVisibleLines() []string {
	if m.height <= 0 {
		return m.lines
	}

	visible := m.height - 2
	if visible <= 0 {
		return nil
	}
	if len(m.lines) <= visible {
		return m.lines
	}
	return m.lines[len(m.lines)-visible:]
}

func (m ConsoleModel) View() string {
	body := strings.Join(m.lastVisibleLines(), "\n")
	footer := "(enter=run, esc=back, ctrl+c=quit)"

	if body == "" {
		return m.input.View() + "\n" + footer
	}
	return body + "\n" + m.input.View() + "\n" + footer
}
