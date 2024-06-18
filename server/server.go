package server

import (
	"net/http"

	"bookstore/database"
)

type Server struct {
	libraryService database.LibraryService
}

// New server
func New(libraryService *database.LibraryService) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, libraryService)
	var handler http.Handler = mux
	// Add middleware
	return handler
}
