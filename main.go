package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"bookstore/config"
	"bookstore/database"
	"bookstore/server"

	"github.com/joho/godotenv"
)

func main() {

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("error loading config", err)
	}

	if err := loadEnv(); err != nil {
		log.Fatal("error loading .env file:", err)
	}

	// Create library service
	libraryService := database.NewLibraryService(cfg)
	// Create server
	srv := server.New(libraryService)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:      srv,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// TODO: Put in goroutine
	log.Printf("listening on %s\n", httpServer.Addr)
	// Serve application
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}

}

func loadConfig() (*config.Config, error) {
	var config config.Config
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")
	config.DBUsername = os.Getenv("DB_USERNAME")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")
	config.DBHost = os.Getenv("DB_HOST")
	config.DBPort = os.Getenv("DB_PORT")
	return &config, nil
}

func loadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	return nil
}

// var wg sync.WaitGroup
// wg.Add(1)
// go func() {
// 	defer wg.Done()
// 	<-ctx.Done()
// 	// make a new context for the shutdown
// 	shutdownCtx := context.Background()
// 	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()
// 	if err := httpServer.Shutdown(shutdownCtx); err != nil {
//
// 		fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
// 	}
// }()
// wg.Wait()
