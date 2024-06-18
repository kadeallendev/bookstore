package server

import (
	"bookstore/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Add all routes of application
func addRoutes(
	mux *http.ServeMux,
	libraryService *database.LibraryService,
) {
	mux.Handle("GET /", handleIndex())
	mux.Handle("GET /book", handleAllBooks(libraryService))
	mux.Handle("GET /book/{isbn}", handleFindBook(libraryService))
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
func handleFindBook(libraryService *database.LibraryService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Find book: ", r.Context())

			// Set header
			w.Header().Set("Content-Type", "application/json")

			// Parse the isbn
			isbn, err := strconv.Atoi(r.PathValue("isbn"))
			if err != nil {
				log.Println("invalid isbn value", err)
				http.Error(w, "Invalid isbn value", http.StatusBadRequest)
				return
			}

			// Find the book
			book, err := libraryService.FindBook(isbn)
			if err != nil {
				// If the book was not found
				if err == sql.ErrNoRows {
					log.Println("failed to find book: ", err)
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(ErrorResponse{
						Error:   "Not Found",
						Message: "The book with the given isbn was not found",
						Status:  http.StatusNotFound,
					})
					return
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Encode book
			encodedBook, err := json.Marshal(book)
			if err != nil {
				log.Println("failed to serialize book", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Write to the response
			bytesWritten, err := w.Write(encodedBook)
			if err != nil {
				log.Println("failed to write to response", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			log.Printf("Find Book: %s\n\tWrote %d bytes:\n%s\n", r.Context(), bytesWritten, string(encodedBook))
		},
	)
}

// Handle the allbooks route
func handleAllBooks(libraryService *database.LibraryService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("All Books: ", r.Context())
			// Set header
			w.Header().Set("Content-Type", "application/json")

			// Get all books
			books, err := libraryService.AllBooks()
			if err != nil {
				log.Println("failed to retreive all books", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Encode books
			encodedBooks, err := json.Marshal(books)
			if err != nil {
				log.Println("failed to encode books", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Write to response
			bytesWritten, err := w.Write(encodedBooks)
			if err != nil {
				log.Println("failed to write to repsonse", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			log.Printf("All Books: %s\n\tWrote %d bytes:\n%s\n", r.Context(), bytesWritten, string(encodedBooks))
		},
	)
}

// func (s *Server) RegisterRoutes() http.Handler {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("GET /", indexHandler)
// 	mux.HandleFunc("GET /books/{isbn}", s.bookLookup)
// 	mux.HandleFunc("GET /foo", fooHandler)
//
// 	return mux
// }
//
// func (s *Server) bookLookup(w http.ResponseWriter, r *http.Request) {
//
// }
//
// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	resp := make(map[string]string)
// 	resp["message"] = "Hello, world"
//
// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		log.Fatalf("error handling JSON marshal. Err: %v", err)
// 	}
// 	_, _ = w.Write(jsonResp)
// }
//
// func fooHandler(w http.ResponseWriter, r *http.Request) {
// 	resp := make(map[string]string)
// 	resp["message"] = "Foo"
// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		log.Fatalf("error handling JSON marshal. Err: %v", err)
// 	}
//
// 	_, _ = w.Write(jsonResp)
// }
