package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ruohao1/penta/internal/engine"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
)

type ScanState int

const (
	ScanStateIdle ScanState = iota
	ScanStateRunning
	ScanStateDone
	ScanStateError
)

type ScanStats struct {
	HostsTotal   int
	HostsDone    int
	PortsOpen    int
	Findings     int
	StartedAt    time.Time
	LastEventAt  time.Time
	LastErrorMsg string
}

type ScanModel struct {
	engine  *engine.Engine
	runOpts engine.RunOptions

	state ScanState
	err   error

	table   table.Model
	logView viewport.Model
	stats   ScanStats

	eventsCh <-chan model.Event
	cancel   context.CancelFunc

	// keymap
	keys ScanKeyMap
}

type ScanKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Back   key.Binding
	Cancel key.Binding
}

func NewScanKeyMap() ScanKeyMap {
	return ScanKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "cancel scan"),
		),
	}
}

func NewScanModel() ScanModel {
	cols := []table.Column{
		{Title: "Host", Width: 18},
		{Title: "Port", Width: 6},
		{Title: "Proto", Width: 6},
		{Title: "Check", Width: 14},
		{Title: "Severity", Width: 8},
		{Title: "Summary", Width: 40},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(nil),
		table.WithFocused(true),
	)

	logView := viewport.New(80, 6)

	return ScanModel{
		state:   ScanStateIdle,
		table:   t,
		logView: logView,
		keys:    NewScanKeyMap(),
	}
}

func (m ScanModel) Render() string {
	var b strings.Builder

	stateStr := "idle"
	switch m.state {
	case ScanStateRunning:
		stateStr = "running"
	case ScanStateDone:
		stateStr = "done"
	case ScanStateError:
		stateStr = "error"
	}

	elapsed := time.Since(m.stats.StartedAt).Truncate(time.Second)
	fmt.Fprintf(&b, "Scan: %s  |  elapsed: %s  | hosts: %d/%d  | ports: %d  | findings: %d\n",
		stateStr,
		elapsed,
		m.stats.HostsDone,
		m.stats.HostsTotal,
		m.stats.PortsOpen,
		m.stats.Findings,
	)

	if m.stats.LastErrorMsg != "" {
		fmt.Fprintf(&b, "last error: %s\n", m.stats.LastErrorMsg)
	}

	b.WriteString("\n")
	b.WriteString(m.table.View())
	b.WriteString("\n\n")
	b.WriteString("Log:\n")
	b.WriteString(m.logView.View())

	return b.String()
}
