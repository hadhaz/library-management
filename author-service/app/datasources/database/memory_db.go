package database

import "context"

func newMemoryDB() Database {
	return &memoryDB{
		records:   make([]Author, 0, 10),
		idCounter: 0,
	}
}

type memoryDB struct {
	records   []Author
	idCounter int
}

func (db *memoryDB) AddAuthor(ctx context.Context, author NewAuthor) error {
	//TODO implement me
	panic("implement me")
}

func (db *memoryDB) UpdateAuthor(ctx context.Context, author Author) error {
	//TODO implement me
	panic("implement me")
}

func (db *memoryDB) DeleteAuthor(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (db *memoryDB) GetAuthor(ctx context.Context, id int) (Author, error) {
	//TODO implement me
	panic("implement me")
}

func (db *memoryDB) ListAuthor(ctx context.Context, filter Author, limit, offset int) ([]Author, error) {
	//TODO implement me
	panic("implement me")
}

func (db *memoryDB) CloseConnections() {
}
