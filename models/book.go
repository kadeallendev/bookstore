package models

type Book struct {
	Isbn        int    `json:"isbn"`
	Title       string `json:"title"`
	Edition     int    `json:"edition"`
	TotalCopies int    `json:"totalCopies"`
	CopiesLeft  int    `json:"copiesLeft"`
}
