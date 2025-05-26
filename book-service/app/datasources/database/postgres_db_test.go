package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDB_GetBookByID_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer mockPool.Close()

	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	query := `
		SELECT id, title, isbn, author_id, category_id, stock, 
		       published_date, description, created_at, updated_at
		FROM books
		WHERE id = $1`
	headerRow := []string{"id", "title", "isbn", "author_id", "category_id", "stock", "published_date",
		"description", "created_at", "updated_at"}

	mockPool.ExpectQuery(EscapeQuery(query)).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows(headerRow).
			AddRow(1, "book1", "1234567890", 1, 2, 10,
				fixedTime, "a book desc", fixedTime, fixedTime))

	db := &postgresDB{
		pool: mockPool,
	}
	result, err := db.GetBookByID(context.Background(), 1)

	assert.Nil(t, err)
	assertBook(t, result, 1, NewBook{
		Title:         "book1",
		ISBN:          "1234567890",
		AuthorID:      1,
		CategoryID:    2,
		Stock:         10,
		PublishedDate: fixedTime,
		Description:   "a book desc",
	})
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_GetBookByID_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer mockPool.Close()

	query := `
		SELECT id, title, isbn, author_id, category_id, stock, 
		       published_date, description, created_at, updated_at
		FROM books
		WHERE id = $1`

	mockPool.ExpectQuery(EscapeQuery(query)).
		WithArgs(999).
		WillReturnError(pgx.ErrNoRows)

	db := &postgresDB{
		pool: mockPool,
	}

	result, err := db.GetBookByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, Book{}, result)

	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_GetBooks_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	defer mockPool.Close()

	mockPool.ExpectQuery(EscapeQuery(`
	SELECT id, title, isbn, author_id, category_id, stock, 
	       published_date, description, created_at, updated_at
	FROM books`)).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "title", "isbn", "author_id", "category_id", "stock",
			"published_date", "description", "created_at", "updated_at"}).
			AddRow(1, "book1", "1234567890", 1, 2, 10,
				fixedTime, "a book desc", time.Now(), time.Now()))

	db := postgresDB{
		pool: mockPool,
	}
	result, err := db.LoadAllBooks(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assertBook(t, result[0], 1,
		NewBook{
			Title:         "book1",
			AuthorID:      1,
			CategoryID:    2,
			Stock:         10,
			PublishedDate: fixedTime,
			Description:   "a book desc",
			ISBN:          "1234567890",
		})

	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_GetBooks_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)

	mockPool.ExpectQuery(EscapeQuery(`
	SELECT id, title, isbn, author_id, category_id, stock, 
	       published_date, description, created_at, updated_at
	FROM books`)).
		WillReturnError(assert.AnError)

	db := postgresDB{
		pool: mockPool,
	}
	result, err := db.LoadAllBooks(context.Background())

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "failed to query books table")
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_CreateBook_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	mockPool.ExpectExec(EscapeQuery(`INSERT INTO books (title, isbn, author_id, category_id, stock, published_date, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`)).
		WithArgs("book1", "1234567890", 1, 2, 10, fixedTime, "a book desc").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	db := postgresDB{
		pool: mockPool,
	}
	err = db.CreateBook(context.Background(), NewBook{
		Title:         "book1",
		ISBN:          "1234567890",
		AuthorID:      1,
		CategoryID:    2,
		Stock:         10,
		PublishedDate: fixedTime,
		Description:   "a book desc",
	})

	assert.Nil(t, err)
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_CreateBook_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	mockPool.ExpectExec(EscapeQuery(`INSERT INTO books (title, isbn, author_id, category_id, stock, published_date, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`)).
		WithArgs("book1", "1234567890", 1, 2, 10, fixedTime, "a book desc").
		WillReturnError(assert.AnError)

	db := postgresDB{
		pool: mockPool,
	}
	err = db.CreateBook(context.Background(), NewBook{
		Title:         "book1",
		ISBN:          "1234567890",
		AuthorID:      1,
		CategoryID:    2,
		Stock:         10,
		PublishedDate: fixedTime,
		Description:   "a book desc",
	})

	assert.ErrorContains(t, err, "failed to insert book")
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_UpdateBook_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	mockPool.ExpectExec(EscapeQuery(`UPDATE books
		 SET title = $1,
		     isbn = $2,
		     author_id = $3,
		     category_id = $4,
		     stock = $5,
		     published_date = $6,
		     description = $7
		 WHERE id = $8`)).
		WithArgs("book1", "1234567890", 1, 2, 10, fixedTime, "a book desc", 21).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	db := postgresDB{
		pool: mockPool,
	}
	err = db.UpdateBook(context.Background(), Book{
		Title:         "book1",
		ISBN:          "1234567890",
		AuthorID:      1,
		CategoryID:    2,
		Stock:         10,
		PublishedDate: fixedTime,
		Description:   "a book desc",
		ID:            21,
	})

	assert.Nil(t, err)
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_UpdateBook_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	mockPool.ExpectExec(EscapeQuery(`UPDATE books
		 SET title = $1,
		     isbn = $2,
		     author_id = $3,
		     category_id = $4,
		     stock = $5,
		     published_date = $6,
		     description = $7
		 WHERE id = $8`)).
		WithArgs("book1", "1234567890", 1, 2, 10, fixedTime, "a book desc", 21).
		WillReturnError(assert.AnError)

	db := postgresDB{
		pool: mockPool,
	}
	err = db.UpdateBook(context.Background(), Book{
		Title:         "book1",
		ISBN:          "1234567890",
		AuthorID:      1,
		CategoryID:    2,
		Stock:         10,
		PublishedDate: fixedTime,
		Description:   "a book desc",
		ID:            21,
	})

	assert.ErrorContains(t, err, "failed to update book")
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_DeleteBook_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)

	mockPool.ExpectExec(EscapeQuery(`DELETE FROM books WHERE id = $1`)).
		WithArgs(21).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	db := postgresDB{
		pool: mockPool,
	}
	err = db.DeleteBook(context.Background(), 21)

	assert.Nil(t, err)
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_DeleteBook_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)

	mockPool.ExpectExec(EscapeQuery(`DELETE FROM books WHERE id = $1`)).
		WithArgs(21).
		WillReturnError(assert.AnError)

	db := postgresDB{
		pool: mockPool,
	}
	err = db.DeleteBook(context.Background(), 21)

	assert.ErrorContains(t, err, "failed to delete book")
	assert.Nil(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_BorrowBook_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	ctx := context.Background()
	bookID := 1
	userID := 123
	borrowedAt := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)
	dueDate := borrowedAt.Add(3 * 24 * time.Hour)

	mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	mockPool.ExpectQuery("SELECT stock FROM books WHERE id = \\$1 FOR UPDATE").
		WithArgs(bookID).
		WillReturnRows(pgxmock.NewRows([]string{"stock"}).AddRow(5))

	mockPool.ExpectExec("UPDATE books SET stock = \\$1 WHERE id = \\$2").
		WithArgs(4, bookID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mockPool.ExpectExec("INSERT INTO borrowing_records").
		WithArgs(userID, bookID, borrowedAt, dueDate).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mockPool.ExpectCommit()

	db := postgresDB{pool: mockPool}
	err = db.BorrowBook(ctx, NewBorrowingRecord{
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: borrowedAt,
	})
	assert.NoError(t, err)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_BorrowBook_Fail(t *testing.T) {
	ctx := context.Background()
	userID := 123
	bookID := 456
	borrowedAt := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)
	dueDate := borrowedAt.Add(3 * 24 * time.Hour)

	book := NewBorrowingRecord{
		UserID:     userID,
		BookID:     bookID,
		BorrowedAt: borrowedAt,
	}

	t.Run("fail to begin transaction", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).WillReturnError(errors.New("begin error"))

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "failed to start transaction")
	})

	t.Run("book not found", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`SELECT stock FROM books WHERE id = $1 FOR UPDATE`)).
			WithArgs(bookID).
			WillReturnError(pgx.ErrNoRows)

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "book not found")
	})

	t.Run("failed to query book", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`SELECT stock FROM books WHERE id = $1 FOR UPDATE`)).
			WithArgs(bookID).
			WillReturnError(errors.New("query error"))

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "failed to query book")
	})

	t.Run("book out of stock", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`SELECT stock FROM books WHERE id = $1 FOR UPDATE`)).
			WithArgs(bookID).
			WillReturnRows(pgxmock.NewRows([]string{"stock"}).AddRow(0))

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "book is not available")
	})

	t.Run("fail to decrement stock", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`SELECT stock FROM books WHERE id = $1 FOR UPDATE`)).
			WithArgs(bookID).
			WillReturnRows(pgxmock.NewRows([]string{"stock"}).AddRow(1))
		mockPool.ExpectExec(EscapeQuery(`UPDATE books SET stock = $1 WHERE id = $2`)).
			WithArgs(0, bookID).
			WillReturnError(errors.New("update stock failed"))

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "failed to decrement stock")
	})

	t.Run("fail to insert borrowing record", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`SELECT stock FROM books WHERE id = $1 FOR UPDATE`)).
			WithArgs(bookID).
			WillReturnRows(pgxmock.NewRows([]string{"stock"}).AddRow(1))
		mockPool.ExpectExec(EscapeQuery(`UPDATE books SET stock = $1 WHERE id = $2`)).
			WithArgs(0, bookID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockPool.ExpectExec(EscapeQuery(`INSERT INTO borrowing_records`)).
			WithArgs(userID, bookID, borrowedAt, dueDate).
			WillReturnError(errors.New("insert fail"))

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "failed to insert borrowing record")
	})

	t.Run("fail to commit transaction", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`SELECT stock FROM books WHERE id = $1 FOR UPDATE`)).
			WithArgs(bookID).
			WillReturnRows(pgxmock.NewRows([]string{"stock"}).AddRow(1))
		mockPool.ExpectExec(EscapeQuery(`UPDATE books SET stock = $1 WHERE id = $2`)).
			WithArgs(0, bookID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockPool.ExpectExec(EscapeQuery(`INSERT INTO borrowing_records`)).
			WithArgs(userID, bookID, borrowedAt, dueDate).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mockPool.ExpectCommit().WillReturnError(errors.New("commit error"))

		db := &postgresDB{pool: mockPool}
		err = db.BorrowBook(ctx, book)

		assert.ErrorContains(t, err, "failed to commit transaction")
	})
}

