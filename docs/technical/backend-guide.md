# Backend Code Guide

## **Code Architecture**

### **Dependency Injection Pattern**

```go
// cmd/server/main.go - Application wiring
func main() {
    // Load environment
    cfg := loadConfig()

    // Initialize database
    db := initDatabase(cfg.DatabaseURL)

    // Create repositories
    userRepo := gormstore.NewUserRepository(db)
    flagRepo := gormstore.NewFlagRepository(db)

    // Create services
    docService := doc.NewDocumentService(flagRepo, userRepo, confluenceClient)

    // Create handlers
    docHandler := httphandler.NewDocHandler(docService)

    // Setup router
    router := setupRouter(docHandler)

    // Start server
    router.Start(":9000")
}
```

### **Clean Architecture Layers**

```
┌─────────────────┐
│   HTTP Layer    │  transport/http/
│   (Echo)        │  - Routes, handlers, middleware
├─────────────────┤
│  Business Layer │  internal/doc/
│   (Use Cases)   │  - Business logic, validation
├─────────────────┤
│  Storage Layer  │  internal/storage/
│   (GORM)       │  - Database operations
├─────────────────┤
│ External APIs   │  internal/providers/
│  (Confluence)   │  - External service clients
└─────────────────┘
```

## **Core Components**

### **HTTP Handlers** (`transport/http/`)

**Router Setup:**

```go
// router.go
func NewRouter(h *DocHandler) *echo.Echo {
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.CORS())

    // Health check
    e.GET("/health", h.Health)

    // API routes
    api := e.Group("/api/v1")
    api.GET("/users", h.GetUsers)
    api.POST("/users", h.CreateUser)
    api.GET("/flags", h.GetFlags)
    api.POST("/flags", h.CreateFlag)

    return e
}
```

**Handler Implementation:**

```go
// doc_handler.go
type DocHandler struct {
    docService *doc.DocumentService
}

func (h *DocHandler) CreateFlag(c echo.Context) error {
    var req CreateFlagRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(400, "Invalid request")
    }

    if err := c.Validate(&req); err != nil {
        return echo.NewHTTPError(400, err.Error())
    }

    flag, err := h.docService.CreateFlag(c.Request().Context(), req)
    if err != nil {
        return echo.NewHTTPError(500, err.Error())
    }

    return c.JSON(201, flag)
}
```

### **Business Logic** (`internal/doc/`)

**Use Cases:**

```go
// usecase.go
type DocumentService struct {
    flagRepo    FlagRepository
    userRepo    UserRepository
    confluence  ConfluenceClient
}

func (s *DocumentService) CreateFlag(ctx context.Context, req CreateFlagRequest) (*Flag, error) {
    // 1. Validate business rules
    if err := s.validateFlagRequest(req); err != nil {
        return nil, err
    }

    // 2. Enrich with external data
    pageInfo, err := s.confluence.GetPageInfo(req.DocURL)
    if err != nil {
        log.Printf("Failed to fetch page info: %v", err)
        // Continue without page info
    }

    // 3. Create domain object
    flag := &Flag{
        ID:          generateID(),
        Title:       req.Title,
        Description: req.Description,
        DocURL:      req.DocURL,
        Status:      StatusPending,
        Priority:    req.Priority,
        CreatedAt:   time.Now(),
    }

    if pageInfo != nil {
        flag.ConfluencePageID = pageInfo.ID
    }

    // 4. Persist to database
    return s.flagRepo.Create(ctx, flag)
}
```

**Domain Models:**

```go
// flag.go
type Flag struct {
    ID                string    `json:"id"`
    UserID           string    `json:"user_id"`
    Title            string    `json:"title"`
    Description      string    `json:"description"`
    DocURL           string    `json:"doc_url"`
    DocType          string    `json:"doc_type"`
    Status           Status    `json:"status"`
    Priority         Priority  `json:"priority"`
    Tags             []string  `json:"tags"`
    Team             string    `json:"team"`
    ConfluencePageID string    `json:"confluence_page_id,omitempty"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type Status string
const (
    StatusPending    Status = "pending"
    StatusInProgress Status = "in_progress"
    StatusResolved   Status = "resolved"
    StatusArchived   Status = "archived"
)

type Priority string
const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
    PriorityUrgent Priority = "urgent"
)
```

### **Data Layer** (`internal/storage/gormstore/`)

**Repository Pattern:**

```go
// flag_repo.go
type FlagRepository struct {
    db *gorm.DB
}

