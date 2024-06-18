package database

import (
	"bookstore/config"
	"bookstore/models"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type DatabaseService interface {
	Health() map[string]string
	Close() error
}

type LibraryService struct {
	db *sql.DB
}

func NewLibraryService(cfg *config.Config) *LibraryService {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	libraryService := &LibraryService{
		db: db,
	}
	return libraryService
}

func (s *LibraryService) FindBook(selected_isbn int) (*models.Book, error) {
	// Run query
	row := s.db.QueryRow("SELECT isbn, title, edition_no, numofcop, numleft FROM book WHERE isbn = $1", selected_isbn)
	// Initilaise values
	var (
		isbn         int
		title        string
		edition_null sql.NullInt32
		numofcop     int
		numleft      int
	)

	// Scan row
	err := row.Scan(&isbn, &title, &edition_null, &numofcop, &numleft)
	if err != nil {
		return nil, err
	}
	// Clean data
	title = strings.TrimSpace(title)
	edition := -1
	if edition_null.Valid {
		edition = int(edition_null.Int32)
	}

	var book = &models.Book{
		Isbn:        isbn,
		Title:       title,
		Edition:     edition,
		TotalCopies: numofcop,
		CopiesLeft:  numleft,
	}
	return book, nil
}

func (s *LibraryService) AllBooks() ([]models.Book, error) {
	// Run query
	rows, err := s.db.Query("SELECT * FROM book")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over all books
	var books []models.Book
	for rows.Next() {
		// Initialise values
		var (
			isbn         int
			title        string
			edition_null sql.NullInt32
			numofcop     int
			numleft      int
		)

		// Scan values
		if err = rows.Scan(&isbn, &title, &edition_null, &numofcop, &numleft); err != nil {
			return nil, err
		}

		// Clean data
		title = strings.TrimSpace(title)
		edition := -1
		if edition_null.Valid {
			edition = int(edition_null.Int32)
		}

		book := models.Book{
			Isbn:        isbn,
			Title:       title,
			Edition:     edition,
			TotalCopies: numofcop,
			CopiesLeft:  numleft,
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return books, err
	}
	return books, nil

}

func (s *LibraryService) Health() map[string]string {
	// TODO: Implement properly

	stats := make(map[string]string)
	stats["message"] = "Health Check"
	return stats
}

func (s *LibraryService) Close() error {
	log.Println("Disconnected from database")
	return s.db.Close()
}
