package session

type Index struct {
	Entries []IndexEntry `json:"entries"`
}

type IndexEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}
