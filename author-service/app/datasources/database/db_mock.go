package database

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) AddAuthor(ctx context.Context, author NewAuthor) error {
	return m.Called(ctx, author).Error(0)
}

func (m *DatabaseMock) UpdateAuthor(ctx context.Context, author Author) error {
	return m.Called(ctx, author).Error(0)
}

func (m *DatabaseMock) DeleteAuthor(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *DatabaseMock) GetAuthor(ctx context.Context, id int) (Author, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Author), args.Error(1)
}

func (m *DatabaseMock) ListAuthor(ctx context.Context, filter Author, limit, offset int) ([]Author, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]Author), args.Error(1)
}

func (m *DatabaseMock) CloseConnections() {
	m.Called()
}
