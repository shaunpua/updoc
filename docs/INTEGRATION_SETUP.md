# UpDoc - Integration Setup Guide

**Version:** 1.0  
**Last Updated:** August 7, 2025

## Complete Setup Process

### 1. **Initial Dashboard Setup** (5 minutes)

```
1. Visit https://updoc.app/setup (your future domain)
2. Sign up with email or SSO
3. Choose your integration stack:

   📋 Common Stacks:
   ☐ Engineering: Teams + GitHub + Confluence
   ☐ Product: Slack + Notion + Linear
   ☐ Operations: Teams + Confluence + ServiceNow
   ☐ Custom: Mix and match
```

### 2. **Teams Integration Setup** (10 minutes)

#### **Step 1: Connect Teams**

```
Dashboard → Integrations → Microsoft Teams → "Connect"
↓
OAuth flow → Authorize UpDoc → Select channels
↓
Bot automatically installed in selected channels
```

#### **Step 2: Configure Teams Bot**

```go
// Your backend handles the OAuth callback
func (h *IntegrationHandler) HandleTeamsCallback(c echo.Context) error {
    code := c.QueryParam("code")

    // Exchange code for tokens
    tokens, err := h.TeamsClient.ExchangeCode(code)

    // Store in database
    integration := Integration{
        Type: "teams",
        TeamID: getCurrentTeam(c),
        Config: map[string]interface{}{
            "access_token": tokens.AccessToken,
            "channels": selectedChannels,
        },
        Status: "connected",
    }

    return h.IntegrationRepo.Save(integration)
}
```

#### **Step 3: Test Connection**

```
Bot posts test message in selected channel:
"🤖 UpDoc is now connected! Type '/flag' to flag a document."
```

### 3. **GitHub Integration Setup** (15 minutes)

#### **Step 1: Install GitHub App**

```
Dashboard → Integrations → GitHub → "Install App"
↓
GitHub permissions page → Select repositories
↓
UpDoc GitHub App installed with webhook permissions
```

#### **Step 2: Configure Trigger Rules**

```jsx
// Frontend configuration UI
<GitHubConfig>
  <RepoSelector repos={userRepos} selected={selectedRepos} />

  <TriggerRules>
    <Rule type="file_pattern" value="docs/**" />
    <Rule type="commit_keywords" value="breaking,deprecate,docs" />
    <Rule type="pr_labels" value="needs-docs,breaking-change" />
  </TriggerRules>

  <DocumentMapping>
    <Mapping
      pattern="api/**"
      document="confluence:api-docs-123"
      priority="high"
    />
  </DocumentMapping>
</GitHubConfig>
```

#### **Step 3: Webhook Processing**

```go
// GitHub webhook handler
func (h *GitHubHandler) HandlePush(c echo.Context) error {
    var payload github.PushPayload
    c.Bind(&payload)

    // Get integration config for this repo
    config := h.getRepoConfig(payload.Repository.FullName)

    // Analyze changes against rules
    flags := h.analyzeChanges(payload, config.Rules)

    // Create flags and notifications
    for _, flag := range flags {
        h.FlagService.CreateFlag(flag)
        h.NotificationService.NotifyTeams(flag)
    }

    return c.JSON(200, map[string]string{"status": "processed"})
}
```

### 4. **Confluence Integration Setup** (5 minutes)

#### **Step 1: API Token**

```
Dashboard → Integrations → Confluence → "Connect"
↓
Form: Confluence URL + API Token
↓
Test connection → Import spaces/pages
```

#### **Step 2: Document Mapping**

```jsx
// Map GitHub patterns to Confluence pages
<DocumentMapper>
  <Mapping
    source="github:api/**"
    target="confluence:space/API Documentation"
    autoFlag={true}
  />

  <Mapping
    source="teams:#product"
    target="confluence:space/Product Requirements"
    autoFlag={false}
  />
</DocumentMapper>
```

### 5. **Slack Integration Setup** (Alternative to Teams)

#### **Step 1: Slack App Installation**

```
Dashboard → Integrations → Slack → "Add to Slack"
↓
Slack OAuth → Select workspace → Authorize permissions
↓
Bot joins selected channels
```

#### **Step 2: Slash Commands**

```
/flag docs.example.com/api "API endpoint changed"
↓
UpDoc creates flag → Notifies doc owner → Posts confirmation
```

## Frontend Dashboard Implementation

### **Tech Stack**

```javascript
// package.json
{
  "dependencies": {
    "next": "^14.0.0",           // React framework
    "tailwindcss": "^3.0.0",     // Styling
    "@tanstack/react-query": "^5.0.0", // API state
    "zustand": "^4.0.0",         // Global state
    "recharts": "^2.0.0",        // Charts
    "lucide-react": "^0.300.0"   // Icons
  }
}
```

### **Key Components**

#### **Integration Manager Component**

```jsx
// components/IntegrationManager.jsx
export function IntegrationManager() {
  const { integrations, connect, disconnect } = useIntegrations();

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {INTEGRATION_TYPES.map((type) => (
        <IntegrationCard
          key={type}
          type={type}
          status={integrations[type]?.status || "disconnected"}
          onConnect={() => connect(type)}
          onDisconnect={() => disconnect(type)}
          onConfigure={() => openConfig(type)}
        />
      ))}
    </div>
  );
}
```

