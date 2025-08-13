package doc

import "time"

// Repository interfaces
type OrganizationRepository interface {
	Create(org *Organization) error
	GetBySlug(slug string) (*Organization, error)
	GetByID(id string) (*Organization, error)
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByOrgID(orgID string) ([]*User, error)
	GetByID(id string) (*User, error)
}

type WorkspaceRepository interface {
	Create(workspace *Workspace) error
	GetByOrgID(orgID string) ([]*Workspace, error)
	GetByID(id string) (*Workspace, error)
	UpdateIntegration(id string, config map[string]interface{}) error
}

type DocumentRepository interface {
	Create(doc *Document) error
	GetByWorkspaceID(workspaceID string) ([]*Document, error)
	GetByURL(url string) (*Document, error)
	GetByID(id string) (*Document, error)
	BulkCreate(docs []*Document) error
}

type FlagRepository interface {
	Create(flag *Flag) error
	GetByID(id string) (*Flag, error)
	GetByDocumentID(documentID string) ([]*Flag, error)
	GetByFilters(filters FlagFilters) ([]*Flag, error)
	Update(flag *Flag) error
}

type NotificationRepository interface {
	Create(notification *Notification) error
	GetByUserID(userID string, limit int) ([]*Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID string) error
}

// Domain Models
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	
	// Confluence Integration
	ConfluenceBaseURL  string `json:"confluence_base_url,omitempty"`
	ConfluenceEmail    string `json:"confluence_email,omitempty"`
	ConfluenceToken    string `json:"-"` // Never expose in JSON
	ConfluenceSpaceKey string `json:"confluence_space_key,omitempty"`
}

type Workspace struct {
	ID                string                 `json:"id"`
	OrgID             string                 `json:"org_id"`
	Name              string                 `json:"name"`
	IntegrationType   string                 `json:"integration_type"`
	IntegrationConfig map[string]interface{} `json:"integration_config"`
	IsDefault         bool                   `json:"is_default"`
	CreatedAt         time.Time              `json:"created_at"`
}

type Document struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	ExternalID  string    `json:"external_id"`
	OwnerID     string    `json:"owner_id"`
	LastChecked time.Time `json:"last_checked"`
	CreatedAt   time.Time `json:"created_at"`
}

type Flag struct {
	ID          string     `json:"id"`
	DocumentID  string     `json:"document_id"`
	CreatedBy   string     `json:"created_by"`
	AssignedTo  *string    `json:"assigned_to"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	Resolution  string     `json:"resolution"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Related entities (populated by repository)
	Document *Document `json:"document,omitempty"`
	Creator  *User     `json:"creator,omitempty"`
	Assignee *User     `json:"assignee,omitempty"`
}

type Notification struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	FlagID    string     `json:"flag_id"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// Request/Response types
type FlagFilters struct {
	WorkspaceID string `json:"workspace_id"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssignedTo  string `json:"assigned_to"`
	CreatedBy   string `json:"created_by"`
	Search      string `json:"search"`
}

type CreateFlagRequest struct {
	DocumentID  string  `json:"document_id" validate:"required"`
	Title       string  `json:"title" validate:"required,min=3,max=200"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Priority    string  `json:"priority" validate:"required,oneof=urgent high medium low"`
	AssignedTo  *string `json:"assigned_to"`
}

type UpdateFlagRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	Status      *string `json:"status"`
	AssignedTo  *string `json:"assigned_to"`
	Resolution  *string `json:"resolution"`
}

// Legacy types (for backward compatibility during migration)
type DocFlag struct {
	ID        string    `json:"id"`
	PageID    string    `json:"page_id"`
	Status    string    `json:"status"`
	Body      string    `json:"body,omitempty"`
	User      User      `json:"user"`
	UpdatedAt time.Time `json:"updated_at"`
}
