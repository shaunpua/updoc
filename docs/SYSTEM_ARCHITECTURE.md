# UpDoc - Component Interaction Architecture

## System Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                Frontend Dashboard                                │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐│
│  │  Integration    │ │   Flag Mgmt     │ │   Document      │ │   Team          ││
│  │   Manager       │ │   Dashboard     │ │   Discovery     │ │   Analytics     ││
│  │                 │ │                 │ │                 │ │                 ││
│  │ • Connect Apps  │ │ • View Flags    │ │ • Search Docs   │ │ • Performance   ││
│  │ • Configure     │ │ • Update Status │ │ • Map Patterns  │ │ • Metrics       ││
│  │ • Test Connect  │ │ • Bulk Actions  │ │ • Auto-Discovery│ │ • Reports       ││
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └─────────────────┘│
└─────────────────────────┬───────────────────────────────────────────────────────┘
                          │ HTTP/WebSocket API
                          │
┌─────────────────────────▼───────────────────────────────────────────────────────┐
│                           Unified Go Backend                                    │
│                                                                                 │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐          │
│  │   Auth &     │ │   Flag       │ │ Integration  │ │   Event      │          │
│  │   Teams      │ │   Engine     │ │   Manager    │ │  Processor   │          │
│  │              │ │              │ │              │ │              │          │
│  │ • JWT Auth   │ │ • Create     │ │ • OAuth      │ │ • GitHub     │          │
│  │ • RBAC       │ │ • Update     │ │ • Webhooks   │ │ • Teams      │          │
│  │ • Team Mgmt  │ │ • Resolve    │ │ • API Calls  │ │ • Slack      │          │
│  │ • User Prefs │ │ • Analytics  │ │ • Health     │ │ • Manual     │          │
│  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘          │
│                                                                                 │
└─────────────────┬───────────────────┬───────────────────────────────────────────┘
                  │                   │
                  ▼                   ▼
    ┌─────────────────────┐  ┌─────────────────────────────────────────────────────┐
    │    PostgreSQL       │  │              Integration Layer                      │
    │    Database         │  │                                                     │
    │                     │  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
    │ • users/teams       │  │  │   Teams     │ │   Slack     │ │   GitHub    │   │
    │ • integrations      │  │  │    Bot      │ │    Bot      │ │  Webhooks   │   │
    │ • flags/documents   │  │  │             │ │             │ │             │   │
    │ • events/audit      │  │  │ • OAuth     │ │ • OAuth     │ │ • App       │   │
    │ • analytics cache   │  │  │ • Commands  │ │ • Commands  │ │ • Push      │   │
    │                     │  │  │ • Cards     │ │ • Blocks    │ │ • PR        │   │
    └─────────────────────┘  │  └─────────────┘ └─────────────┘ └─────────────┘   │
                             │                                                     │
                             │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
                             │  │ Confluence  │ │   Notion    │ │   Custom    │   │
                             │  │     API     │ │     API     │ │     APIs    │   │
                             │  │             │ │             │ │             │   │
                             │  │ • Get Page  │ │ • Get Page  │ │ • Webhook   │   │
                             │  │ • Update    │ │ • Update    │ │ • REST API  │   │
                             │  │ • Search    │ │ • Search    │ │ • GraphQL   │   │
                             │  └─────────────┘ └─────────────┘ └─────────────┘   │
                             └─────────────────────────────────────────────────────┘
```

## Event Flow Examples

### **1. GitHub Code Change → Auto Flag**

```
Developer commits code
        ↓
GitHub Webhook → UpDoc API
        ↓
Event Processor analyzes diff
        ↓
Matches patterns → Creates flags
        ↓
Notification Service → Teams/Slack
        ↓
Doc owner gets notification
        ↓
Updates doc → Marks flag resolved
```

### **2. Manual Flag from Teams**

```
User types "/flag docs.company.com API docs outdated"
        ↓
Teams webhook → UpDoc API
        ↓
Parse command → Create flag
        ↓
Store in database
        ↓
Send adaptive card to channel
        ↓
Notify doc owner via DM
```

### **3. Bulk Document Update**

```
Frontend dashboard → Select multiple flags
        ↓
Bulk action API call
        ↓
For each flag: Update status
        ↓
Send notifications to stakeholders
        ↓
Update analytics metrics
```

## Data Flow Architecture

### **Input Sources**

```
Manual Flags (Frontend/Chat) → Event Processor
GitHub Webhooks              → Event Processor
Teams/Slack Commands         → Event Processor
Scheduled Tasks              → Event Processor
API Integrations             → Event Processor
```

### **Processing Pipeline**

```
Event Processor → Validation → Business Logic → Database → Notifications
                      ↓              ↓             ↓           ↓
                  Schema Check   Flag Creation   Persist    Teams/Slack
                  Auth Check     Document Map    Analytics  Email/Web
                  Rate Limit     Priority Calc   Audit Log  Webhooks
