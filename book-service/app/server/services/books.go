package services

import (
	"context"
	"fmt"

	"app/datasources/database"
	"app/server/domain"
)

type BooksService interface {
	GetBooks(ctx context.Context) ([]domain.Book, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	SaveBook(ctx context.Context, newBook domain.Book) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, book domain.Book) error
}

type booksService struct {
	db database.Database
}

func (s *booksService) GetBook(ctx context.Context, id int) (domain.Book, error) {
	record, err := s.db.GetBookByID(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}
	book := domain.Book{
		ID:          record.ID,
		Title:       record.Title,
		AuthorID:    record.AuthorID,
		CategoryID:  record.CategoryID,
		PublishDate: record.PublishedDate,
		Description: record.Description,
		ISBN:        record.ISBN,
	}

	return book, nil
}

func NewBooksService(db database.Database) BooksService {
	return &booksService{db: db}
}

func (s *booksService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	dbRecords, err := s.db.LoadAllBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load books: %w", err)
	}

	books := make([]domain.Book, 0, len(dbRecords))
	for _, record := range dbRecords {
		books = append(books, domain.Book{
			ID:          record.ID,
			Title:       record.Title,
			AuthorID:    record.AuthorID,
			ISBN:        record.ISBN,
			PublishDate: record.PublishedDate,
			Description: record.Description,
			CategoryID:  record.CategoryID,
		})
	}

	return books, nil
}

func (s *booksService) SaveBook(ctx context.Context, book domain.Book) error {
	dbBook := database.NewBook{
		Title:         book.Title,
		AuthorID:      book.AuthorID,
		ISBN:          book.ISBN,
		Description:   book.Description,
		CategoryID:    book.CategoryID,
		PublishedDate: book.PublishDate,
		Stock:         12,
	}

	err := s.db.CreateBook(ctx, dbBook)
	if err != nil {
		return fmt.Errorf("failed to save book: %w", err)
	}

	return nil
}

func (s *booksService) DeleteBook(ctx context.Context, id int) error {
	err := s.db.DeleteBook(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

func (s *booksService) UpdateBook(ctx context.Context, book domain.Book) error {
	dbBook := database.Book{
		Title:       book.Title,
		AuthorID:    book.AuthorID,
		ISBN:        book.ISBN,
		Description: book.Description,
		ID:          book.ID,
	}

	err := s.db.UpdateBook(ctx, dbBook)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}
