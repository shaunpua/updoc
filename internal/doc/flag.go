package doc

import "time"

type FlagStore interface {
	Save(DocFlag) error
	ByPage(pageID string) ([]DocFlag, error)
}

// -------- Flag DTO & interface --------
type DocFlag struct {
	ID        string    `json:"id"`
	PageID    string    `json:"page_id"`
	Status    string    `json:"status"` // pending-update | stale | fresh
	Body      string    `json:"body,omitempty"`
	User      User      `json:"user"`
	UpdatedAt time.Time `json:"updated_at"`
}
