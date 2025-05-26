package services

import (
	"context"

	"app/server/domain"

	"github.com/stretchr/testify/mock"
)

type AuthorsServiceMock struct {
	mock.Mock
}

func (m *AuthorsServiceMock) GetAuthor(ctx context.Context, id int) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (m *AuthorsServiceMock) UpdateAuthor(ctx context.Context, author domain.Author) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (m *AuthorsServiceMock) DeleteAuthor(ctx context.Context, id int) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (m *AuthorsServiceMock) CreateAuthor(ctx context.Context, author domain.Author) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}
