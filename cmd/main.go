package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Rafli-Dewanto/go-template/internal/config"
	"github.com/Rafli-Dewanto/go-template/internal/router"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const serverAddr = ":8080"

func main() {
	// Load database configuration
	configPath := filepath.Join("config", "database.ini")
	dbConfig, err := config.LoadDatabaseConfig(configPath)
	if err != nil {
		log.Fatalf("cannot load database config: %v", err)
	}

	// Connect to database using the configuration
	db, err := sqlx.Connect(dbConfig.Driver, dbConfig.GetDSN())
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer db.Close()

	// Create a channel to listen for interrupt signal to gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	router := router.NewRouter(db)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router.SetupRoutes(),
	}

	go func() {
		log.Printf("Server is running on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-sigChan
	fmt.Println("\nShutting down server...")

	// Gracefully shutdown the server
	if err := server.Close(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}
