package session

import "time"

type Session struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Scope     Scope
	// Targets   []Target
	// Env       Env
	// Artefacts []Artefact
}

type SessionSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}
