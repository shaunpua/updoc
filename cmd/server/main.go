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
	transport "github.com/shaunpua/updoc/internal/transport/http"
)

func main() {
	_ = godotenv.Load()

	confClient := confluence.New() // your existing wrapper
	store := doc.NewInMemStore()
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
