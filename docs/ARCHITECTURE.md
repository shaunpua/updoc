# UpDoc - Simple Architecture Guide

**Version:** 1.1  
**Last Updated:** August 7, 2025

## What UpDoc Does (Simply)

**Problem:** Teams update code but forget to update docs  
**Solution:** Easy way to flag docs as "needs update" and track who fixes them

## Core User Flow

```
1. Someone notices doc is outdated
   ↓
2. They flag it (via Teams/API)
   ↓
3. Doc owner gets notified
   ↓
4. Owner updates doc and marks "fixed"
   ↓
5. Flag cleared, everyone happy
```

## Current Weak Points & Fixes

### ❌ **WEAK:** People won't remember to flag docs

**✅ FIX:** Teams integration makes flagging 2 clicks from chat

### ❌ **WEAK:** No context on WHY doc needs update

**✅ FIX:** Add "reason" field to flags (quick win)

### ❌ **WEAK:** Hard to discover which docs need flagging

**✅ FIX:** Start simple - let users flag what they find. Add automation later.

## Solo Developer Strategy

### Phase 1: MVP (What you have now) ✅

```
[Teams Message] → [Flag API] → [Database] → [Confluence Update]
```

**Cost:** $0 (free tiers everywhere)
**Time:** Add Teams bot (1-2 weeks)

### Phase 2: Make it sticky (6 months)

- Add "reason" field to flags
- Teams bot for notifications
- Simple dashboard to see all flags

### Phase 3: Growth (1 year)

- GitHub webhook auto-flagging
- Basic AI suggestions (OpenAI API)
- Multi-team support

## Current System (MVP)

**What you have now:**

```
Teams Bot ←→ Your Go API ←→ PostgreSQL
    ↑              ↓
    ↓         Confluence API
User Clicks
```

**Status:** ✅ Working foundation for Confluence flagging

## Future Complete System

**Vision:** Universal documentation flagging across all company knowledge sources

```
Frontend Dashboard → Unified Backend → Multiple Integrations
       ↓                   ↓              ↓
   All Teams         PostgreSQL     Teams|Slack|GitHub
   All Flags         All Data       Confluence|Notion
   All Analytics     One Source     Custom APIs
```

**See [Complete Architecture](COMPLETE_ARCHITECTURE.md) for full system design.**

## Current Components (Simplified)

### 1. **Your Go API** (`cmd/server/main.go`)

- Handles Teams webhooks
- Manages flags in database
- Updates Confluence pages
- **Status:** ✅ Working

### 2. **PostgreSQL Database**

- Stores flags and users
- Simple 2-table design
- **Status:** ✅ Working

### 3. **Confluence Integration**

- Reads/updates pages
- **Status:** ✅ Working

### 4. **Teams Bot**

- Let users flag docs from chat
- Send notifications
- **Status:** ❌ Need to build

## Teams Integration (Simple Version)

### What You Need

1. **Azure Bot Registration** (Free)
2. **Add one endpoint** to your API: `/api/messages`
3. **Teams App Manifest** (JSON file)

### User Experience

```
User: Sees outdated info in Teams
User: Right-clicks → "Flag Doc"
Bot:  Opens form → "Which doc? Why outdated?"
User: Submits form
Bot:  Posts card → "📋 API docs flagged by @user"
Owner: Clicks "Mark Fixed" → Flag cleared
```

### Implementation (1-2 weeks)

```go
// Add to your existing router.go
v1.POST("/api/messages", h.HandleTeamsWebhook)

// New handler in doc_handler.go
func (h *DocHandler) HandleTeamsWebhook(c echo.Context) error {
    // Parse Teams message
    // Create flag in database
    // Send response card
}
```

## Cost Breakdown (Solo Developer)

| Service            | Free Tier               | Paid Tier     |
| ------------------ | ----------------------- | ------------- |
| **PostgreSQL**     | Free (Railway/Supabase) | $5/month      |
| **Server Hosting** | Free (Railway/Render)   | $5/month      |
| **Teams Bot**      | Free (Azure)            | Free          |
| **Confluence API** | Free (existing)         | Free          |
| **Domain**         | Optional                | $10/year      |
| **Total**          | **$0/month**            | **$10/month** |

## Data Model (Keep Simple)

### Flags Table

```sql
flags:
- id (uuid)
- page_id (confluence page)
- status (pending/fixed)
- reason (NEW: why it needs update)
- created_by (user)
- created_at
```

### Users Table

```sql
users:
- id (teams user id)
- name
- email
```

**That's it. No complex schemas.**

## Addressing Weak Points

### 🎯 **Make Flagging Effortless**

```
Current: Remember to go to website, find doc, create flag
Better: Right-click in Teams → Flag doc (2 seconds)
```

### 🎯 **Add Context**

```sql
-- Add this field to flags table
reason TEXT -- "API changed", "Process updated", etc.
```

### 🎯 **Show Impact**

Simple dashboard showing:

- How many flags per team
- Average time to fix
- Most problematic docs

### 🎯 **Prevent Notification Fatigue**

- Only notify doc owner (not whole team)
- Group multiple flags for same doc
- Send weekly summary, not instant pings

## Next Features (Priority Order)

### Week 1-2: Teams Bot MVP

- Basic flag creation from Teams
- Simple notification cards

### Month 1: Polish

- Add "reason" field
- Better notification messages
- Simple web dashboard

### Month 3: Growth Features

- GitHub webhook (auto-flag when code changes)
- Team analytics
- Slack integration

### Month 6: AI Features

- OpenAI integration for smart suggestions
- Auto-categorize flag reasons
- Predict which docs need updates

## Success Metrics (Keep Simple)

1. **Flags created per week** (adoption)
2. **Time from flag → fixed** (effectiveness)
3. **Repeat users** (stickiness)

Start with just these 3 numbers.

## Deployment Strategy

### Development

```bash
# Your current setup works
docker compose up -d db
go run ./cmd/server/
```

### Production (Free/Cheap)

```bash
# Deploy to Railway (free tier)
railway login
railway new updoc
railway add postgresql
railway deploy
```

## Risk Mitigation

### "What if Teams integration fails?"

- Keep API working independently
- Teams is just another client

### "What if I can't get users?"

- Start with your own team
- Solve your own problem first

### "What if costs spiral?"

- Everything has generous free tiers
- Revenue before scale

## The Real Strategy

1. **Build Teams integration** (make flagging effortless)
2. **Add reason field** (give context)
3. **Get 5 teams using it** (prove value)
4. **Add GitHub webhooks** (show automation potential)
5. **Charge $5/team/month** (sustainability)

**Focus:** Solve the adoption problem first, features second.
