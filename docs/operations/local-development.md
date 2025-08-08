# Local Development

## **Quick Start (5 minutes)**

### **Prerequisites**

- Go 1.23+
- Docker
- Git

### **Setup Commands**

```bash
# 1. Clone and setup
git clone <repo>
cd updoc/backend

# 2. Environment configuration
cp .env.example .env
# Edit .env with your Confluence credentials

# 3. Start database
docker compose up -d db
sleep 3

# 4. Start API server
go run ./cmd/server/
```

**✅ Success**: Server starts on http://localhost:9000

## **Environment Configuration**

### **Required Variables (.env)**

```bash
# Confluence API (get token: https://id.atlassian.com/manage-profile/security/api-tokens)
CONF_BASE=https://your-company.atlassian.net/wiki
CONF_EMAIL=your-email@company.com
CONF_TOKEN=your-confluence-api-token

# Database (auto-configured for Docker)
DATABASE_URL=postgres://updoc:updoc@localhost:5433/updoc?sslmode=disable

# Server
PORT=9000
LOG_LEVEL=info
```

## **Development Commands**

### **Server Management**

```bash
# Start server (auto-reload on changes)
go run ./cmd/server/

# Build binary
go build -o server ./cmd/server/

# Run binary
./server

# Stop server
Ctrl+C or pkill -f "cmd/server"
```

### **Database Management**

```bash
# Start database
docker compose up -d db

# Stop database
docker compose down

# Reset database (lose all data)
docker compose down -v
docker compose up -d db

# Connect to database
docker exec -it $(docker compose ps -q db) psql -U updoc updoc
```

### **Testing Commands**

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test ./internal/doc/

# Run tests with coverage
go test -cover ./...
```

## **API Testing**

### **Health Check**

```bash
curl http://localhost:9000/health
```

### **Create Test Data**

```bash
# Create user
curl -X POST http://localhost:9000/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User","teams":["Engineering"]}'

# Create flag
curl -X POST http://localhost:9000/flags \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Flag","description":"Testing","doc_url":"https://example.com","priority":"high","team":"Engineering"}'

# List flags
curl http://localhost:9000/flags

# List users
curl http://localhost:9000/users
```

### **Test Different Scenarios**

```bash
# Filter flags by team
curl "http://localhost:9000/flags?team=Engineering"

# Filter by status
curl "http://localhost:9000/flags?status=pending"

# Filter by priority
curl "http://localhost:9000/flags?priority=high"

# Search flags
curl "http://localhost:9000/flags?search=api"
```

## **Database Inspection**

### **Connect to Database**

```bash
docker exec -it $(docker compose ps -q db) psql -U updoc updoc
```

### **Useful SQL Queries**

```sql
-- View all tables
\dt

-- Count records
SELECT 'users' as table_name, COUNT(*) FROM users
UNION ALL
SELECT 'flags' as table_name, COUNT(*) FROM flags;

-- View recent flags
SELECT title, status, priority, team, created_at
FROM flags
ORDER BY created_at DESC
LIMIT 10;

-- View flags with user info
SELECT f.title, f.status, u.name, u.email
FROM flags f
JOIN users u ON f.user_id = u.id;

-- Exit database
\q
```

## **Code Development**

### **Project Structure**

```
backend/
├── cmd/server/          # Application entry point
│   └── main.go         # Server startup, dependency injection
├── internal/           # Private application code
│   ├── doc/           # Business logic
│   │   ├── flag.go    # Flag domain model
│   │   ├── usecase.go # Business operations
│   │   └── user.go    # User domain model
│   ├── providers/     # External API clients
│   │   └── confluence/ # Confluence API integration
│   ├── storage/       # Database layer
│   │   └── gormstore/ # GORM implementation
│   └── transport/     # HTTP layer
│       └── http/      # Echo handlers and routing
```

### **Making Changes**

**1. Add New API Endpoint**

```bash
# 1. Add route in transport/http/router.go
v1.POST("/new-endpoint", h.NewHandler)

