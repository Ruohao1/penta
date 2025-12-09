package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Next       key.Binding
	Prev       key.Binding
	Quit       key.Binding
	ToggleHelp key.Binding
}

type HelpComponent struct {
	Help  help.Model
	Keys  KeyMap
	Shown bool
}

func NewHelpComponent() HelpComponent {
	h := help.New()
	h.ShowAll = false
	return HelpComponent{
		Help:  h,
		Keys:  KeyMap{}, // will be set per-view
		Shown: true,
	}
}

// isZero tells us whether a binding is effectively unset.
func isZero(b key.Binding) bool {
	return len(b.Keys()) == 0 &&
		b.Help().Key == "" &&
		b.Help().Desc == ""
}

// ShortHelp implements help.Helper.
// This is the compact line shown at the bottom.
func (k KeyMap) ShortHelp() []key.Binding {
	var bindings []key.Binding

	add := func(b key.Binding) {
		if !isZero(b) {
			bindings = append(bindings, b)
		}
	}

	// Order for short help: nav + quit + toggle
	add(k.Up)
	add(k.Down)
	add(k.Next)
	add(k.Prev)
	add(k.Quit)
	add(k.ToggleHelp)

	return bindings
}

// FullHelp implements help.Helper.
// This is the expanded multi-line help.
func (k KeyMap) FullHelp() [][]key.Binding {
	var nav, global []key.Binding

	add := func(dst *[]key.Binding, b key.Binding) {
		if !isZero(b) {
			*dst = append(*dst, b)
		}
	}

	// group 1: navigation
	add(&nav, k.Up)
	add(&nav, k.Down)
	add(&nav, k.Next)
	add(&nav, k.Prev)

	// group 2: global actions
	add(&global, k.Quit)
	add(&global, k.ToggleHelp)

	var rows [][]key.Binding
	if len(nav) > 0 {
		rows = append(rows, nav)
	}
	if len(global) > 0 {
		rows = append(rows, global)
	}
	return rows
}

func (hc *HelpComponent) SetKeys(k KeyMap) {
	hc.Keys = k
}

func (hc *HelpComponent) Toggle() {
	if !hc.Shown {
		hc.Shown = true
		hc.Help.ShowAll = false
		return
	}
	hc.Help.ShowAll = !hc.Help.ShowAll
}

func (hc HelpComponent) View() string {
	if !hc.Shown {
		return ""
	}
	return hc.Help.View(hc.Keys)
}
