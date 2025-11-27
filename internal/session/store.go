package session

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/Ruohao1/penta/internal/utils"
)

type Store interface {
	LoadIndex(ctx context.Context) (*Index, error)
	SaveIndex(ctx context.Context, idx Index) error

	CurrentSession() (*Session, error)
	SetCurrentSession(sessionID string) error

	CreateSession(ctx context.Context, session Session) error
	DeleteSession(ctx context.Context, sessionID string) error
	LastUpdatedSession(ctx context.Context) (string, error)
	ListSessions(ctx context.Context) ([]IndexEntry, error)
}

type FileStore struct {
	BaseDirPath    string
	SessionDirPath string
	LockDirPath    string
	IndexPath      string
	CurrentPath    string
}

func NewFileStore() (*FileStore, error) {
	baseDirPath, err := utils.EnsureStoreLayout()
	if err != nil {
		return nil, err
	}

	sessionDirPath := filepath.Join(baseDirPath, "sessions")
	lockDirPath := filepath.Join(baseDirPath, "locks")
	indexPath := filepath.Join(baseDirPath, "index.json")
	currentPath := filepath.Join(baseDirPath, "current")

	return &FileStore{
		baseDirPath,
		sessionDirPath,
		lockDirPath,
		indexPath,
		currentPath,
	}, nil
}

func (s *FileStore) LoadIndex(ctx context.Context) (*Index, error) {
	logger := utils.LoggerFrom(ctx)
	indexBytes, err := os.ReadFile(s.IndexPath)
	if err != nil {
		logger.Info().Msgf("Index not found, creating new index...")
		return &Index{}, nil
	}

	index, err := utils.FromJSON[Index](indexBytes)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to deserialize index")
		return &Index{}, err
	}
	logger.Debug().Msgf("Index loaded successfully")
	return index, nil
}

func (s *FileStore) SaveIndex(ctx context.Context, idx Index) error {
	logger := utils.LoggerFrom(ctx)
	indexBytes, err := utils.ToJSON(idx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to serialize index")
		return err
	}

	if err = os.WriteFile(s.IndexPath, indexBytes, 0o600); err != nil {
		logger.Error().Err(err).Msg("Failed to save index")
		return err
	}
	logger.Debug().Msgf("Index saved successfully")
	return nil
}

func (s *FileStore) CurrentSession() (*Session, error) {
	currentFile, err := os.ReadFile(s.CurrentPath)
	if err != nil {
		return &Session{}, err
	}

	sessionID := string(currentFile)
	sessionPath, err := utils.LookupSessionPath(sessionID)
	if err != nil {
		return &Session{}, err
	}
	sessionData, err := os.ReadFile(sessionPath)
	if err != nil {
		return &Session{}, err
	}

	return utils.FromJSON[Session](sessionData)
}

func (s *FileStore) SetCurrentSession(ctx context.Context, sessionID string) error {
	logger := utils.LoggerFrom(ctx)

	index, err := s.LoadIndex(ctx)
	if err != nil {
		return err
	}
	for _, entry := range index.Entries {
		if entry.ID == sessionID {
			if err := os.WriteFile(s.CurrentPath, []byte(sessionID), 0o600); err != nil {
				logger.Error().Err(err).Msg("Failed to write current session")
				return err
			}
			return nil
		}
	}
	logger.Error().Msgf("Session not found: %s", sessionID)
	return utils.ErrSessionNotFound
}

func (s *FileStore) CreateSession(ctx context.Context, session Session) error {
	logger := utils.LoggerFrom(ctx)

	info, err := os.Stat(session.Workspace)
	if err == nil && info.IsDir() {
		logger.Error().Msgf("Workspace already exists: %s", session.Workspace)
		return utils.ErrSessionExists
	}

	index, err := s.LoadIndex(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load index")
		return err
	}
	index.Entries = append(index.Entries,
		IndexEntry{
			ID:        session.ID,
			Name:      session.Name,
			Workspace: session.Workspace,
		})

	if err := s.SaveIndex(ctx, *index); err != nil {
		logger.Error().Err(err).Msg("Failed to save index")
		return err
	}

	sessionBytes, err := utils.ToJSON(session)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to serialize session")
		return err
	}

	sessionPath, err := utils.LookupSessionPath(session.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to lookup session path")
		return err
	}

	if err := os.WriteFile(
		sessionPath,
		sessionBytes, 0o600); err != nil {
		logger.Error().Err(err).Msg("Failed to write session")
		return err
	}

	if err := os.MkdirAll(session.Workspace, 0o700); err != nil {
		logger.Error().Err(err).Msg("Failed to create workspace directory")
		return err
	}

	pentaWorkspaceDir, err := utils.LookupPentaWorkspaceDir(session.Workspace)
	if err == utils.ErrWorkspaceNotFound {
		if err := os.MkdirAll(pentaWorkspaceDir, 0o700); err != nil {
			logger.Error().Err(err).Msg("Failed to create penta workspace directory")
			return err
		}
	} else if err != nil {
		logger.Error().Err(err).Msg("Failed to lookup penta workspace directory")
		return err
	}

	sessionSummary := SessionSummary{
		ID:        session.ID,
		Name:      session.Name,
		Workspace: session.Workspace,
	}

	sessionSummaryBytes, err := utils.ToJSON(sessionSummary)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to serialize session summary")
		return err
	}

	if err := os.WriteFile(filepath.Join(pentaWorkspaceDir, "session.json"), sessionSummaryBytes, 0o600); err != nil {
		logger.Error().Err(err).Msg("Failed to write session summary")
		return err
	}

	return nil
}

func (s *FileStore) DeleteSession(ctx context.Context, sessionID string) error {
	logger := utils.LoggerFrom(ctx)

	index, err := s.LoadIndex(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load index")
		return err
	}
	for i, entry := range index.Entries {
		if entry.ID == sessionID {
			if err := os.RemoveAll(entry.Workspace); err != nil {
				logger.Error().Err(err).Msg("Failed to remove session directory")
				return err
			}

			index.Entries = append(index.Entries[:i], index.Entries[i+1:]...)
			break
		}
	}
	if err := s.SaveIndex(ctx, *index); err != nil {
		logger.Error().Err(err).Msg("Failed to save index")
		return err
	}

	return nil
}

func (s *FileStore) LastUpdatedSession(ctx context.Context) (string, error) {
	logger := utils.LoggerFrom(ctx)

	index, err := s.LoadIndex(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load index")
		return "", err
	}

	if index.Entries == nil {
		return "", nil
	}

	var latestTime time.Time
	var sessionID string

	for _, entry := range index.Entries {
		sessionPath, err := utils.LookupSessionPath(entry.ID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to lookup session path")
			return "", err
		}
		sessionBytes, err := os.ReadFile(sessionPath)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read session")
			return "", err
		}

		session, err := utils.FromJSON[Session](sessionBytes)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get session info")
			return "", err
		}
		if session.UpdatedAt.After(latestTime) {
			latestTime = session.UpdatedAt
			sessionID = entry.ID
		}
	}

	return sessionID, nil
}
