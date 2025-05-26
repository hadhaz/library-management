package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

func newPostgresDB(ctx context.Context, databaseURL string) (Database, error) {
	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}
	return &postgresDB{
		pool: dbpool,
	}, nil
}

type postgresDB struct {
	pool PostgresPool
}

func (db *postgresDB) GetBookByID(ctx context.Context, bookID int) (Book, error) {
	query := `
		SELECT id, title, isbn, author_id, category_id, stock, 
		       published_date, description, created_at, updated_at
		FROM books
		WHERE id = $1`

	var book Book
	err := db.pool.QueryRow(ctx, query, bookID).Scan(
		&book.ID,
		&book.Title,
		&book.ISBN,
		&book.AuthorID,
		&book.CategoryID,
		&book.Stock,
		&book.PublishedDate,
		&book.Description,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Book{}, fmt.Errorf("book with ID %d not found", bookID)
		}
		return Book{}, fmt.Errorf("unable to query book: %w", err)
	}

	return book, nil
}

func (db *postgresDB) LoadAllBooks(ctx context.Context) ([]Book, error) {
	query := `
		SELECT id, title, isbn, author_id, category_id, stock, 
		       published_date, description, created_at, updated_at
		FROM books`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query books table: %w", err)
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}
	return books, nil
}

func (db *postgresDB) CreateBook(ctx context.Context, newBook NewBook) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO books (title, isbn, author_id, category_id, stock, published_date, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		newBook.Title,
		newBook.ISBN,
		newBook.AuthorID,
		newBook.CategoryID,
		newBook.Stock,
		newBook.PublishedDate,
		newBook.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to insert book: %w", err)
	}
	return nil
}

func (db *postgresDB) UpdateBook(ctx context.Context, book Book) error {
	_, err := db.pool.Exec(ctx,
		`UPDATE books
		 SET title = $1,
		     isbn = $2,
		     author_id = $3,
		     category_id = $4,
		     stock = $5,
		     published_date = $6,
		     description = $7
		 WHERE id = $8`,
		book.Title,
		book.ISBN,
		book.AuthorID,
		book.CategoryID,
		book.Stock,
		book.PublishedDate,
		book.Description,
		book.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}
	return nil
}

func (db *postgresDB) DeleteBook(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

func (db *postgresDB) BorrowBook(ctx context.Context, book NewBorrowingRecord) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var stock int
	err = tx.QueryRow(ctx, "SELECT stock FROM books WHERE id = $1 FOR UPDATE", book.BookID).Scan(&stock)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("book not found")
		}
		return fmt.Errorf("failed to query book: %w", err)
	}

	if stock == 0 {
		return fmt.Errorf("book is not available")
	}

	_, err = tx.Exec(ctx, "UPDATE books SET stock = $1 WHERE id = $2", stock-1, book.BookID)
	if err != nil {
		return fmt.Errorf("failed to decrement stock: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO borrowing_records (user_id, book_id, borrowed_at, due_date)
		VALUES ($1, $2, $3, $4)`, book.UserID, book.BookID, book.BorrowedAt, book.BorrowedAt.Add(3*24*time.Hour))
	if err != nil {
		return fmt.Errorf("failed to insert borrowing record: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *postgresDB) ReturnBook(ctx context.Context, book BorrowingRecord) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM borrowing_records
			WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
		)
	`, book.UserID, book.BookID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check borrowing record: %w", err)
	}
	if !exists {
		return fmt.Errorf("borrowing record not found or already returned")
	}

	_, err = tx.Exec(ctx, `
		UPDATE borrowing_records
		SET returned_at = $1
		WHERE user_id = $2 AND book_id = $3 AND returned_at IS NULL
	`, time.Now(), book.UserID, book.BookID)
	if err != nil {
		return fmt.Errorf("failed to update borrowing record: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE books
		SET stock = stock + 1
		WHERE id = $1
	`, book.BookID)
	if err != nil {
		return fmt.Errorf("failed to increment book stock: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *postgresDB) GetRecommendedBooks(ctx context.Context, bookID int) ([]BookRecommendation, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, book_id, recommended_book_id, score  FROM book_recommendation WHERE book_id = $1", bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to query book recommendations: %w", err)
	}

	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[BookRecommendation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return books, nil
}

func (db *postgresDB) AddRecommendedBook(ctx context.Context, book NewBookRecommendation) error {
	_, err := db.pool.Exec(ctx, "INSERT INTO book_recommendation (book_id, recommended_book_id, score) VALUES ($1, $2, $3)", book.BookID, book.RecommendedBookID, book.Score)
	if err != nil {
		return fmt.Errorf("failed to add recommended book: %w", err)
	}

	return nil
}

func (db *postgresDB) CloseConnections() {
	db.pool.Close()
}
