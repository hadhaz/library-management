package domain

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ISBN        string    `json:"isbn"`
	AuthorID    int       `json:"author_id"`
	CategoryID  int       `json:"category_id"`
	PublishDate time.Time `json:"publish_date"`
	Description string    `json:"description"`
}

// BooksResponse represents a response containing a list of books
type BooksResponse struct {
	Books []Book `json:"books"`
}
