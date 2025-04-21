package main

import (
	"log"
	"logs-grpc-server/grpc-kafka-server/internal/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	application, err := grpc.InitApp()
	if err != nil {
		application.Logger.Fatalf("Error initializing application: %v", err)
	}
	go func() {
		application.Run()
	}()
	// Set up a signal channel to handle graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait for a termination signal.

	log.Println("Shutting down server...")
	log.Println("Application shut down gracefully.")
}
