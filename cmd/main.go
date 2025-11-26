package main

import (
	"context"
	"go-first-project/client"
	"go-first-project/database"
	"go-first-project/handlers"
	helper "go-first-project/internal"
	"go-first-project/repository"
	"go-first-project/services"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	dbConfig := database.Config{
		Host:     helper.GetEnv("DB_HOST", ""),
		Port:     helper.GetEnv("DB_PORT", ""),
		User:     helper.GetEnv("DB_USER", ""),
		Password: helper.GetEnv("DB_PASSWORD", ""),
		DBName:   helper.GetEnv("DB_NAME", ""),
	}

	pool, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.InitSchema(pool); err != nil {
		log.Fatal("Failed to initialize database schema:", err)
	}

	weatherClient := client.NewClient(
		helper.GetEnv("API_KEY_OPEN_WEATHER", ""),
		helper.GetEnv("API_KEY_WEATHER_API", ""),
	)

	wRepo := repository.NewWeatherRepository(pool)
	wServices := services.NewWeatcherService(wRepo, weatherClient)
	wHandlers := handlers.WeatherHandler{WeatcherService: *wServices}

	http.HandleFunc("/weather", wHandlers.GetWeatherHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	serverDone := make(chan struct{}, 1)

	go func() {
		defer func() {
			serverDone <- struct{}{}
		}()

		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	select {
	case sig := <-sigChan:
		log.Printf("Received shutdown signal: %v", sig)
		cancel()

	case <-ctx.Done():
		log.Printf("Context cancelled, initiating shutdown")
	}

	log.Printf("Initiating graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer shutdownCancel()

	log.Printf("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
		database.Close(pool)
		os.Exit(1)
	}

	select {
	case <-serverDone:
		log.Printf("HTTP server stopped gracefully")
	case <-time.After(1 * time.Second):
		log.Printf("HTTP server goroutine cleanup timeout")
	}

	log.Printf("Closing database connections...")
	database.Close(pool) // may be use defer database.Close(pool)
	log.Printf("Database connections closed")

	log.Printf("Shutdown complete")
	os.Exit(0)
}
