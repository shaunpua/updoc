package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/shaunpua/updoc/internal/services"
	"github.com/shaunpua/updoc/internal/storage/gormstore"
	transport "github.com/shaunpua/updoc/internal/transport/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Database configuration with environment variables
	dsn := getDatabaseURL()
	log.Printf("Connecting to database with DSN (masked): %s", maskDSN(dsn))
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := gormstore.AutoMigrate(gormDB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	orgRepo := gormstore.NewOrganizationRepo(gormDB)
	userRepo := gormstore.NewUserRepo(gormDB)

	// Initialize services
	orgService := services.NewOrganizationService(orgRepo, userRepo)
	confluenceService := services.NewConfluenceService(orgRepo)

	// Setup HTTP router
	e := transport.NewRouter()
	
	// Add organization endpoints
	orgHandler := transport.NewOrganizationHandler(orgService, confluenceService)
	api := e.Group("/api/v1")
	api.POST("/orgs", orgHandler.CreateOrganization)
	api.GET("/orgs/:slug", orgHandler.GetOrganization)
	api.POST("/orgs/:id/test-confluence", orgHandler.TestConfluence)
	api.GET("/orgs/:id/confluence/pages", orgHandler.ListConfluencePages)

	// Get port from environment
	port := getEnv("PORT", "9000")

	log.Printf("Starting server on port %s", port)
	log.Printf("Database: %s", maskDSN(dsn))
	log.Printf("Confluence: %s", getEnv("CONF_BASE", "not configured"))

	// Start server in background
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}

// getDatabaseURL constructs database URL from environment variables
func getDatabaseURL() string {
	// Try DATABASE_URL first (for production)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	// Fallback to individual components (for development)
	host := getEnv("POSTGRES_HOST", "localhost")
	user := getEnv("POSTGRES_USER", "updoc")
	password := getEnv("POSTGRES_PASSWORD", "updoc")
	dbname := getEnv("POSTGRES_DB", "updoc")
	// Default to 5433 because docker-compose maps 5433->5432
	port := getEnv("POSTGRES_PORT", "5433")
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// maskDSN masks sensitive information in database URL for logging
func maskDSN(dsn string) string {
	if len(dsn) > 20 {
		return dsn[:20] + "***"
	}
	return "***"
}
