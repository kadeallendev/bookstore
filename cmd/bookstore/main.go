package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kadeallendev/bookstore/config"
	"github.com/kadeallendev/bookstore/internal/datastore"
	"github.com/kadeallendev/bookstore/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("error loading config: ", err)
	}

	// Create library database object
	libraryDB, err := datastore.NewDB(*cfg)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}

	// Create library repository
	libraryStore := datastore.NewLibraryStore(libraryDB)

	// Create server
	srv := server.New(cfg, libraryStore)

	// Run server
	go func() {
		if err = srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("error running server: ", err)
		}
	}()

	// Create context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Wait for the context to be cancelled
	<-ctx.Done()

	// Shutdown server
	log.Println("server shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalln("error shutting down server: ", err)
	}

	// Shutdown database connection
	log.Println("database shutting down")
	if err := libraryDB.Close(); err != nil {
		log.Fatalln("error shutting down database: ", err)
	}
}
