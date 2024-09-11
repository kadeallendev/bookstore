package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/kadeallendev/bookstore/config"
	"github.com/kadeallendev/bookstore/internal/datastore"
)

// Could add the handlers as methods to the Server struct, instead of passing librarystore to all routes

type Server struct {
	httpServer *http.Server

	libraryStore *datastore.LibraryStore
}

// Create a new server with the given configuration
func New(cfg *config.Config, libraryStore *datastore.LibraryStore) *Server {

	router := createRouter(libraryStore)
	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.ServerAddress, cfg.ServerPort),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		httpServer:   srv,
		libraryStore: libraryStore,
	}

}

// Run the server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func createRouter(libraryStore *datastore.LibraryStore) http.Handler {
	// TODO: Add routes
}

// rem: func NewServer()....datastore.LibraryStore
// // Add routes
// mux := http.NewServeMux()
// addRoutes(mux, libraryStore)
// var handler http.Handler = mux
//
// // TODO: Add middleware
//
// // Create server object
// httpServer := &http.Server{
// 	Addr:         net.JoinHostPort(cfg.ServerAddress, cfg.ServerPort),
// 	Handler:      handler,
// 	IdleTimeout:  time.Minute,
// 	ReadTimeout:  10 * time.Second,
// 	WriteTimeout: 30 * time.Second,
// }
//
// // Return new server object
// return &Server{
// 	httpServer:   httpServer,
// 	libraryStore: libraryStore,
// }