# 2. Add handler in transport/http/doc_handler.go
func (h *DocHandler) NewHandler(c echo.Context) error {
    // Implementation
}

# 3. Restart server
go run ./cmd/server/
```

**2. Modify Database Schema**

```bash
# 1. Update model in internal/storage/gormstore/
# 2. GORM auto-migrates on server start
# 3. For breaking changes, reset database:
docker compose down -v && docker compose up -d db
```

**3. Add Business Logic**

```bash
# 1. Add function to internal/doc/usecase.go
# 2. Update handler to call new function
# 3. Add tests in internal/doc/usecase_test.go
```

## **Debugging**

### **Common Issues**

**Port 5433 already in use:**

```bash
docker compose down
docker compose up -d db
```

**Database connection failed:**

```bash
# Check if container is running
docker compose ps

# Check container logs
docker compose logs db

# Verify environment variables
echo $DATABASE_URL
```

**API returns 500 error:**

```bash
# Check server logs in terminal
# Look for specific error messages
# Common causes:
# - Missing environment variables
# - Database connection issues
# - Invalid JSON in request
```

**Confluence API not working:**

```bash
# Test credentials manually
curl -u "$CONF_EMAIL:$CONF_TOKEN" \
  "$CONF_BASE/rest/api/content?limit=1"

# Check environment variables
echo $CONF_BASE
echo $CONF_EMAIL
# Don't echo CONF_TOKEN for security
```

### **Log Analysis**

```bash
# Server logs show:
# - HTTP requests and responses
# - Database operations
# - External API calls
# - Error stack traces

# Enable debug logging
LOG_LEVEL=debug go run ./cmd/server/
```

## **Development Workflow**

### **Daily Development**

```bash
# Morning setup (once)
cd backend
docker compose up -d db

# Development cycle (repeat)
1. Edit code in IDE
2. Save file
3. Server auto-restarts (if using air/fresh)
4. Test with curl/Postman
5. Check logs for errors

# Evening cleanup
docker compose down
```

### **Feature Development**

```bash
# 1. Create feature branch
git checkout -b feature/new-endpoint

# 2. Develop and test locally
go run ./cmd/server/
curl http://localhost:9000/new-endpoint

# 3. Run tests
go test ./...

# 4. Commit and push
git add .
git commit -m "Add new endpoint"
git push origin feature/new-endpoint
```

### **Hot Reload (Optional)**

```bash
# Install air for auto-reload
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air

# Now server restarts automatically on code changes
```

## **Performance Monitoring**

### **Basic Metrics**

```bash
# Memory usage
ps aux | grep "cmd/server"

# CPU usage
top -p $(pgrep -f "cmd/server")

# Database connections
docker exec $(docker compose ps -q db) \
  psql -U updoc updoc -c "SELECT count(*) FROM pg_stat_activity;"
```

### **Database Performance**

```sql
-- Slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Database size
SELECT pg_size_pretty(pg_database_size('updoc'));

-- Table sizes
SELECT schemaname,tablename,pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname='public';
```

## **Integration Testing**

### **End-to-End Test Script**

```bash
#!/bin/bash
# test-e2e.sh

echo "🧪 Running E2E tests..."

# 1. Health check
curl -f http://localhost:9000/health || exit 1

# 2. Create user
USER_ID=$(curl -s -X POST http://localhost:9000/users \
  -H "Content-Type: application/json" \
  -d '{"email":"e2e@test.com","name":"E2E User","teams":["Test"]}' \
  | jq -r '.id')

# 3. Create flag
FLAG_ID=$(curl -s -X POST http://localhost:9000/flags \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"E2E Test\",\"description\":\"End-to-end test flag\",\"doc_url\":\"https://example.com\",\"priority\":\"high\",\"team\":\"Test\"}" \
  | jq -r '.id')

# 4. Update flag
curl -f -X PUT http://localhost:9000/flags/$FLAG_ID \
  -H "Content-Type: application/json" \
  -d '{"status":"resolved"}' || exit 1

echo "✅ E2E tests passed!"
```

Run with: `chmod +x test-e2e.sh && ./test-e2e.sh`
