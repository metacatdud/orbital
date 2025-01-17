package domain

import (
	"orbital/web/wasm/dom"
	"orbital/web/wasm/storage"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Access    string `json:"access"`
}

type UserRepository struct {
	db storage.Storage
}

func (repo UserRepository) Save(u User) error {
	return repo.db.Set("user", u)
}

func (repo UserRepository) Get() (*User, error) {
	u := &User{}
	if err := repo.db.Get("user", u); err != nil {
		return nil, err
	}

	return u, nil
}

func (repo UserRepository) Remove() error {
	return repo.db.Del("user")
}

func (repo UserRepository) HasSession() bool {
	u := &User{}
	if err := repo.db.Get("user", u); err != nil {
		dom.PrintToConsole("user not found")
		return false
	}

	return u.PublicKey != ""
}

func NewUserRepository(db storage.Storage) UserRepository {
	return UserRepository{db: db}
}