#### **Flag Dashboard Component**

```jsx
// components/FlagDashboard.jsx
export function FlagDashboard() {
  const { flags, filters, updateFilters } = useFlags();

  return (
    <div>
      <FlagFilters filters={filters} onChange={updateFilters} />
      <FlagTable
        flags={flags}
        onStatusChange={handleStatusChange}
        onBulkAction={handleBulkAction}
      />
    </div>
  );
}
```

### **API Integration Pattern**

```javascript
// hooks/useIntegrations.js
export function useIntegrations() {
  const { data: integrations } = useQuery({
    queryKey: ["integrations"],
    queryFn: () => api.get("/integrations"),
  });

  const connectMutation = useMutation({
    mutationFn: (type) => {
      // Redirect to OAuth flow
      window.location.href = `/api/integrations/${type}/connect`;
    },
  });

  return {
    integrations,
    connect: connectMutation.mutate,
    // ... other methods
  };
}
```

## Backend Architecture Changes

### **New Database Tables**

```sql
-- Add to your existing schema
CREATE TABLE integrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL,
    type VARCHAR NOT NULL, -- teams|slack|github|confluence|notion
    status VARCHAR NOT NULL DEFAULT 'disconnected',
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR NOT NULL, -- confluence|notion|github|url
    external_id VARCHAR NOT NULL,
    url VARCHAR,
    title VARCHAR,
    type VARCHAR, -- api-doc|process|policy|guide
    team_id UUID NOT NULL,
    metadata JSONB DEFAULT '{}'
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR NOT NULL, -- github.push|teams.flag|manual.create
    source VARCHAR NOT NULL,
    payload JSONB NOT NULL,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### **Integration Service Layer**

```go
// internal/integrations/service.go
type Service struct {
    repo   Repository
    teams  *teams.Client
    slack  *slack.Client
    github *github.Client
}

func (s *Service) ConnectIntegration(teamID, integrationType string, config map[string]interface{}) error {
    integration := Integration{
        TeamID: teamID,
        Type:   integrationType,
        Config: config,
        Status: "connected",
    }

    // Test connection
    if err := s.testConnection(integration); err != nil {
        integration.Status = "error"
        return err
    }

    return s.repo.Save(integration)
}

func (s *Service) ProcessEvent(event Event) error {
    switch event.Type {
    case "github.push":
        return s.handleGitHubPush(event)
    case "teams.flag":
        return s.handleTeamsFlag(event)
    default:
        return fmt.Errorf("unknown event type: %s", event.Type)
    }
}
```

## Complete User Journey Examples

### **Example 1: Engineering Team Setup**

```
Day 1: Setup
1. Engineering Manager visits UpDoc dashboard
2. Connects Teams (OAuth flow)
3. Installs GitHub App on repos
4. Connects Confluence with API token
5. Maps code patterns to doc pages

Day 2: First Flag
1. Developer pushes breaking API change
2. GitHub webhook fires → UpDoc analyzes diff
3. Auto-flags "API Documentation" page
4. Teams message posted in #engineering
5. Tech writer updates docs → marks resolved

Week 1: Team Adoption
1. Developers see value → start manual flagging
2. PM flags product docs from Teams
3. Support team flags troubleshooting guides
4. Analytics show 47 flags created, avg 2.3 days to resolve
```

### **Example 2: Product Team (Non-Technical)**

```
Setup:
1. Product Manager connects Slack
2. Links Notion workspace
3. No GitHub (non-technical team)
4. Sets up manual flagging workflows

Usage:
1. PM posts in Slack: "Pricing changed"
2. Uses /flag command → flags pricing docs
3. Sales and Marketing notified
4. Multiple teams update their docs
5. All changes tracked in UpDoc dashboard
```

## Implementation Timeline

### **Phase 1: MVP + Frontend** (6 weeks)

- Week 1-2: Basic React dashboard
- Week 3-4: Teams integration UI
- Week 5-6: GitHub webhook processing

### **Phase 2: Multi-Integration** (10 weeks)

- Week 7-8: Slack integration
- Week 9-10: Notion API
- Week 11-12: Advanced document mapping

### **Phase 3: Intelligence** (16 weeks)

- Week 13-14: Smart change analysis
- Week 15-16: Predictive flagging
- Week 17-20: AI-powered suggestions

## Cost Considerations

### **Development Costs**

- Frontend development: 40 hours
- Integration APIs: 60 hours
- Backend extensions: 40 hours
- **Total: ~140 hours** (3-4 weeks solo)

### **Runtime Costs**

```
Free Tier (0-5 teams):
- Database: Free (Supabase/PlanetScale)
- Hosting: Free (Vercel/Railway)
- Integrations: Free (OAuth apps)
Total: $0/month

Pro Tier (5-50 teams):
- Database: $25/month
- Hosting: $20/month
- Third-party APIs: $10/month
Total: $55/month
```

This setup process gives you a complete integration platform where teams can connect all their tools through a single dashboard, with your backend handling all the complexity behind a unified API.
