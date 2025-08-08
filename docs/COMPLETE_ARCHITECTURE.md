# UpDoc - Complete Integration Architecture

**Version:** 2.0  
**Last Updated:** August 7, 2025

## Expanded Vision

**UpDoc is a universal documentation flagging system** for any company's source-of-truth documents:

- **Engineering**: API docs, architecture, runbooks
- **Product**: PRDs, feature specs, user guides
- **Operations**: Processes, policies, SOPs
- **HR**: Onboarding, policies, org charts
- **Sales**: Playbooks, competitive intel, pricing

## Complete System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend Dashboard                       │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐   │
│  │  Integration    │ │   Flag Mgmt     │ │   Analytics     │   │
│  │   Manager       │ │   Dashboard     │ │   & Reports     │   │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘   │
└─────────────────────────┬───────────────────────────────────────┘
                          │ REST API
┌─────────────────────────▼───────────────────────────────────────┐
│                     Go Backend API                              │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐            │
│  │   Auth &     │ │   Flag       │ │ Integration  │            │
│  │   Users      │ │   Engine     │ │   Manager    │            │
│  └──────────────┘ └──────────────┘ └──────────────┘            │
└─────────────────┬───────────────────┬───────────────────────────┘
                  │                   │
                  ▼                   ▼
    ┌─────────────────────┐  ┌─────────────────────────────────────┐
    │    PostgreSQL       │  │        Integration Layer            │
    │                     │  │  ┌─────────┐ ┌─────────┐ ┌────────┐│
    │  • Users/Teams      │  │  │ Teams   │ │ Slack   │ │GitHub  ││
    │  • Flags/Docs       │  │  │ Bot     │ │ Bot     │ │Webhooks││
    │  • Integrations     │  │  └─────────┘ └─────────┘ └────────┘│
    │  • Audit Log        │  │  ┌─────────┐ ┌─────────┐ ┌────────┐│
    └─────────────────────┘  │  │Confluence│ │ Notion  │ │Custom  ││
                             │  │   API   │ │   API   │ │  APIs  ││
                             │  └─────────┘ └─────────┘ └────────┘│
                             └─────────────────────────────────────┘
```

## Setup Process & User Journey

### 1. **Admin Setup Flow**

```
Admin visits dashboard → Signs up → Configures integrations → Invites team
```

**Step-by-step:**

1. **Visit `https://updoc.app/setup`**
2. **Choose your stack:**
   ```
   ☐ Teams + Confluence + GitHub
   ☐ Slack + Notion + GitLab
   ☐ Custom combination
   ```
3. **Connect each integration** (OAuth flows)
4. **Configure rules:** Which repos/channels trigger flags
5. **Invite team members** via email/SSO

### 2. **Integration Connection Flow**

#### **Teams Integration Setup**

```
Dashboard → Integrations → Teams → Install App → Authorize → Configure Channels
```

**Technical Flow:**

1. Click "Connect Teams" in dashboard
2. OAuth redirect to Microsoft
3. User authorizes UpDoc app
4. Store Teams tokens in database
5. Install bot in selected channels
6. Test connection with ping

#### **GitHub Integration Setup**

```
Dashboard → Integrations → GitHub → Install App → Select Repos → Configure Rules
```

**Technical Flow:**

1. Click "Connect GitHub" in dashboard
2. Install GitHub App on organization
3. Select repositories to monitor
4. Configure trigger rules:
   - File patterns (`docs/**`, `README.md`)
   - Commit message keywords (`breaking`, `deprecate`)
   - PR labels (`needs-docs`)

#### **Confluence/Notion Setup**

```
Dashboard → Integrations → Confluence → API Token → Test Connection → Import Spaces
```

## Unified Backend Architecture

### Core Components

#### **1. Integration Manager** (`internal/integrations/`)

```go
type IntegrationManager struct {
    Teams      *teams.Client
    Slack      *slack.Client
    GitHub     *github.Client
    Confluence *confluence.Client
    Notion     *notion.Client
}

type Integration struct {
    ID       string `json:"id"`
    Type     string `json:"type"`     // teams|slack|github|confluence|notion
    Status   string `json:"status"`   // connected|disconnected|error
    Config   map[string]interface{} `json:"config"`
    TeamID   string `json:"team_id"`
}
```

#### **2. Unified Flag Engine** (`internal/flagging/`)

