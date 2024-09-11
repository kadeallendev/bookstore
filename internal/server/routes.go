package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kadeallendev/bookstore/internal/datastore"
)

// Add all routes of application
func addRoutes(
	mux *http.ServeMux,
	libraryStore *datastore.LibraryStore,
) {
	mux.Handle("GET /", handleIndex())
	mux.Handle("GET /book", handleAllBooks(libraryStore))
	mux.Handle("GET /book/{isbn}", handleFindBook(libraryStore))
}

func handleIndex() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Index: ", r.Context())
			_, _ = w.Write([]byte("Hello, world!"))
		},
	)
}

// Handle the index route
func handleFindBook(libraryStore *datastore.LibraryStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("handling find book at %s\n", r.URL.Path)

			// Parse the isbn
			isbn, err := strconv.Atoi(r.PathValue("isbn"))
			if err != nil {
				log.Println("invalid isbn value", err)
				RespondWithError(w, "invalid isbn value", http.StatusBadRequest)
				return
			}

			// Find the book
			book, err := libraryStore.GetBook(isbn)
			if err != nil {
				// If the book was not found
				if err == sql.ErrNoRows {
					log.Println("failed to find book: ", err)
					errMsg := fmt.Sprintf("The book %d was not found", isbn)
					RespondWithError(w, errMsg, http.StatusNotFound)
					return
				}
				RespondWithError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Write to the response
			if err = WriteJSON(w, book, http.StatusOK); err != nil {
				log.Println("failed to write to response", err)
				RespondWithError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("responded to find book request at: %s\n", r.URL.Path)
		},
	)
}

// Handle the allbooks route
func handleAllBooks(libraryStore *datastore.LibraryStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("handling all books at %s\n", r.URL.Path)

			// Get all books
			books, err := libraryStore.GetAllBooks()
			if err != nil {
				log.Println("failed to retreive all books", err)
				RespondWithError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Write to response
			if err = WriteJSON(w, books, http.StatusOK); err != nil {
				log.Println("failed to write response", err)
				RespondWithError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			log.Printf("responded to all books request at %s\n", r.URL.Path)
		},
	)
}
