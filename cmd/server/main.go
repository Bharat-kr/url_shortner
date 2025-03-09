package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	database "github.com/Bharat-kr/url-shortner/internal/storage"
	routes "github.com/Bharat-kr/url-shortner/internal/routes"
	"github.com/rs/cors"
)

func main() {

	database.ConnectDb()

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	serverPort := "8000"
	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      corsHandler.Handler(routes.RegisterRoutes()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Server start channel for error
	serverErrors := make(chan error, 1)

	// Start HTTP server
	go func() {
		log.Printf("Starting server on port %s", serverPort)
		serverErrors <- srv.ListenAndServe()
	}()

	// Shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Block main and wait for shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, initiating shutdown", sig)

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown error: %v", err)
			if err := srv.Close(); err != nil {
				log.Printf("Failed to stop server: %v", err)
			}
		}
	}

	log.Println("Server stopped")
}