```go
type FlagRequest struct {
    Source      string `json:"source"`       // teams|slack|github|manual
    DocumentID  string `json:"document_id"`  // confluence:123 | notion:abc | file:path
    Reason      string `json:"reason"`
    Priority    string `json:"priority"`     // low|medium|high|urgent
    Context     map[string]interface{} `json:"context"` // commit hash, PR link, etc.
    CreatedBy   User   `json:"created_by"`
}

type DocumentReference struct {
    Provider string `json:"provider"` // confluence|notion|file|url
    ID       string `json:"id"`
    URL      string `json:"url"`
    Title    string `json:"title"`
    Type     string `json:"type"`     // api-doc|process|policy|guide
}
```

#### **3. Event Processing** (`internal/events/`)

```go
type EventProcessor struct {
    Handlers map[string]EventHandler
}

type Event struct {
    Type      string      `json:"type"`        // github.push|teams.message|manual.flag
    Source    string      `json:"source"`      // github|teams|slack|web
    Payload   interface{} `json:"payload"`
    Timestamp time.Time   `json:"timestamp"`
}

// Examples:
// github.push → analyze diff → flag related docs
// teams.flag → create flag → notify doc owner
// confluence.update → resolve flags → update status
```

## GitHub Integration Deep Dive

### **Trigger Mechanisms**

#### **1. GitHub Actions Workflow** (Recommended)

```yaml
# .github/workflows/doc-check.yml
name: Documentation Check
on: [push, pull_request]

jobs:
  doc-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 2 # Get diff

      - name: Check for doc-impacting changes
        run: |
          # Analyze git diff
          git diff HEAD~1 --name-only > changed_files.txt

          # Call UpDoc API with changes
          curl -X POST "$UPDOC_API/events/github/push" \
            -H "Authorization: Bearer $UPDOC_TOKEN" \
            -d @- << EOF
          {
            "repository": "$GITHUB_REPOSITORY",
            "commit": "$GITHUB_SHA",
            "changed_files": $(cat changed_files.txt | jq -R . | jq -s .),
            "commit_message": "$GITHUB_EVENT_PATH"
          }
          EOF
```

#### **2. GitHub Webhooks** (Alternative)

```go
// internal/transport/http/webhook_handler.go
func (h *WebhookHandler) HandleGitHubPush(c echo.Context) error {
    var payload github.PushPayload
    if err := c.Bind(&payload); err != nil {
        return err
    }

    // Analyze what changed
    flags := h.GitHubSvc.AnalyzeChanges(payload)

    // Create flags for impacted docs
    for _, flag := range flags {
        h.FlagSvc.CreateFlag(flag)
    }

    return c.JSON(200, map[string]string{"status": "processed"})
}
```

#### **3. Smart Change Detection**

```go
type GitHubAnalyzer struct {
    Patterns map[string][]DocumentReference // file pattern → docs
    AI       *openai.Client // Optional: AI analysis
}

func (g *GitHubAnalyzer) AnalyzeChanges(push github.PushPayload) []FlagRequest {
    flags := []FlagRequest{}

    for _, commit := range push.Commits {
        for _, file := range commit.Modified {
            // Rule-based matching
            if strings.Contains(file, "api/") {
                flags = append(flags, FlagRequest{
                    DocumentID: "confluence:api-docs-123",
                    Reason: fmt.Sprintf("API file %s changed in %s", file, commit.ID),
                    Priority: "medium",
                    Context: map[string]interface{}{
                        "commit": commit.ID,
                        "file": file,
                        "author": commit.Author,
                    },
                })
            }

            // AI-powered analysis (optional)
            if g.AI != nil {
                aiFlags := g.analyzeWithAI(commit, file)
                flags = append(flags, aiFlags...)
            }
        }
    }

    return flags
}
```

## Frontend Dashboard Architecture

### **Technology Stack**

```
Frontend: React/Next.js (or Vue/Nuxt)
State: Zustand/Redux Toolkit
UI: Tailwind + Shadcn/UI
Charts: Recharts
API: React Query/SWR
```

### **Key Pages**

#### **1. Integration Dashboard** (`/integrations`)

```jsx
// Components for each integration type
<IntegrationCard
  type="teams"
  status="connected"
  onConfigure={() => openTeamsConfig()}
  onDisconnect={() => handleDisconnect('teams')}
/>

<IntegrationCard
  type="github"
  status="disconnected"
  onConnect={() => initiateGitHubOAuth()}
/>
```

#### **2. Flag Management** (`/flags`)

```jsx
// Unified flag view across all sources
<FlagTable
  filters={{
    source: ["teams", "github", "manual"],
    status: ["pending", "in-progress", "resolved"],
    priority: ["high", "medium", "low"],
  }}
  onBulkAction={handleBulkAction}
/>
```