func (r *FlagRepository) Create(ctx context.Context, flag *Flag) (*Flag, error) {
    model := flagToModel(flag)

    if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
        return nil, fmt.Errorf("failed to create flag: %w", err)
    }

    return modelToFlag(&model), nil
}

func (r *FlagRepository) GetByFilters(ctx context.Context, filters FlagFilters) ([]*Flag, error) {
    var models []FlagModel

    query := r.db.WithContext(ctx)

    if filters.Team != "" {
        query = query.Where("team = ?", filters.Team)
    }

    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }

    if filters.Priority != "" {
        query = query.Where("priority = ?", filters.Priority)
    }

    if filters.Search != "" {
        query = query.Where("title ILIKE ? OR description ILIKE ?",
            "%"+filters.Search+"%", "%"+filters.Search+"%")
    }

    if err := query.Find(&models).Error; err != nil {
        return nil, err
    }

    return modelsToFlags(models), nil
}
```

**Database Models:**

```go
// flag_model.go
type FlagModel struct {
    ID                string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID           string    `gorm:"type:uuid;not null"`
    Title            string    `gorm:"not null"`
    Description      string    `gorm:"type:text"`
    DocURL           string    `gorm:"not null"`
    DocType          string    `gorm:"default:'confluence'"`
    Status           string    `gorm:"default:'pending'"`
    Priority         string    `gorm:"default:'medium'"`
    Tags             postgres.Jsonb `gorm:"type:jsonb"`
    Team             string    `gorm:"not null"`
    ConfluencePageID *string   `gorm:"type:varchar(50)"`
    CreatedAt        time.Time `gorm:"autoCreateTime"`
    UpdatedAt        time.Time `gorm:"autoUpdateTime"`

    User UserModel `gorm:"foreignKey:UserID"`
}

func (FlagModel) TableName() string {
    return "flags"
}
```

### **External APIs** (`internal/providers/confluence/`)

**Client Implementation:**

```go
// client.go
type Client struct {
    baseURL    string
    email      string
    token      string
    httpClient *resty.Client
}

func NewClient(baseURL, email, token string) *Client {
    client := resty.New().
        SetBaseURL(baseURL).
        SetBasicAuth(email, token).
        SetHeader("Accept", "application/json").
        SetTimeout(30 * time.Second)

    return &Client{
        baseURL:    baseURL,
        email:      email,
        token:      token,
        httpClient: client,
    }
}

func (c *Client) GetPageInfo(pageID string) (*PageInfo, error) {
    var response PageResponse

    resp, err := c.httpClient.R().
        SetResult(&response).
        Get(fmt.Sprintf("/rest/api/content/%s", pageID))

    if err != nil {
        return nil, fmt.Errorf("failed to get page info: %w", err)
    }

    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("confluence API error: %d", resp.StatusCode())
    }

    return &PageInfo{
        ID:       response.ID,
        Title:    response.Title,
        Space:    response.Space.Key,
        URL:      response.Links.WebUI,
        Modified: response.Version.When,
    }, nil
}
```

## **Error Handling**

### **HTTP Error Responses**

```go
// Error response format
type ErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}

// Custom error middleware
func ErrorMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            err := next(c)
            if err == nil {
                return nil
            }

            // Handle different error types
            switch e := err.(type) {
            case *echo.HTTPError:
                return c.JSON(e.Code, ErrorResponse{
                    Error:   http.StatusText(e.Code),
                    Message: fmt.Sprintf("%v", e.Message),
                })
            case *ValidationError:
                return c.JSON(400, ErrorResponse{
                    Error:   "Validation Failed",
                    Message: e.Message,
                    Details: e.Details,
                })
            default:
                log.Printf("Unhandled error: %v", err)
                return c.JSON(500, ErrorResponse{
                    Error:   "Internal Server Error",
                    Message: "An unexpected error occurred",
                })
            }
        }
    }
}
```

### **Business Logic Errors**

```go
// Custom error types
type ValidationError struct {
    Message string
    Details map[string]string
}

func (e *ValidationError) Error() string {
    return e.Message
}

type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
}

