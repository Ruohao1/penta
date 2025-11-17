package session

import (
	"os"

	"github.com/Ruohao1/penta/internal/utils"
)

type Store interface {
	LoadIndex() (*Index, error)
	SaveIndex(idx Index) error

	CurrentSession() (Session, error)
	SetCurrentSession(sessionID string) error

	AddSession(session Session) error
	RemoveSession(sessionID string) error
}

type FileStore struct {
	BaseDirPath    string
	SessionDirPath string
	LockDirPath    string

	Index Index
}

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (s *FileStore) LoadIndex() (*Index, error) {
	indexPath, err := utils.LookupIndexPath()
	if err != nil {
		return &Index{}, err
	}

	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		return &Index{}, err
	}

	return DeserializeToIndex(indexBytes)
}
