package database

import "context"

func newMemoryDB() Database {
	return &memoryDB{
		records:   make([]Book, 0, 10),
		idCounter: 0,
	}
}

type memoryDB struct {
	records   []Book
	idCounter int
}

func (db *memoryDB) GetBookByID(ctx context.Context, bookID int) (Book, error) {
	return db.records[db.idCounter], nil
}

func (db *memoryDB) BorrowBook(ctx context.Context, book NewBorrowingRecord) error {
	return nil
}

func (db *memoryDB) ReturnBook(ctx context.Context, books BorrowingRecord) error {
	return nil
}

func (db *memoryDB) AddRecommendedBook(ctx context.Context, book NewBookRecommendation) error {
	return nil
}

func (db *memoryDB) GetRecommendedBooks(ctx context.Context, bookID int) ([]BookRecommendation, error) {
	return []BookRecommendation{}, nil
}

func (db *memoryDB) LoadAllBooks(_ context.Context) ([]Book, error) {
	return db.records, nil
}

func (db *memoryDB) CreateBook(_ context.Context, newBook NewBook) error {
	db.records = append(db.records, Book{
		ID:    db.idCounter,
		Title: newBook.Title,
	})
	db.idCounter++
	return nil
}

func (db *memoryDB) UpdateBook(_ context.Context, book Book) error {
	return nil
}

func (db *memoryDB) DeleteBook(_ context.Context, id int) error {
	return nil
}

func (db *memoryDB) CloseConnections() {
}
