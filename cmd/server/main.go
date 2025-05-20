package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/shaunpua/updoc/internal/doc"
	"github.com/shaunpua/updoc/internal/providers/confluence"
	"github.com/shaunpua/updoc/internal/storage/gormstore"
	transport "github.com/shaunpua/updoc/internal/transport/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	dsn := "host=localhost user=updoc password=updoc dbname=updoc port=5432 sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := gormstore.AutoMigrate(gormDB); err != nil {
		log.Fatal(err)
	}

	userRepo := gormstore.NewUserRepo(gormDB)
	_ = userRepo.Ensure("u-123", "demo user")

	confClient := confluence.New() // your existing wrapper
	store := gormstore.NewFlagRepo(gormDB)
	svc := doc.New(confClient, store)
	e := transport.NewRouter(svc)

	// ---- graceful shutdown (Echo cookbook) ----
	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = e.Shutdown(ctx) // :contentReference[oaicite:4]{index=4}
}