func TestPostgresDB_ReturnBook_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	db := &postgresDB{pool: mockPool}
	ctx := context.Background()
	userID := 1
	bookID := 101

	record := BorrowingRecord{
		UserID: userID,
		BookID: bookID,
	}

	mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})

	mockPool.ExpectQuery(EscapeQuery(`
		SELECT EXISTS (
			SELECT 1 FROM borrowing_records
			WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
		)`)).
		WithArgs(userID, bookID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	mockPool.ExpectExec(EscapeQuery(`
		UPDATE borrowing_records
		SET returned_at = $1
		WHERE user_id = $2 AND book_id = $3 AND returned_at IS NULL
	`)).
		WithArgs(pgxmock.AnyArg(), userID, bookID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mockPool.ExpectExec(EscapeQuery(`
		UPDATE books
		SET stock = stock + 1
		WHERE id = $1
	`)).
		WithArgs(bookID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mockPool.ExpectCommit()

	err = db.ReturnBook(ctx, record)
	assert.NoError(t, err)
	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_ReturnBook_Fail(t *testing.T) {
	ctx := context.Background()
	userID := 123
	bookID := 456

	book := BorrowingRecord{
		UserID: userID,
		BookID: bookID,
	}

	t.Run("fail to begin transaction", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).WillReturnError(errors.New("begin error"))

		db := &postgresDB{pool: mockPool}
		err = db.ReturnBook(ctx, book)

		assert.ErrorContains(t, err, "failed to start transaction")
	})

	t.Run("fail to check borrowing record existence", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`
			SELECT EXISTS (
				SELECT 1 FROM borrowing_records
				WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
			)
		`)).
			WithArgs(userID, bookID).
			WillReturnError(errors.New("query error"))

		db := &postgresDB{pool: mockPool}
		err = db.ReturnBook(ctx, book)

		assert.ErrorContains(t, err, "failed to check borrowing record")
	})

	t.Run("borrowing record not found or already returned", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`
			SELECT EXISTS (
				SELECT 1 FROM borrowing_records
				WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
			)
		`)).
			WithArgs(userID, bookID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		db := &postgresDB{pool: mockPool}
		err = db.ReturnBook(ctx, book)

		assert.ErrorContains(t, err, "borrowing record not found or already returned")
	})

	t.Run("fail to update borrowing record", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`
			SELECT EXISTS (
				SELECT 1 FROM borrowing_records
				WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
			)
		`)).
			WithArgs(userID, bookID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mockPool.ExpectExec(EscapeQuery(`
			UPDATE borrowing_records
			SET returned_at = $1
			WHERE user_id = $2 AND book_id = $3 AND returned_at IS NULL
		`)).
			WithArgs(mock.Anything, userID, bookID).
			WillReturnError(errors.New("update error"))

		db := &postgresDB{pool: mockPool}
		err = db.ReturnBook(ctx, book)

		assert.ErrorContains(t, err, "failed to update borrowing record")
	})

	t.Run("fail to increment book stock", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`
			SELECT EXISTS (
				SELECT 1 FROM borrowing_records
				WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
			)
		`)).
			WithArgs(userID, bookID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mockPool.ExpectExec(EscapeQuery(`
			UPDATE borrowing_records
			SET returned_at = $1
			WHERE user_id = $2 AND book_id = $3 AND returned_at IS NULL
		`)).
			WithArgs(pgxmock.AnyArg(), userID, bookID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mockPool.ExpectExec(EscapeQuery(`
			UPDATE books
			SET stock = stock + 1
			WHERE id = $1
		`)).
			WithArgs(bookID).
			WillReturnError(fmt.Errorf("failed to increment book stock"))

		db := &postgresDB{pool: mockPool}
		err = db.ReturnBook(ctx, book)

		assert.ErrorContains(t, err, "failed to increment book stock")
	})

	t.Run("fail to commit transaction", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockPool.Close()

		mockPool.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
		mockPool.ExpectQuery(EscapeQuery(`
			SELECT EXISTS (
				SELECT 1 FROM borrowing_records
				WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
			)
		`)).
			WithArgs(userID, bookID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mockPool.ExpectExec(EscapeQuery(`
			UPDATE borrowing_records
			SET returned_at = $1
			WHERE user_id = $2 AND book_id = $3 AND returned_at IS NULL
		`)).
			WithArgs(pgxmock.AnyArg(), userID, bookID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mockPool.ExpectExec(EscapeQuery(`
			UPDATE books
			SET stock = stock + 1
			WHERE id = $1
		`)).
			WithArgs(bookID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mockPool.ExpectCommit().WillReturnError(errors.New("commit error"))

		db := &postgresDB{pool: mockPool}
		err = db.ReturnBook(ctx, book)

		assert.ErrorContains(t, err, "failed to commit transaction")
	})
}

func TestPostgresDB_GetRecommendedBooks_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	db := &postgresDB{pool: mockPool}
	ctx := context.Background()
	bookID := 1

	mockPool.ExpectQuery(EscapeQuery(`
			SELECT id, book_id, recommended_book_id, score  FROM book_recommendation WHERE book_id = $1`)).
		WithArgs(bookID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "book_id", "recommended_book_id", "score"}).
			AddRow(bookID, 1, 2, 0.9))

	books, err := db.GetRecommendedBooks(ctx, bookID)

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_GetRecommendedBooks_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	db := &postgresDB{pool: mockPool}
	ctx := context.Background()
	bookID := 1

	mockPool.ExpectQuery(EscapeQuery(`
			SELECT id, book_id, recommended_book_id, score  FROM book_recommendation WHERE book_id = $1`)).
		WithArgs(bookID).
		WillReturnError(fmt.Errorf("failed to query recommended_books"))
	books, err := db.GetRecommendedBooks(ctx, bookID)

	assert.ErrorContains(t, err, "failed to query recommended_books")
	assert.Nil(t, books)
}