// Usage in service layer
func (s *DocumentService) GetFlag(ctx context.Context, id string) (*Flag, error) {
    flag, err := s.flagRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &NotFoundError{Resource: "flag", ID: id}
        }
        return nil, fmt.Errorf("failed to get flag: %w", err)
    }

    return flag, nil
}
```

## **Validation**

### **Request Validation**

```go
// Request structs with validation tags
type CreateFlagRequest struct {
    Title       string   `json:"title" validate:"required,min=3,max=200"`
    Description string   `json:"description" validate:"required,min=10,max=1000"`
    DocURL      string   `json:"doc_url" validate:"required,url"`
    Priority    Priority `json:"priority" validate:"required,oneof=low medium high urgent"`
    Team        string   `json:"team" validate:"required,min=2,max=50"`
    Tags        []string `json:"tags" validate:"dive,min=1,max=30"`
}

// Custom validator
func setupValidator() *validator.Validate {
    v := validator.New()

    // Register custom validations
    v.RegisterValidation("priority", validatePriority)
    v.RegisterValidation("status", validateStatus)

    return v
}

func validatePriority(fl validator.FieldLevel) bool {
    priority := fl.Field().String()
    return priority == "low" || priority == "medium" || priority == "high" || priority == "urgent"
}
```

## **Configuration Management**

### **Environment Configuration**

```go
// config.go
type Config struct {
    Port        string
    DatabaseURL string
    LogLevel    string
    Confluence  ConfluenceConfig
}

type ConfluenceConfig struct {
    BaseURL string
    Email   string
    Token   string
}

func LoadConfig() *Config {
    godotenv.Load() // Load .env file

    return &Config{
        Port:        getEnv("PORT", "9000"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        Confluence: ConfluenceConfig{
            BaseURL: getEnv("CONF_BASE", ""),
            Email:   getEnv("CONF_EMAIL", ""),
            Token:   getEnv("CONF_TOKEN", ""),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

## **Database Migrations**

### **GORM Auto-Migration**

```go
// init.go
func InitDatabase(databaseURL string) *gorm.DB {
    db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate schemas
    err = db.AutoMigrate(
        &UserModel{},
        &FlagModel{},
    )
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    return db
}
```

## **Testing Strategy**

### **Unit Tests**

```go
// usecase_test.go
func TestDocumentService_CreateFlag(t *testing.T) {
    // Setup
    mockFlagRepo := &MockFlagRepository{}
    mockUserRepo := &MockUserRepository{}
    mockConfluence := &MockConfluenceClient{}

    service := doc.NewDocumentService(mockFlagRepo, mockUserRepo, mockConfluence)

    // Test case
    req := CreateFlagRequest{
        Title:       "Test Flag",
        Description: "Test description",
        DocURL:      "https://example.com",
        Priority:    "high",
        Team:        "Engineering",
    }

    mockFlagRepo.On("Create", mock.Anything, mock.Anything).Return(&Flag{ID: "123"}, nil)

    // Execute
    flag, err := service.CreateFlag(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "123", flag.ID)
    assert.Equal(t, "Test Flag", flag.Title)
    mockFlagRepo.AssertExpectations(t)
}
```

### **Integration Tests**

```go
// integration_test.go
func TestCreateFlagIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Setup services
    flagRepo := gormstore.NewFlagRepository(db)
    userRepo := gormstore.NewUserRepository(db)
    service := doc.NewDocumentService(flagRepo, userRepo, nil)

    // Create test user
    user := &User{Email: "test@example.com", Name: "Test User"}
    createdUser, err := userRepo.Create(context.Background(), user)
    require.NoError(t, err)

    // Test flag creation
    req := CreateFlagRequest{
        UserID:      createdUser.ID,
        Title:       "Integration Test Flag",
        Description: "Testing integration",
        DocURL:      "https://example.com",
        Priority:    "medium",
        Team:        "Test",
    }

    flag, err := service.CreateFlag(context.Background(), req)

    assert.NoError(t, err)
    assert.NotEmpty(t, flag.ID)
    assert.Equal(t, req.Title, flag.Title)
}
```

## **Logging**

### **Structured Logging**

```go
// logger.go
var log = logrus.New()

func init() {
    log.SetFormatter(&logrus.JSONFormatter{})

    level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
    if err != nil {
        level = logrus.InfoLevel
    }
    log.SetLevel(level)
}

// Usage in handlers
func (h *DocHandler) CreateFlag(c echo.Context) error {
    log.WithFields(logrus.Fields{
        "method": "POST",
        "path":   "/flags",
        "user_id": c.Get("user_id"),
    }).Info("Creating flag")

    // ... handler logic

    log.WithFields(logrus.Fields{
        "flag_id": flag.ID,
        "title":   flag.Title,
    }).Info("Flag created successfully")
}
```

This backend code guide provides a comprehensive overview of the Go codebase structure, patterns, and best practices used in the UpDoc project.
