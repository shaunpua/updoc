# System Architecture

## **Component Overview**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   React/Next.js │◄──►│   Go/Echo       │◄──►│   PostgreSQL    │
│   Port 3000     │    │   Port 9000     │    │   Port 5433     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   Confluence    │
                    │   REST API      │
                    └─────────────────┘
```

## **Tech Stack**

### **Backend (Go)**

- **Framework**: Echo v4 (HTTP router)
- **ORM**: GORM (PostgreSQL driver)
- **Environment**: godotenv
- **HTTP Client**: go-resty (Confluence API)
- **Validation**: go-playground/validator

### **Database (PostgreSQL)**

- **Version**: 16
- **ORM**: GORM with auto-migrations
- **Connection**: pgx driver
- **Deployment**: Docker container

### **Frontend (React/Next.js)**

- **Framework**: Next.js 14 + TypeScript
- **Styling**: Tailwind CSS
- **State**: React Query + Zustand
- **API Client**: Axios with interceptors

### **External Integrations**

- **Confluence**: REST API v2 (read/write pages)
- **Teams**: Bot Framework (planned)
- **GitHub**: Webhooks (planned)

## **Data Flow**

### **Flag Creation Flow**

```
1. User creates flag via API
2. Backend validates data
3. Stores in PostgreSQL
4. Optionally fetches Confluence page info
5. Returns flag with metadata
```

### **Confluence Integration Flow**

```
1. API receives page URL
2. Extracts page ID from URL
3. Calls Confluence REST API
4. Caches page metadata
5. Links flag to page
```

## **Repository Structure**

```
updoc/
├── backend/                 # Go API server
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Private packages
│   │   ├── doc/           # Business logic
│   │   ├── providers/     # External API clients
│   │   ├── storage/       # Database layer
│   │   └── transport/     # HTTP handlers
│   ├── go.mod             # Dependencies
│   └── docker-compose.yaml # Local database
├── frontend/               # React application
│   ├── src/               # Source code
│   ├── public/            # Static assets
│   └── package.json       # Dependencies
├── shared/                 # Shared types
│   └── types/             # TypeScript definitions
└── deploy/                 # Infrastructure
    ├── terraform/         # AWS deployment
    └── k8s/              # Kubernetes configs
```

## **API Design**

### **REST Endpoints**

```
GET    /health              # Service health
GET    /users               # List users
POST   /users               # Create user
GET    /flags               # List flags (with filters)
POST   /flags               # Create flag
PUT    /flags/{id}          # Update flag
DELETE /flags/{id}          # Delete flag
```

### **Request/Response Format**

- **Content-Type**: `application/json`
- **Authentication**: Bearer token (planned)
- **Error Format**: Consistent JSON error responses
- **Validation**: Request body validation with detailed errors

## **Database Schema**

### **Core Tables**

```sql
users:
  id          UUID PRIMARY KEY
  email       VARCHAR UNIQUE
  name        VARCHAR
  teams       JSONB
  created_at  TIMESTAMP
  updated_at  TIMESTAMP

flags:
  id                 UUID PRIMARY KEY
  user_id           UUID REFERENCES users(id)
  title             VARCHAR
  description       TEXT
  doc_url           VARCHAR
  doc_type          VARCHAR
  status            VARCHAR (pending|in_progress|resolved|archived)
  priority          VARCHAR (low|medium|high|urgent)
  tags              JSONB
  team              VARCHAR
  confluence_page_id VARCHAR
  due_date          TIMESTAMP
  resolved_at       TIMESTAMP
  created_at        TIMESTAMP
  updated_at        TIMESTAMP
```

## **Environment Configuration**

### **Required Variables**

```bash
# Database
DATABASE_URL=postgres://user:pass@host:port/db

# Confluence API
CONF_BASE=https://company.atlassian.net/wiki
CONF_EMAIL=email@company.com
CONF_TOKEN=api-token

# Server
PORT=9000
ENV=development|production
LOG_LEVEL=info|debug
```

## **Integration Points**

### **Confluence API**

- **Base URL**: `/wiki/rest/api/content`
- **Authentication**: Basic auth (email + API token)
- **Operations**: Read page content, update pages
- **Rate Limiting**: 1000 requests/hour per token

### **Teams Integration (Planned)**

- **Bot Framework**: Azure Bot Service
- **Webhooks**: `/api/messages` endpoint
- **Cards**: Adaptive Cards for rich UI
- **Authentication**: Teams SSO

### **GitHub Integration (Planned)**

- **Webhooks**: Push events, PR events
- **Auto-flagging**: Based on file changes
- **API**: GitHub Apps for repository access

## **Security Considerations**

### **Current**

- Environment variables for secrets
- Input validation on all endpoints
- PostgreSQL parameterized queries
- CORS configuration

### **Planned**

- JWT authentication
- Rate limiting per user
- Audit logging
- Secrets management (AWS/Azure)

## **Deployment Architecture**

### **Development**

```
Local Machine:
├── Docker (PostgreSQL)
├── Go server (localhost:9000)
└── React dev server (localhost:3000)
```

### **Production (AWS)**

```
Internet → ALB → ECS Fargate (Go) → RDS PostgreSQL
                     ↓
              S3 + CloudFront (React)
```

## **Scalability Considerations**

### **Current Scale**

- Single server instance
- One database connection
- Suitable for small teams (<100 users)

### **Future Scale**

- Horizontal scaling with load balancer
- Connection pooling
- Read replicas for queries
- Redis caching layer
- CDN for static assets

## **Monitoring & Observability**

### **Current**

- Application logs to stdout
- Health check endpoint
- Basic error handling

### **Planned**

- Structured logging (JSON)
- Metrics collection (Prometheus)
- Distributed tracing
- Error tracking (Sentry)
- Performance monitoring