#### **3. Analytics** (`/analytics`)

```jsx
// Cross-integration metrics
<DashboardGrid>
  <MetricCard title="Flags Created This Week" value="47" trend="+23%" />
  <MetricCard title="Avg Resolution Time" value="2.3 days" trend="-15%" />
  <SourceBreakdown data={flagsBySource} />
  <TeamPerformance data={teamMetrics} />
</DashboardGrid>
```

## Database Schema Evolution

```sql
-- Core tables (existing)
users (id, name, email, role)
flags (id, document_id, status, reason, priority, created_by, created_at)

-- New integration tables
integrations (
  id UUID PRIMARY KEY,
  team_id UUID REFERENCES teams(id),
  type VARCHAR NOT NULL, -- teams|slack|github|confluence|notion
  status VARCHAR NOT NULL, -- connected|disconnected|error
  config JSONB, -- tokens, settings, rules
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

teams (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  plan VARCHAR DEFAULT 'free', -- free|pro|enterprise
  settings JSONB
);

documents (
  id UUID PRIMARY KEY,
  provider VARCHAR NOT NULL, -- confluence|notion|github|url
  external_id VARCHAR NOT NULL, -- provider's document ID
  url VARCHAR,
  title VARCHAR,
  type VARCHAR, -- api-doc|process|policy|guide
  team_id UUID REFERENCES teams(id),
  metadata JSONB
);

events (
  id UUID PRIMARY KEY,
  type VARCHAR NOT NULL, -- github.push|teams.flag|manual.create
  source VARCHAR NOT NULL,
  payload JSONB,
  processed_at TIMESTAMP,
  created_at TIMESTAMP
);
```

## API Architecture

### **Unified REST API**

```
/api/v1/
├── auth/
│   ├── login
│   └── logout
├── integrations/
│   ├── GET /               # List all integrations
│   ├── POST /teams/connect # Connect Teams
│   ├── POST /github/connect # Connect GitHub
│   └── DELETE /:id         # Disconnect
├── flags/
│   ├── GET /               # List flags (with filters)
│   ├── POST /              # Create flag
│   └── PATCH /:id          # Update flag status
├── documents/
│   ├── GET /               # List documents
│   └── POST /search        # Search across providers
├── events/
│   ├── POST /github/push   # GitHub webhook
│   ├── POST /teams/message # Teams webhook
│   └── GET /audit          # Audit log
└── teams/
    ├── GET /members
    └── GET /analytics
```

## Real-World Integration Examples

### **Example 1: Engineering Team**

```
Setup:
- Teams channels: #engineering, #api-team
- GitHub repos: api-service, frontend, docs
- Confluence space: Engineering

Flow:
1. Developer pushes API change to GitHub
2. GitHub Action calls UpDoc API
3. UpDoc flags "API Documentation" in Confluence
4. Teams message posted in #api-team
5. Tech writer updates docs, marks resolved
```

### **Example 2: Product Team**

```
Setup:
- Slack channel: #product
- Notion workspace: Product Hub
- No GitHub (non-technical team)

Flow:
1. PM posts in Slack: "New pricing model approved"
2. Team member flags pricing docs via Slack command
3. Sales and Marketing teams get notified
4. Multiple docs updated across teams
```

### **Example 3: Operations Team**

```
Setup:
- Teams channels: #ops, #security
- Confluence: Operations Runbooks
- GitHub: infrastructure repos

Flow:
1. Infrastructure code changes (Terraform)
2. Auto-flag deployment runbooks
3. Operations team updates procedures
4. Security team reviews compliance docs
```

## Implementation Phases

### **Phase 1: Foundation** (Current + 4 weeks)

- Frontend dashboard (basic)
- Teams integration
- GitHub webhook processing
- Flag management UI

### **Phase 2: Multi-Provider** (8 weeks)

- Slack integration
- Notion API integration
- Document search across providers
- Advanced flag filtering

### **Phase 3: Intelligence** (12 weeks)

- AI-powered change analysis
- Smart doc-to-code mapping
- Predictive flagging
- Advanced analytics

### **Phase 4: Enterprise** (16 weeks)

- SSO integration
- Advanced permissions
- Custom integrations API
- Compliance features

This architecture supports your vision of a universal documentation flagging system that works across all company knowledge sources, not just code. The unified backend can handle any integration while the frontend provides a single pane of glass for managing everything.
