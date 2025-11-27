package utils

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type locker interface {
	WithStoreLock(ctx context.Context, fn func() error) error
	WithSessionLock(ctx context.Context, name string, fn func() error) error

	withLock(ctx context.Context, name string, fn func() error) error
	acquire(ctx context.Context, lockPath string) error
	stale()
}

type FileLocker struct {
	locksDir     string
	staleAfter   time.Duration
	pollInterval time.Duration
}

func NewFileLocker(locksDir string) *FileLocker {
	return &FileLocker{
		locksDir:     locksDir,
		staleAfter:   5 * time.Minute,
		pollInterval: 100 * time.Millisecond,
	}
}

func (l *FileLocker) WithStoreLock(ctx context.Context, fn func() error) error {
	return l.withLock(ctx, "store", fn)
}

func (l *FileLocker) WithSessionLock(ctx context.Context, name string, fn func() error) error {
	return l.withLock(ctx, "session-"+name, fn)
}

func (l *FileLocker) withLock(ctx context.Context, name string, fn func() error) error {
	lockPath := filepath.Join(l.locksDir, name+".lock")

	if err := l.acquire(ctx, lockPath); err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(lockPath) }() // release

	return fn()
}

func (l *FileLocker) acquire(ctx context.Context, lockPath string) error {
	pid := os.Getpid()
	host, _ := os.Hostname()

	for {
		// Respect ctx cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Try to create the lock directory (atomic)
		err := os.Mkdir(lockPath, 0o700)
		if err == nil {
			// We own the lock, write metadata for debugging / stale detection
			ownerPath := filepath.Join(lockPath, "owner")
			data := fmt.Sprintf(
				"pid=%d\nhost=%s\nat=%s\nver=%s\n",
				pid,
				host,
				time.Now().UTC().Format(time.RFC3339Nano),
				runtime.Version(),
			)

			_ = os.WriteFile(ownerPath, []byte(data), 0o600)
			return nil
		}

		// If it's not "already exists", something else is wrong
		if !errors.Is(err, fs.ErrExist) {
			return err
		}

		// Lock directory exists → check if it looks stale
		info, statErr := os.Stat(lockPath)
		if statErr == nil {
			age := time.Since(info.ModTime())
			if age > l.staleAfter {
				// Best-effort stale cleanup.
				// Race is fine: mkdir is still atomic, either we or someone else will win after this.
				_ = os.RemoveAll(lockPath)
				continue
			}
		}

		// Not stale (or can't stat) → back off and retry
		time.Sleep(l.pollInterval)
	}
}
