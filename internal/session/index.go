package session

import "encoding/json"

type Index struct {
	Path    string
	Entries []IndexEntry
}

type IndexEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}

func (idx *Index) ToJSON() ([]byte, error) {
	return json.MarshalIndent(idx, "", "  ")
}

func DeserializeToIndex(data []byte) (*Index, error) {
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}
