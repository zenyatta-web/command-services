package repository

import "context"

type Repository interface {
	Close()
	Insert(ctx context.Context) error
	Update(ctx context.Context) error
}

var repository Repository

func SetRepository(r Repository) {
	repository = r
}

func Close() {
	repository.Close()
}

func Insert(ctx context.Context) error {
	return repository.Insert(ctx)
}

func Update(ctx context.Context) error {
	return repository.Update(ctx)
}
