# UpDoc Repository Guide - Simple Structure

## What This Project Does
UpDoc helps teams track documentation that needs updates. Think of it as a "flag system" for docs - when you notice something is outdated, you can flag it for someone to fix later.

## Repository Structure (Simple Explanation)

```
updoc/
├── README_SIMPLE.md              # Quick start guide (start here!)
├── docs/
│   ├── API_SIMPLE.md            # How to use the web API
│   └── CLEANUP_SUMMARY.md       # What we cleaned up
└── backend/                     # The main application code
    ├── cmd/server/main.go       # Starts the web server
    ├── docker-compose.yaml      # Database setup
    ├── go.mod & go.sum          # Dependencies list
    └── internal/                # Application code
        ├── doc/                 # Data models (what info we store)
        ├── services/            # Business logic (what the app does)
        ├── storage/             # Database code (how we save data)
        └── transport/http/      # Web API (how users interact)
```

## Core Files Explained

### 🚀 Entry Points
- **`README_SIMPLE.md`** - Start here! Quick setup instructions
- **`backend/cmd/server/main.go`** - The main program that starts everything

### 📊 Data Models (`backend/internal/doc/`)
- **`flag.go`** - Defines what data we store (organizations, users, etc.)
- **`user.go`** - User information structure

### 🏢 Business Logic (`backend/internal/services/`)
- **`organization_service.go`** - Creates organizations and users
- **`confluence_service.go`** - Connects to Confluence (documentation platform)

### 💾 Database Layer (`backend/internal/storage/gormstore/`)
- **`organization_repo.go`** - Saves/loads organizations from database
- **`user_repo.go`** - Saves/loads users from database  
- **`user_model.go`** - Database table structure
- **`init.go`** - Sets up database tables

### 🌐 Web API (`backend/internal/transport/http/`)
- **`organization_handler.go`** - Handles web requests for organizations
- **`router.go`** - Directs web requests to right handler

### ⚙️ Infrastructure
- **`docker-compose.yaml`** - PostgreSQL database setup
- **`go.mod`** - Lists required libraries

## What Each File Does (Non-Technical)

### The Main Program (`main.go`)
1. Connects to the database
2. Sets up the organization and user management systems
3. Starts a web server that listens for requests
4. Handles shutdown gracefully

### Organization Service (`organization_service.go`)
- Creates new companies/teams in the system
- Creates the first admin user for each organization
- Validates that organization names are unique
- Handles Confluence integration setup

### Confluence Service (`confluence_service.go`)  
- Tests if Confluence credentials work
- Fetches list of documents from Confluence
- Handles authentication with Confluence API

### Database Repositories (`*_repo.go`)
- Translates between our program's data and database storage
- Handles creating, reading, updating data
- Ensures data consistency and relationships

### Web Handlers (`*_handler.go`)
- Receives web requests (like from curl or browser)
- Validates input data  
- Calls the appropriate service
- Returns results as JSON

## Data Flow Example

When you create an organization:

1. **Web Request** → `organization_handler.go` receives HTTP POST
2. **Validation** → Handler checks required fields are present  
3. **Business Logic** → `organization_service.go` creates org + admin user
4. **Database** → `organization_repo.go` and `user_repo.go` save to PostgreSQL
5. **Response** → Handler returns success with organization details

## Key Features Working

✅ **Create Organizations** - Companies can sign up  
✅ **Create Users** - Admin users created automatically  
✅ **Confluence Integration** - Store API credentials per organization  
✅ **Test Connections** - Verify Confluence access works  
✅ **List Documents** - Fetch pages from Confluence spaces  

## What We Removed (Cleanup)

❌ **Legacy Flag System** - Old complex flagging code  
❌ **Multiple Document Types** - Simplified to just organizations  
❌ **Complex Service Interfaces** - Direct service calls now  
❌ **Unused Provider Abstractions** - Kept simple Confluence integration  
❌ **Excessive Documentation** - Condensed to essentials  

## Environment Setup

The app needs these environment variables (with defaults):
```bash
POSTGRES_HOST=localhost     # Database location
POSTGRES_PORT=5433         # Database port  
POSTGRES_USER=updoc        # Database username
POSTGRES_PASSWORD=updoc    # Database password
POSTGRES_DB=updoc         # Database name
PORT=9000                 # Web server port
```

## Next Development Steps

The current clean foundation supports adding:
1. **Document Management** - Import docs from Confluence
2. **Flagging System** - Mark docs that need updates
3. **Team Management** - Add more users to organizations  
4. **Notifications** - Alert when docs need attention
5. **Dashboard** - Web interface for managing everything

## Getting Help

- **Quick Start**: Read `README_SIMPLE.md`
- **API Usage**: Check `docs/API_SIMPLE.md`  
- **Database**: Use `docker compose exec db psql -U updoc updoc`
- **Logs**: Check `docker compose logs -f`

The codebase is now clean and focused on the core MVP features, making it much easier to understand and extend.
