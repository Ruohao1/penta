package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func storeDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".local", "share", "penta"), nil
}

func EnsureStoreLayout() (string, error) {
	storeDir, err := storeDirPath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(storeDir, "locks"), 0o700); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(storeDir, "sessions"), 0o700); err != nil {
		return "", err
	}

	return storeDir, nil
}

func LookupSessionPath(sessionID string) (string, error) {
	storeDir, err := storeDirPath()
	if err != nil {
		return "", err
	}

	p := filepath.Join(storeDir, "sessions", sessionID+".json")

	_, err = os.Stat(p)
	if err == nil {
		return p, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return p, ErrSessionNotFound
	}
	return "", err
}

func LookupIndexPath() (string, error) {
	storeDir, err := storeDirPath()
	if err != nil {
		return "", err
	}
	p := filepath.Join(storeDir, "index.json")

	_, err = os.Stat(p)
	if err == nil {
		return p, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return p, ErrIndexNotFound
	}
	return "", err
}

func LookupPentaWorkspaceDir(workspace string) (string, error) {
	p := filepath.Join(workspace, ".penta")

	_, err := os.Stat(p)
	if err == nil {
		return p, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return p, ErrWorkspaceNotFound
	}
	return "", err
}
