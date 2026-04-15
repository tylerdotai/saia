package memory

import "fmt"

type Store struct {
	// TODO: FTS5-backed memory store
}

func New() *Store {
	return &Store{}
}

func (s *Store) Search(query string, limit int) ([]string, error) {
	// TODO: FTS5 search
	return nil, fmt.Errorf("not implemented")
}

func (s *Store) Add(content, contextID string) error {
	// TODO: Add to FTS5
	return fmt.Errorf("not implemented")
}

func (s *Store) Recall(ctxID string) ([]string, error) {
	// TODO: Retrieve by context
	return nil, fmt.Errorf("not implemented")
}
