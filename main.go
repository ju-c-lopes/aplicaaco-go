package main

import (
	"fmt"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"lanchonete/bootstrap"
	"lanchonete/internal/interfaces/http/server"
	_"lanchonete/docs"
)

// @title Lanchonete API - Tech Challenge 2
// @version 1.0
// @description API para o Tech Challenge 2 da FIAP - SOAT

// @host localhost:8080
// @BasePath /
//go:generate go run github.com/swaggo/swag/cmd/swag@latest init
func main() {
	fmt.Println("🔧 Iniciando aplicação...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize application
	app, err := bootstrap.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create and configure HTTP server
	srv := server.NewServer(app)
	srv.SetupRoutes()

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// Perform any cleanup here if needed
}