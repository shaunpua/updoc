# UpDoc API Reference - Simplified

## Base URL
```
http://localhost:9000/api/v1
```

## Organizations

### Create Organization + Admin User
Creates a new organization and an admin user in one operation.

```http
POST /orgs
Content-Type: application/json

{
  "name": "Acme Corporation",
  "user_name": "John Doe",
  "user_email": "john@acme.com",
  "confluence_base_url": "https://acme.atlassian.net/wiki",  // optional
  "confluence_email": "john@acme.com",                       // optional
  "confluence_token": "ATATT3xFfGF0...",                     // optional
  "confluence_space_key": "ENG"                              // optional
}
```

**Response 201:**
```json
{
  "organization": {
    "id": "132fa32f-b3cd-42d4-a4db-ed14539208af",
    "name": "Acme Corporation",
    "slug": "acme-corporation",
    "created_at": "2025-08-13T14:56:08Z",
    "confluence_base_url": "https://acme.atlassian.net/wiki",
    "confluence_email": "john@acme.com",
    "confluence_space_key": "ENG"
  },
  "user": {
    "id": "c14ac557-8e49-47a4-b8ee-42de579b9b28",
    "email": "john@acme.com",
    "name": "John Doe",
    "org_id": "132fa32f-b3cd-42d4-a4db-ed14539208af",
    "role": "admin",
    "is_active": true,
    "created_at": "2025-08-13T14:56:08Z"
  }
}
```

### Get Organization
Retrieves an organization by its slug.

```http
GET /orgs/{slug}
```

**Response 200:**
```json
{
  "id": "132fa32f-b3cd-42d4-a4db-ed14539208af",
  "name": "Acme Corporation", 
  "slug": "acme-corporation",
  "created_at": "2025-08-13T14:56:08Z",
  "confluence_base_url": "https://acme.atlassian.net/wiki",
  "confluence_email": "john@acme.com",
  "confluence_space_key": "ENG"
}
```

## Confluence Integration

### Test Confluence Connection
Tests if the organization's Confluence credentials are valid.

```http
POST /orgs/{org_id}/test-confluence
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Connection successful",
  "details": "Successfully authenticated with Confluence"
}
```

**Response 200 (Failure):**
```json
{
  "success": false,
  "message": "Authentication failed",
  "details": "HTTP 403: Current user not permitted to use Confluence"
}
```

### List Confluence Pages
Gets pages from the organization's configured Confluence space.

```http
GET /orgs/{org_id}/confluence/pages?limit=10
```

**Response 200:**
```json
{
  "pages": [
    {
      "id": "123456",
      "title": "API Documentation",
      "url": "https://acme.atlassian.net/wiki/spaces/ENG/pages/123456",
      "space": "ENG"
    },
    {
      "id": "789012", 
      "title": "Deployment Guide",
      "url": "https://acme.atlassian.net/wiki/spaces/ENG/pages/789012",
      "space": "ENG"
    }
  ],
  "count": 2
}
```

## Error Responses

All endpoints may return these error responses:

**400 Bad Request:**
```json
{
  "message": "name, user_name, and user_email are required"
}
```

**404 Not Found:**
```json
{
  "message": "Organization not found"
}
```

**500 Internal Server Error:**
```json
{
  "message": "Database connection failed"
}
```

## Testing Examples

### 1. Create Organization (Basic)
```bash
curl -X POST http://localhost:9000/api/v1/orgs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Startup Inc",
    "user_name": "Alice Smith", 
    "user_email": "alice@techstartup.com"
  }'
```

### 2. Create Organization (With Confluence)
```bash
curl -X POST http://localhost:9000/api/v1/orgs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "DevOps Team",
    "user_name": "Bob Wilson",
    "user_email": "bob@devops.com",
    "confluence_base_url": "https://mycompany.atlassian.net/wiki",
    "confluence_email": "bob@devops.com", 
    "confluence_token": "ATATT3xFfGF0T8ZGgvJNOEZl...",
    "confluence_space_key": "DEV"
  }'
```

### 3. Get Organization
```bash
curl http://localhost:9000/api/v1/orgs/devops-team
```

### 4. Test Confluence Connection
```bash
curl -X POST http://localhost:9000/api/v1/orgs/9a470035-19c7-4187-82d6-c6a25db03e84/test-confluence
```

### 5. List Confluence Pages
```bash
curl http://localhost:9000/api/v1/orgs/9a470035-19c7-4187-82d6-c6a25db03e84/confluence/pages?limit=5
```

---

**Note**: Replace `{org_id}` with actual organization UUID from the creation response.
