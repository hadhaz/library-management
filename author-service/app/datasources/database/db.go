package database

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

type Author struct {
	ID          int        `db:"id"`
	FirstName   string     `db:"first_name"`
	LastName    string     `db:"last_name"`
	BirthDate   *time.Time `db:"birth_date"`
	Nationality string     `db:"nationality"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type NewAuthor struct {
	FirstName   string     `db:"first_name"`
	LastName    string     `db:"last_name"`
	BirthDate   *time.Time `db:"birth_date"`
	Nationality string     `db:"nationality"`
}

type Database interface {
	AddAuthor(ctx context.Context, author NewAuthor) error
	UpdateAuthor(ctx context.Context, author Author) error
	DeleteAuthor(ctx context.Context, id int) error
	GetAuthor(ctx context.Context, id int) (Author, error)
	ListAuthor(ctx context.Context, filter Author, limit, offset int) ([]Author, error)

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
