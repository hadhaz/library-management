package database

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) BorrowBook(ctx context.Context, book NewBorrowingRecord) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *DatabaseMock) ReturnBook(ctx context.Context, book BorrowingRecord) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *DatabaseMock) GetRecommendedBooks(ctx context.Context, bookID int) ([]BookRecommendation, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).([]BookRecommendation), args.Error(1)
}

func (m *DatabaseMock) AddRecommendedBook(ctx context.Context, book NewBookRecommendation) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *DatabaseMock) LoadAllBooks(ctx context.Context) ([]Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Book), args.Error(1)
}

func (m *DatabaseMock) CreateBook(ctx context.Context, newBook NewBook) error {
	args := m.Called(ctx, newBook)
	return args.Error(0)
}

func (m *DatabaseMock) DeleteBook(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DatabaseMock) UpdateBook(ctx context.Context, book Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *DatabaseMock) CloseConnections() {
}
