package session

import "time"

type Session struct {
	ID        string
	Name      string
	Workspace string

	CreatedAt time.Time
	UpdatedAt time.Time

	// Scope     Scope
	// Targets   []Target
	// Env       Env
	// Artefacts []Artefact
}
