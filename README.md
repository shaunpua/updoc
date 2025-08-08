# UpDoc - Universal Documentation Management

A monorepo containing the full-stack UpDoc application for flagging and tracking documentation updates across teams.

## Quick Start

### 1. **Environment Setup**

```bash
# Copy environment templates
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local
# Edit .env files with your configuration
```

### 2. **Development (Full Stack)**

```bash
# One command to rule them all
make quickstart

# Or manually:
make setup  # Install dependencies
make dev    # Start everything
```

### 3. **Access the Application**

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:9000
- **Database**: localhost:5432

2. **Test everything works:**

   ```bash
   ./test-mvp.sh  # Automated validation script
   ```

3. **Start development:**

   ```bash
   docker compose up -d db
   go run ./cmd/server/
   ```

4. **Test the API:**
   ```bash
   curl http://localhost:9000/health
   curl http://localhost:9000/v1/docs/YOUR_CONFLUENCE_PAGE_ID
   ```

**See [Setup Guide](docs/SETUP_GUIDE.md) for detailed instructions.**

## Documentation

| Document                                                   | Purpose                                                  |
| ---------------------------------------------------------- | -------------------------------------------------------- |
| **[Setup Guide](docs/SETUP_GUIDE.md)**                     | Complete setup instructions and testing (start here)     |
| **[MVP Roadmap](docs/MVP_ROADMAP.md)**                     | 6-week plan from foundation to validation                |
| **[Architecture Guide](docs/ARCHITECTURE.md)**             | Simple overview and MVP strategy                         |
| **[Complete Architecture](docs/COMPLETE_ARCHITECTURE.md)** | Full system design for universal doc management          |
| **[Integration Setup](docs/INTEGRATION_SETUP.md)**         | Step-by-step integration connection process              |
| **[System Architecture](docs/SYSTEM_ARCHITECTURE.md)**     | Technical component interactions and data flows          |
| **[Developer Guide](docs/DEVELOPER.md)**                   | Code structure, API reference, and development workflows |
| **[Operations Guide](docs/OPERATIONS.md)**                 | Commands, deployment, monitoring, and troubleshooting    |

## Features

- **Document Flagging**: Mark Confluence pages as needing updates
- **Status Tracking**: Track document states (pending-update, stale, fresh)
- **REST API**: Simple HTTP interface for document operations
- **Database Persistence**: PostgreSQL storage with automatic migrations
- **Extensible Design**: Ready for Teams integration and additional doc providers

## Core API

- `GET /v1/docs/:id` - Get document details and flags
- `POST /v1/docs/:id` - Update document and/or create flags
- `GET /health` - Health check endpoint

## Tech Stack

- **Backend**: Go 1.23 with Echo framework
- **Database**: PostgreSQL with GORM ORM
- **External**: Confluence REST API integration
- **Deployment**: Docker Compose ready

## Next Steps

This project is architected to support:

- Microsoft Teams integration (message extensions, bots, adaptive cards)
- Additional document providers (Notion, etc.)
- Webhook-based automatic flagging
- Notification systems and workflows

See the [Architecture Guide](docs/ARCHITECTURE.md) for detailed integration planning.
