package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
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

func (db *postgresDB) AddAuthor(ctx context.Context, author NewAuthor) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO authors (first_name, last_name, birth_date, nationality)`,
		author.FirstName, author.LastName, author.BirthDate, author.Nationality)
	if err != nil {
		return fmt.Errorf("unable to add author: %v", err)
	}

	return nil
}

func (db *postgresDB) UpdateAuthor(ctx context.Context, author Author) error {
	_, err := db.pool.Exec(ctx,
		`UPDATE authors 
		 SET first_name = $1,
		     last_name = $2,
		     birth_date = $3,
		     nationality = $4
		 WHERE id = $5`,
		author.FirstName, author.LastName, author.BirthDate, author.Nationality, author.ID)
	if err != nil {
		return fmt.Errorf("unable to update author: %v", err)
	}

	return nil
}

func (db *postgresDB) DeleteAuthor(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx,
		`DELETE FROM authors WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("unable to delete author: %v", err)
	}

	return nil
}

func (db *postgresDB) GetAuthor(ctx context.Context, id int) (Author, error) {
	query := `
		SELECT author_id, first_name, last_name, birth_date,
		       nationality, created_at, updated_at 
		FROM authors WHERE author_id = $1`

	var author Author
	err := db.pool.QueryRow(ctx, query, id).Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.BirthDate,
		&author.Nationality,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		return Author{}, fmt.Errorf("unable to get author: %v", err)
	}

	return author, nil
}

func (db *postgresDB) ListAuthor(ctx context.Context, filter Author, limit, offset int) ([]Author, error) {
	var (
		args  []interface{}
		where []string
	)

	if filter.FirstName != "" {
		args = append(args, "%"+strings.ToLower(filter.FirstName)+"%")
		where = append(where, fmt.Sprintf("LOWER(first_name) LIKE $%d", len(args)))
	}
	if filter.LastName != "" {
		args = append(args, "%"+strings.ToLower(filter.LastName)+"%")
		where = append(where, fmt.Sprintf("LOWER(last_name) LIKE $%d", len(args)))
	}

	query := `
        SELECT author_id, first_name, last_name, birth_date, nationality, created_at, updated_at
        FROM authors
    `
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " OR ")
	}

	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY first_name LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to list authors: %w", err)
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var a Author
		if err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.BirthDate,
			&a.Nationality,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning author row: %w", err)
		}
		authors = append(authors, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading author rows: %w", err)
	}

	return authors, nil
}

func (db *postgresDB) CloseConnections() {
	db.pool.Close()
}
