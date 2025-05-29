package domain

import (
	"errors"
	"fmt"
	"orbital/web/wasm/pkg/storage"
)

type RepositoryKey string

func (r RepositoryKey) String() string {
	return string(r)
}

type Repository[T any] struct {
	db  storage.Storage
	key RepositoryKey
}

func NewRepository[T any](db storage.Storage, key RepositoryKey) Repository[T] {
	return Repository[T]{
		db:  db,
		key: key,
	}
}

func (repo Repository[T]) Save(value T) error {
	return repo.db.Set(repo.key.String(), value)
}

func (repo Repository[T]) Get() (*T, error) {
	var result T
	if err := repo.db.Get(repo.key.String(), &result); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%w:[%s]", ErrKeyNotFound, repo.key)
		}
		return nil, err
	}

	return &result, nil
}

func (repo Repository[T]) Remove() error {
	return repo.db.Del(repo.key.String())
}
