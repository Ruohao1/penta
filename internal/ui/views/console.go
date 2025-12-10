package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type ConsoleModel struct {
	textInput textinput.Model
	err       error

	prompt string
	lines  []string

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
		textInput: ti,
		err:       nil,
		prompt:    "> ",
		focused:   true,
	}
}

func (m ConsoleModel) Init() tea.Cmd {
	if m.focused {
		return textinput.Blink
	}
	return nil
}

func (m ConsoleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m ConsoleModel) View() string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
