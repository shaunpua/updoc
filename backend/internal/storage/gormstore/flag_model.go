package gormstore

import "time"

// Flag represents an issue with documentation that needs to be addressed
type Flag struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DocumentID  string     `json:"document_id" gorm:"not null;type:uuid"`
	CreatedBy   string     `json:"created_by" gorm:"not null;type:uuid"`
	AssignedTo  *string    `json:"assigned_to" gorm:"type:uuid"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text;not null"`
	Priority    string     `json:"priority" gorm:"default:'medium'"` // urgent, high, medium, low
	Status      string     `json:"status" gorm:"default:'pending'"`  // pending, in_progress, resolved, archived
	Resolution  string     `json:"resolution" gorm:"type:text"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Document      Document       `gorm:"foreignKey:DocumentID"`
	Creator       User           `gorm:"foreignKey:CreatedBy"`
	Assignee      *User          `gorm:"foreignKey:AssignedTo"`
	Notifications []Notification `gorm:"foreignKey:FlagID"`
}

func (Flag) TableName() string {
	return "flags"
}