```

### **Output Channels**

```
Notifications → Teams Cards, Slack Messages, Email, Webhooks
Dashboard     → Real-time updates via WebSocket
Analytics     → Metrics aggregation and reporting
Audit Logs    → Compliance and debugging
```

## API Architecture

### **REST Endpoints**

```
Auth:
POST /api/auth/login
POST /api/auth/logout
GET  /api/auth/me

Integrations:
GET    /api/integrations           # List all
POST   /api/integrations/teams     # Connect Teams
POST   /api/integrations/github    # Connect GitHub
DELETE /api/integrations/:id       # Disconnect
GET    /api/integrations/:id/test  # Test connection

Flags:
GET    /api/flags                  # List with filters
POST   /api/flags                  # Create flag
PATCH  /api/flags/:id              # Update status
DELETE /api/flags/:id              # Delete flag
POST   /api/flags/bulk             # Bulk operations

Documents:
GET    /api/documents              # List documents
POST   /api/documents/search       # Cross-provider search
GET    /api/documents/:id          # Get document details
POST   /api/documents/map          # Map code to docs

Events:
POST   /api/events/github          # GitHub webhook
POST   /api/events/teams           # Teams webhook
POST   /api/events/slack           # Slack webhook
GET    /api/events/audit           # Audit log

Analytics:
GET    /api/analytics/overview     # Dashboard metrics
GET    /api/analytics/flags        # Flag statistics
GET    /api/analytics/teams        # Team performance
GET    /api/analytics/documents    # Document metrics
```

### **WebSocket Events**

```
Client → Server:
- flag:subscribe      # Subscribe to flag updates
- analytics:subscribe # Subscribe to live metrics

Server → Client:
- flag:created        # New flag created
- flag:updated        # Flag status changed
- flag:resolved       # Flag marked as resolved
- analytics:update    # Live metrics update
```

## Integration Patterns

### **OAuth Flow Pattern**

```go
// Standardized OAuth handler
type OAuthHandler struct {
    provider Provider // teams|slack|github
    config   OAuthConfig
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
    state := generateState()
    url := h.provider.GetAuthURL(state)
    return c.Redirect(302, url)
}

func (h *OAuthHandler) Callback(c echo.Context) error {
    token := h.provider.ExchangeCode(c.QueryParam("code"))
    integration := Integration{
        Type:   h.provider.Name(),
        Config: map[string]interface{}{"token": token},
        Status: "connected",
    }
    return h.saveIntegration(integration)
}
```

### **Event Processing Pattern**

```go
type EventProcessor struct {
    handlers map[string]EventHandler
}

type EventHandler interface {
    Handle(Event) error
    CanHandle(eventType string) bool
}

func (ep *EventProcessor) Process(event Event) error {
    handler, exists := ep.handlers[event.Type]
    if !exists {
        return fmt.Errorf("no handler for event type: %s", event.Type)
    }
    return handler.Handle(event)
}

// Example handlers
type GitHubPushHandler struct{}
type TeamsMessageHandler struct{}
type SlackCommandHandler struct{}
```

### **Provider Abstraction Pattern**

```go
type DocumentProvider interface {
    GetDocument(id string) (*Document, error)
    UpdateDocument(id string, content string) error
    SearchDocuments(query string) ([]Document, error)
    GetMetadata(id string) (*DocumentMetadata, error)
}

type ConfluenceProvider struct{}
type NotionProvider struct{}
type GitHubFileProvider struct{}

// Unified document service
type DocumentService struct {
    providers map[string]DocumentProvider
}

func (ds *DocumentService) GetDocument(ref DocumentReference) (*Document, error) {
    provider := ds.providers[ref.Provider]
    return provider.GetDocument(ref.ID)
}
```

## Security Architecture

### **Authentication Flow**

```
Frontend → JWT Token → Backend API
Backend → Validate Token → Process Request
         → Check Permissions → Return Response
```

### **Integration Security**

```
OAuth Tokens → Encrypted Storage (Database)
API Keys     → Environment Variables
Webhooks     → HMAC Signature Verification
Rate Limiting → Per-team quotas
```

### **Data Security**

```
Database → TLS encryption in transit
Secrets  → HashiCorp Vault or AWS Secrets Manager
Audit    → All actions logged with user context
RBAC     → Role-based access control per team
```

This architecture provides a solid foundation for your universal documentation flagging system, with clear separation of concerns and extensible patterns for adding new integrations.
