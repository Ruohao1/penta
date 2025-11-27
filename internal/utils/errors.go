package utils

import "errors"

var (
	ErrSessionExists     = errors.New("session already exists")
	ErrSessionNotFound   = errors.New("session not found")
	ErrIndexNotFound     = errors.New("index not found")
	ErrWorkspaceNotFound = errors.New("workspace not found")
)
