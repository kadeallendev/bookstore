package datastore

import (
	"database/sql"
	"strings"

	"github.com/kadeallendev/bookstore/internal/models"
)

type LibraryStore struct {
	db *sql.DB
}

func NewLibraryStore(db *sql.DB) *LibraryStore {
	libraryRepo := &LibraryStore{
		db: db,
	}
	return libraryRepo
}

func (ls *LibraryStore) FindBook(selected_isbn int) (*models.Book, error) {
	// Run query
	row := ls.db.QueryRow("SELECT isbn, title, edition_no, numofcop, numleft FROM book WHERE isbn = $1", selected_isbn)
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
