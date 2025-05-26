package services

import (
	"app/datasources/database"
	"app/server/domain"
	"context"
)

type AuthorsService interface {
	GetAuthors(ctx context.Context) ([]domain.Author, error)
	GetAuthor(ctx context.Context, id int) (domain.Author, error)
	UpdateAuthor(ctx context.Context, author domain.Author) (domain.Author, error)
	DeleteAuthor(ctx context.Context, id int) (domain.Author, error)
	CreateAuthor(ctx context.Context, author domain.Author) (domain.Author, error)
}

type authorsService struct {
	db database.Database
}

func (a authorsService) GetAuthors(ctx context.Context) ([]domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (a authorsService) GetAuthor(ctx context.Context, id int) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (a authorsService) UpdateAuthor(ctx context.Context, author domain.Author) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (a authorsService) DeleteAuthor(ctx context.Context, id int) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (a authorsService) CreateAuthor(ctx context.Context, author domain.Author) (domain.Author, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthorsService(db database.Database) AuthorsService {
	return &authorsService{db: db}
}
