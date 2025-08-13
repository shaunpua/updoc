package gormstore

import "time"

// Organization represents a company/team using UpDoc
type Organization struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Slug      string    `json:"slug" gorm:"unique;not null"` // acme-corp
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Confluence Integration (simplified)
	ConfluenceBaseURL string `json:"confluence_base_url" gorm:"column:confluence_base_url"`
	ConfluenceEmail   string `json:"confluence_email" gorm:"column:confluence_email"`
	ConfluenceToken   string `json:"confluence_token" gorm:"column:confluence_token"`
	ConfluenceSpaceKey string `json:"confluence_space_key" gorm:"column:confluence_space_key"`

	// Relationships
	Users      []User      `gorm:"foreignKey:OrgID"`
	Workspaces []Workspace `gorm:"foreignKey:OrgID"`
}

func (Organization) TableName() string {
	return "organizations"
}

// User represents a person with access to documentation workspaces
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Name      string    `json:"name" gorm:"not null"`
	OrgID     string    `json:"org_id" gorm:"not null;type:uuid"`
	Role      string    `json:"role" gorm:"default:'member'"` // admin, member
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Organization   Organization   `gorm:"foreignKey:OrgID"`
	CreatedFlags   []Flag         `gorm:"foreignKey:CreatedBy"`
	AssignedFlags  []Flag         `gorm:"foreignKey:AssignedTo"`
	OwnedDocuments []Document     `gorm:"foreignKey:OwnerID"`
	Notifications  []Notification `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

// Workspace represents a collection of documents (e.g., "Engineering Docs")
type Workspace struct {
	ID                string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrgID             string                 `json:"org_id" gorm:"not null;type:uuid"`
	Name              string                 `json:"name" gorm:"not null"`
	IntegrationType   string                 `json:"integration_type"` // confluence, notion, github
	IntegrationConfig map[string]interface{} `json:"integration_config" gorm:"type:jsonb"`
	IsDefault         bool                   `json:"is_default" gorm:"default:false"`
	CreatedAt         time.Time              `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrgID"`
	Documents    []Document   `gorm:"foreignKey:WorkspaceID"`
}

func (Workspace) TableName() string {
	return "workspaces"
}

// Document represents a trackable piece of documentation
type Document struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	WorkspaceID string    `json:"workspace_id" gorm:"not null;type:uuid"`
	Title       string    `json:"title" gorm:"not null"`
	URL         string    `json:"url" gorm:"not null;unique"`
	ExternalID  string    `json:"external_id"` // page_id, file_path, etc.
	OwnerID     string    `json:"owner_id" gorm:"type:uuid"`
	LastChecked time.Time `json:"last_checked"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Workspace Workspace `gorm:"foreignKey:WorkspaceID"`
	Owner     *User     `gorm:"foreignKey:OwnerID"`
	Flags     []Flag    `gorm:"foreignKey:DocumentID"`
}

func (Document) TableName() string {
	return "documents"
}

// Notification represents an alert sent to a user about flag activity
type Notification struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string     `json:"user_id" gorm:"not null;type:uuid"`
	FlagID    string     `json:"flag_id" gorm:"not null;type:uuid"`
	Type      string     `json:"type" gorm:"not null"` // flag_created, flag_assigned, flag_resolved
	Message   string     `json:"message" gorm:"type:text"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
	Flag Flag `gorm:"foreignKey:FlagID"`
}

func (Notification) TableName() string {
	return "notifications"
}
