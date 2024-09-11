package datastore

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kadeallendev/bookstore/internal/models"
)

type LibraryStore interface {
	GetBook(int) (*models.Book, error)
	GetAllBooks() ([]models.Book, error)
	// GetAllBooksOnLoan() ([]models.Book, error)
	// DeleteBook(int) error

	// BorrowBook(int) (*models.Book, error)
	// ReturnBook(int) error

	GetAuthor(int) (*models.Author, error)
	GetAllAuthors() ([]models.Author, error)
	GetBooksByAuthor(int) ([]models.Book, error)
	// DeleteAuthor(int) error

	GetCustomer(int) (*models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
	GetBooksBorrowedByCustomer(int) ([]models.Book, error)
	// DeleteCustomer(int) error
}

// A library datastore
type libraryStore struct {
	db *sql.DB
}

// Create a new library store with the given database
func NewLibraryStore(db *sql.DB) LibraryStore {
	ls := libraryStore{
		db: db,
	}
	return ls
}

// Find a book by isbn
func (ls libraryStore) GetBook(selected_isbn int) (*models.Book, error) {
	// Run query
	row := ls.db.QueryRow("SELECT isbn, title, edition_no, numofcop, numleft FROM book WHERE isbn = $1", selected_isbn)

	// Initialise values
	var (
		isbn     int
		title    string
		edition  sql.NullInt32
		numofcop int
		numleft  int
	)

	// Scan row
	err := row.Scan(&isbn, &title, &edition, &numofcop, &numleft)
	if err != nil {
		return nil, fmt.Errorf("FindBook: error scanning result: %w", err)
	}

	// Create new book
	return newBook(isbn, title, edition, numofcop, numleft), nil
}

func (ls libraryStore) GetAllBooks() ([]models.Book, error) {
	// Run query
	rows, err := ls.db.Query("SELECT * FROM book")
	if err != nil {
		return nil, fmt.Errorf("AllBooks: error executing query: %w", err)
	}
	defer rows.Close()

	var books []models.Book

	// Iterate over all rows
	for rows.Next() {
		// Initialise values
		var (
			isbn     int
			title    string
			edition  sql.NullInt32
			numofcop int
			numleft  int
		)

		// Scan values
		if err = rows.Scan(&isbn, &title, &edition, &numofcop, &numleft); err != nil {
			return nil, fmt.Errorf("AllBooks: error scanning row: %w", err)
		}

		// Create new book and add to list
		book := newBook(isbn, title, edition, numofcop, numleft)
		books = append(books, *book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("AllBooks: %w", err)
	}
	return books, nil

}

func (ls libraryStore) GetBooksByAuthor(authorID int) ([]models.Book, error) {
	rows, err := ls.db.Query(`SELECT isbn, title, edition_no, numofcop, numleft
        FROM book b
        NATURAL JOIN book_author ba
        NATURAL JOIN author a
        WHERE a.authorid = $1
        ORDER BY authorseqno`, authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a nil slice of books, and no error
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var (
			isbn     int
			title    string
			edition  sql.NullInt32
			numofcop int
			numleft  int
		)
		err := rows.Scan(&isbn, &title, &edition, &numofcop, &numleft)
		if err != nil {
			return nil, fmt.Errorf("FindBooksByAuthor: error scanning row: %w", err)
		}

		book := newBook(isbn, title, edition, numofcop, numleft)
		books = append(books, *book)
	}

	return books, nil
}

func (ls libraryStore) GetAuthor(authorID int) (*models.Author, error) {
	row := ls.db.QueryRow("SELECT authorid, surname, name FROM author WHERE authorid = $1", authorID)

	// Initialise values
	var (
		id        int
		lastname  string
		firstname sql.NullString
	)
	// Scan values
	err := row.Scan(&id, &lastname, &firstname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("FindAuthor: error scanning row: %w", err)
	}

	return newAuthor(id, lastname, firstname), nil
}

func (ls libraryStore) GetAllAuthors() ([]models.Author, error) {
	rows, err := ls.db.Query("SELECT authorid, surname, name FROM author")
	if err != nil {
		return nil, fmt.Errorf("AllAuthors: error executing query: %w", err)
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var (
			id        int
			lastname  string
			firstname sql.NullString
		)
		err = rows.Scan(&id, &lastname, &firstname)
		if err != nil {
			return nil, fmt.Errorf("AllAuthors: there was an error scanning rows: %w", err)
		}

		author := newAuthor(id, lastname, firstname)
		authors = append(authors, *author)
	}

	return authors, nil
}

func (ls libraryStore) GetCustomer(customerID int) (*models.Customer, error) {
	row := ls.db.QueryRow("SELECT customerid, l_name, f_name, city FROM customer WHERE customerid = $1", customerID)

	var (
		id        int
		lastname  string
		firstname sql.NullString
		city      sql.NullString
	)
	err := row.Scan(&id, &lastname, &firstname, &city)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("FindCustomer: the customer with id %d does not exist: %w", err)
		}
		return nil, err
	}

	return newCustomer(id, lastname, firstname, city), nil
}

