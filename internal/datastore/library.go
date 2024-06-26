package datastore

import (
	"database/sql"
	"strings"

	"github.com/kadeallendev/bookstore/internal/models"
)

// A library datastore
type LibraryStore struct {
	db *sql.DB
}

// Create a new library store with the given database
func NewLibraryStore(db *sql.DB) *LibraryStore {
	libraryStore := &LibraryStore{
		db: db,
	}
	return libraryStore
}

// Find a book by isbn
func (ls *LibraryStore) FindBook(selected_isbn int) (*models.Book, error) {
	// Run query
	row := ls.db.QueryRow("SELECT isbn, title, edition_no, numofcop, numleft FROM book WHERE isbn = $1", selected_isbn)
	// Initialise values
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

func (ls *LibraryStore) AllBooks() ([]models.Book, error) {
	// Run query
	rows, err := ls.db.Query("SELECT * FROM book")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	// Iterate over all rows
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

		// Create new book and add to list
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
