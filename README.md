# UpDoc - Simplified Documentation Flagging System

**Clean, focused MVP for organization and user management with Confluence integration.**

## Quick Start

1. **Start Database:**
   ```bash
   cd backend && docker compose up -d
   ```

2. **Start Server:**
   ```bash
   cd backend && go run ./cmd/server
   ```

3. **Create Organization:**
   ```bash
   curl -X POST http://localhost:9000/api/v1/orgs \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Your Company",
       "user_name": "Admin User", 
       "user_email": "admin@company.com"
     }'
   ```

## Core Features

- ✅ **Organization Management**: Create organizations with admin users
- ✅ **User Creation**: Automatic admin user creation with organizations  
- ✅ **Confluence Integration**: Store Confluence credentials per organization
- ✅ **Connection Testing**: Test Confluence API connectivity
- ✅ **PostgreSQL Storage**: Persistent data with GORM

## API Endpoints

### Organizations

**Create Organization + Admin User:**
```bash
POST /api/v1/orgs
{
  "name": "Acme Corp",
  "user_name": "John Doe",
  "user_email": "john@acme.com",
  "confluence_base_url": "https://acme.atlassian.net/wiki",
  "confluence_email": "john@acme.com", 
  "confluence_token": "your-token",
  "confluence_space_key": "ENG"
}
```

**Get Organization:**
```bash
GET /api/v1/orgs/{slug}
```

**Test Confluence Connection:**
```bash
POST /api/v1/orgs/{id}/test-confluence
```

**List Confluence Pages:**
```bash
GET /api/v1/orgs/{id}/confluence/pages?limit=10
```

## Database Schema

```sql
-- Organizations table
organizations:
  id (uuid, primary key)
  name (text)
  slug (text, unique)
  confluence_base_url (text)
  confluence_email (text)
  confluence_token (text)
  confluence_space_key (text)
  created_at (timestamp)

-- Users table  
users:
  id (uuid, primary key)
  email (text, unique)
  name (text)
  org_id (uuid, foreign key)
  role (text, default: 'member')
  is_active (boolean, default: true)
  created_at (timestamp)
```

## Tech Stack

- **Backend**: Go 1.23 + Echo framework
- **Database**: PostgreSQL with GORM
- **Integration**: Confluence REST API via Resty
- **Deployment**: Docker Compose

## Environment Variables

```bash
# Database (optional - defaults provided)
POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=updoc  
POSTGRES_PASSWORD=updoc
POSTGRES_DB=updoc

# Server
PORT=9000
```

## Development

1. **Run Tests:**
   ```bash
   go test ./...
   ```

2. **Database Reset:**
   ```bash
   docker compose down -v && docker compose up -d
   ```

3. **View Logs:**
   ```bash
   docker compose logs -f
   ```

## Testing Examples

```bash
# Create organization
curl -X POST http://localhost:9000/api/v1/orgs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "DevOps Team",
    "user_name": "Bob Wilson",
    "user_email": "bob@devops.com",
    "confluence_base_url": "https://mycompany.atlassian.net/wiki",
    "confluence_email": "bob@devops.com",
    "confluence_token": "your-real-token",
    "confluence_space_key": "DEV"
  }'

# Get organization
curl http://localhost:9000/api/v1/orgs/devops-team

# Test Confluence connection  
curl -X POST http://localhost:9000/api/v1/orgs/{org-id}/test-confluence

# List Confluence pages
curl http://localhost:9000/api/v1/orgs/{org-id}/confluence/pages?limit=5
```

## Next Steps

This is a clean MVP focused on the core features. Future enhancements:

- Document management and flagging
- Team workspace creation  
- Integration with Teams/Slack
- Advanced user permissions
- Notification system

---

**Status**: ✅ **Ready for testing organization and user creation with Confluence integration**