func (ls libraryStore) GetAllCustomers() ([]models.Customer, error) {
	rows, err := ls.db.Query("SELECT * FROM customer")
	if err != nil {
		return nil, fmt.Errorf("GetAllCustomers: error executing query: %w", err)
	}

	var customers []models.Customer
	for rows.Next() {
		var (
			id        int
			lastname  string
			firstname sql.NullString
			city      sql.NullString
		)
		err = rows.Scan(&id, &lastname, &firstname, &city)
		if err != nil {
			return nil, fmt.Errorf("GetAllCustomers: error scanning row: %w", err)
		}

		customer := newCustomer(id, lastname, firstname, city)
		customers = append(customers, *customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllCustomers: error iterating over rows: %w", err)
	}

	return customers, nil
}

func (ls libraryStore) GetBooksBorrowedByCustomer(customerID int) ([]models.Book, error) {
	rows, err := ls.db.Query(`SELECT b.isbn, b.title, b.edition_no, b.numofcop, b.numleft
        FROM book b
        NATURAL JOIN cust_book cb
        WHERE c.customerid = $1`, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("FindBooksBorrowedByCustomer: error executing query: %w", err)
	}

	var books []models.Book
	for rows.Next() {
		var (
			isbn     int
			title    string
			edition  sql.NullInt32
			numofcop int
			numleft  int
		)
		err = rows.Scan(&isbn, &title, &edition, &numofcop, &numleft)
		if err != nil {
			return nil, fmt.Errorf("FindBooksBorrowedByCustomer: error scanning row: %w", err)
		}

		book := newBook(isbn, title, edition, numofcop, numleft)
		books = append(books, *book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FindBooksBorrowedByCustomer: error iterating over rows: %w", err)
	}

	return books, nil
}

// Helper functions

func newBook(isbn int, title string, edition sql.NullInt32, totalCopies int, copiesLeft int) *models.Book {
	book := models.Book{
		Isbn:  isbn,
		Title: strings.TrimSpace(title),
		// Edition may be null, so leave empty here
		TotalCopies: totalCopies,
		CopiesLeft:  copiesLeft,
	}

	if edition.Valid {
		book.Edition = int(edition.Int32)
	}

	return &book
}

func newCustomer(id int, lastname string, firstname, city sql.NullString) *models.Customer {
	customer := models.Customer{
		ID:       id,
		Lastname: lastname,
		// Ignore firstname and city as may be null
	}

	if firstname.Valid {
		customer.Firstname = firstname.String
	}
	if city.Valid {
		customer.City = city.String
	}

	return &customer
}

func newAuthor(id int, lastname string, firstname sql.NullString) *models.Author {
	author := models.Author{
		ID:       id,
		Lastname: lastname,
		// ignore firstname as may be null
	}

	if firstname.Valid {
		author.Firstname = firstname.String
	}

	return &author
}
