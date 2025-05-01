package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	app "pg-consumer/internal"
	"sync"
	"syscall"
)

func main() {
	app, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("❌ Failed to initialize application: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	app.Start(ctx, &wg)

	// Wait for SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Termination signal received. Shutting down...")
	cancel()
	wg.Wait()

	// Close DB
	sqlDB, _ := app.DB.DB()
	if err := sqlDB.Close(); err != nil {
		log.Printf("❌ DB close error: %v", err)
	}

	log.Println("✅ Application shutdown complete.")
}