func TestPostgresDB_AddRecommendedBook_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	db := &postgresDB{pool: mockPool}
	ctx := context.Background()
	bookID := 1
	bookRecommendationId := 2
	bookRecommendation := NewBookRecommendation{
		BookID:            bookID,
		RecommendedBookID: bookRecommendationId,
		Score:             0.85,
	}

	mockPool.ExpectExec(EscapeQuery("INSERT INTO book_recommendation (book_id, recommended_book_id, score) VALUES ($1, $2, $3)")).
		WithArgs(bookID, bookRecommendationId, float32(0.85)).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = db.AddRecommendedBook(ctx, bookRecommendation)
	assert.NoError(t, err)
	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_AddRecommendedBook_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	db := &postgresDB{pool: mockPool}
	ctx := context.Background()
	bookID := 1
	bookRecommendationId := 2
	bookRecommendation := NewBookRecommendation{
		BookID:            bookID,
		RecommendedBookID: bookRecommendationId,
		Score:             0.85,
	}

	mockPool.ExpectExec(EscapeQuery("INSERT INTO book_recommendation (book_id, recommended_book_id, score) VALUES ($1, $2, $3)")).
		WithArgs(bookID, bookRecommendationId, float32(0.85)).
		WillReturnError(fmt.Errorf("failed to query recommended_books"))

	err = db.AddRecommendedBook(ctx, bookRecommendation)
	assert.Error(t, err)
}
