package doc

import (
	"sync"
	"time"
)

// User represents the actor who flagged / updated docs.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DocFlag is one change request on a Confluence page.
type DocFlag struct {
	ID        string    `json:"id"`      // ← new UUID
	PageID    string    `json:"page_id"` // redundant but handy
	Status    string    `json:"status"`  // pending-update | stale | fresh
	Body      string    `json:"body,omitempty"`
	User      User      `json:"user"` // who flagged it
	UpdatedAt time.Time `json:"updated_at"`
}

type InMemStore struct {
	mu   sync.RWMutex
	data map[string][]DocFlag // pageID → slice of flags
}

func NewInMemStore() *InMemStore {
	return &InMemStore{data: make(map[string][]DocFlag)}
}

func (s *InMemStore) Save(flag DocFlag) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[flag.PageID] = append(s.data[flag.PageID], flag)
}

func (s *InMemStore) FlagsByPage(pageID string) []DocFlag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]DocFlag(nil), s.data[pageID]...) // copy
}
