package database

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

// Book represents a book in the database
type Book struct {
	ID            int       `db:"id"`
	Title         string    `db:"title"`
	ISBN          string    `db:"isbn"`
	AuthorID      int       `db:"author_id"`
	CategoryID    int       `db:"category_id"`
	Stock         int       `db:"stock"`
	PublishedDate time.Time `db:"published_date"`
	Description   string    `db:"description"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// NewBook represents a new book to be created to the database
type NewBook struct {
	Title         string
	ISBN          string
	AuthorID      int
	CategoryID    int
	Stock         int
	PublishedDate time.Time
	Description   string
}

type BorrowingRecord struct {
	ID         int
	BookID     int
	UserID     int
	BorrowedAt time.Time
	ReturnedAt time.Time
	DueDate    time.Time
}

type NewBorrowingRecord struct {
	BookID     int
	UserID     int
	BorrowedAt time.Time
	ReturnedAt time.Time
	Status     string
}

type BookRecommendation struct {
	ID                int
	BookID            int
	RecommendedBookID int
	Score             float32
}

type NewBookRecommendation struct {
	BookID            int
	RecommendedBookID int
	Score             float32
}

type Database interface {
	LoadAllBooks(ctx context.Context) ([]Book, error)

	GetBookByID(ctx context.Context, bookID int) (Book, error)

	CreateBook(ctx context.Context, newBook NewBook) error

	UpdateBook(ctx context.Context, book Book) error

	DeleteBook(ctx context.Context, id int) error

	BorrowBook(ctx context.Context, book NewBorrowingRecord) error

	ReturnBook(ctx context.Context, book BorrowingRecord) error

	GetRecommendedBooks(ctx context.Context, bookID int) ([]BookRecommendation, error)

	AddRecommendedBook(ctx context.Context, book NewBookRecommendation) error

	CloseConnections()
}

// NewDatabase creates a new Database instance
func NewDatabase(ctx context.Context, databaseURL string) (Database, error) {
	if databaseURL == "" {
		slog.Info("Using in-memory database implementation")
		return newMemoryDB(), nil
	}

	if strings.HasPrefix(databaseURL, "postgres://") {
		db, err := newPostgresDB(ctx, databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL database connection: %w", err)
		}
		slog.Info("Using PostgreSQL database implementation")
		return db, nil
	}

	return nil, fmt.Errorf("unsupported database URL scheme: %s", databaseURL)
}

func EscapeQuery(query string) string {
	return regexp.QuoteMeta(query)
}
