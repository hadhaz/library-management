package database

import (
	"context"
	"fmt"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPostgresDB_AddAuthor_Success(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	db := &postgresDB{pool: mock}
	timeNow := time.Now()
	author := NewAuthor{"Jane", "Doe", &timeNow, "USA"}

	query := `INSERT INTO authors (first_name, last_name, birth_date, nationality)`
	mock.ExpectExec(EscapeQuery(query)).
		WithArgs(author.FirstName, author.LastName, author.BirthDate, author.Nationality).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := db.AddAuthor(context.Background(), author)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_AddAuthor_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer mockPool.Close()

	db := &postgresDB{pool: mockPool}
	timeNow := time.Now()
	author := NewAuthor{"Jane", "Doe", &timeNow, "USA"}

	query := `INSERT INTO authors (first_name, last_name, birth_date, nationality)`
	mockPool.ExpectExec(EscapeQuery(query)).
		WithArgs(author.FirstName, author.LastName, author.BirthDate, author.Nationality).
		WillReturnError(fmt.Errorf("insert failed"))

	err = db.AddAuthor(context.Background(), author)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to add author")

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPostgresDB_GetAuthor_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	timeNow := time.Now()
	authorID := 1

	expected := Author{
		ID:          authorID,
		FirstName:   "Jane",
		LastName:    "Doe",
		BirthDate:   &timeNow,
		Nationality: "USA",
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	query := `
		SELECT author_id, first_name, last_name, birth_date,
		       nationality, created_at, updated_at 
		FROM authors WHERE author_id = $1`

	mock.ExpectQuery(EscapeQuery(query)).
		WithArgs(authorID).
		WillReturnRows(pgxmock.NewRows([]string{
			"author_id", "first_name", "last_name", "birth_date",
			"nationality", "created_at", "updated_at",
		}).AddRow(expected.ID, expected.FirstName, expected.LastName, expected.BirthDate, expected.Nationality, expected.CreatedAt, expected.UpdatedAt))

	result, err := db.GetAuthor(context.Background(), authorID)
	assert.NoError(t, err)
	assert.Equal(t, expected.FirstName, result.FirstName)
	assert.Equal(t, expected.LastName, result.LastName)
	assert.Equal(t, expected.Nationality, result.Nationality)
	assert.Equal(t, expected.BirthDate, result.BirthDate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_GetAuthor_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	authorID := 999

	query := `
		SELECT author_id, first_name, last_name, birth_date,
		       nationality, created_at, updated_at 
		FROM authors WHERE author_id = $1`

	mock.ExpectQuery(EscapeQuery(query)).
		WithArgs(authorID).
		WillReturnError(fmt.Errorf("author not found"))

	result, err := db.GetAuthor(context.Background(), authorID)
	assert.Error(t, err)
	assert.Equal(t, "unable to get author: author not found", err.Error())
	assert.Equal(t, Author{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_ListAuthor_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	timeNow := time.Now()
	authorID := 1

	expected := Author{
		ID:          authorID,
		FirstName:   "Jane",
		LastName:    "Doe",
		BirthDate:   &timeNow,
		Nationality: "USA",
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	filter := Author{FirstName: "Jane"}
	limit := 10
	offset := 0

	query := `
        SELECT author_id, first_name, last_name, birth_date, nationality, created_at, updated_at
        FROM authors WHERE LOWER(first_name) LIKE $1 ORDER BY first_name LIMIT $2 OFFSET $3`

	mock.ExpectQuery(EscapeQuery(query)).
		WithArgs("%jane%", limit, offset).
		WillReturnRows(pgxmock.NewRows([]string{
			"author_id", "first_name", "last_name", "birth_date",
			"nationality", "created_at", "updated_at",
		}).AddRow(
			expected.ID,
			expected.FirstName,
			expected.LastName,
			expected.BirthDate,
			expected.Nationality,
			expected.CreatedAt,
			expected.UpdatedAt,
		))

	result, err := db.ListAuthor(context.Background(), filter, limit, offset)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, expected, result[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_ListAuthor_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	filter := Author{FirstName: "Jane"}
	limit := 10
	offset := 0

	query := `
        SELECT author_id, first_name, last_name, birth_date, nationality, created_at, updated_at
        FROM authors WHERE LOWER(first_name) LIKE $1 ORDER BY first_name LIMIT $2 OFFSET $3`

	mock.ExpectQuery(EscapeQuery(query)).
		WithArgs("%jane%", limit, offset).
		WillReturnError(fmt.Errorf("query failed"))

	authors, err := db.ListAuthor(context.Background(), filter, limit, offset)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to list authors")
	assert.Nil(t, authors)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_DeleteAuthor_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	authorID := 1

	query := `DELETE FROM authors WHERE id = $1`
	mock.ExpectExec(EscapeQuery(query)).
		WithArgs(authorID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = db.DeleteAuthor(context.Background(), authorID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_DeleteAuthor_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	authorID := 99

	query := `DELETE FROM authors WHERE id = $1`
	mock.ExpectExec(EscapeQuery(query)).
		WithArgs(authorID).
		WillReturnError(fmt.Errorf("some db error"))

	err = db.DeleteAuthor(context.Background(), authorID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to delete author")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_UpdateAuthor_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	timeNow := time.Now()

	author := Author{
		ID:          1,
		FirstName:   "Jane",
		LastName:    "Doe",
		BirthDate:   &timeNow,
		Nationality: "USA",
	}

	query := `
		UPDATE authors 
		 SET first_name = $1,
		     last_name = $2,
		     birth_date = $3,
		     nationality = $4
		 WHERE id = $5`

	mock.ExpectExec(EscapeQuery(query)).
		WithArgs(author.FirstName, author.LastName, author.BirthDate, author.Nationality, author.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1)) // simulate 1 row updated

	err = db.UpdateAuthor(context.Background(), author)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresDB_UpdateAuthor_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &postgresDB{pool: mock}
	timeNow := time.Now()

	author := Author{
		ID:          1,
		FirstName:   "Jane",
		LastName:    "Doe",
		BirthDate:   &timeNow,
		Nationality: "USA",
	}

	query := `
		UPDATE authors 
		 SET first_name = $1,
		     last_name = $2,
		     birth_date = $3,
		     nationality = $4
		 WHERE id = $5`

	mock.ExpectExec(EscapeQuery(query)).
		WithArgs(author.FirstName, author.LastName, author.BirthDate, author.Nationality, author.ID).
		WillReturnError(fmt.Errorf("update failed"))

	err = db.UpdateAuthor(context.Background(), author)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to update author")
	assert.NoError(t, mock.ExpectationsWereMet())
}
