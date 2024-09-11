package models

type Book struct {
	Isbn        int    `json:"isbn"`
	Title       string `json:"title"`
	Edition     int    `json:"edition"`
	TotalCopies int    `json:"totalCopies"`
	CopiesLeft  int    `json:"copiesLeft"`
}

type BookWithAuthors struct {
	Book
	Authors []Author
}

type Customer struct {
	ID        int    `json:"customer_id"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	City      string `json:"city"`
}

type CustomerWithBooks struct {
	Customer
	Books []Book
}

type Author struct {
	ID        int    `json:"author_id"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}

type AuthorWithBooks struct {
	Author
	Books []Book
}

type LoanedBook struct {
	Book
	Customers []Customer
	Authors   []Author
}
