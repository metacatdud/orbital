package domain

import "orbital/web/wasm/pkg/storage"

const (
	UserStorageKey RepositoryKey = "user"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access string `json:"access"`
}

type UserRepository struct {
	base Repository[User]
}

func NewUserRepository(db storage.Storage) *UserRepository {
	return &UserRepository{
		base: NewRepository[User](db, UserStorageKey),
	}
}

func (u *UserRepository) Save(user User) error {
	return u.base.Save(user)
}

func (u *UserRepository) Get() (*User, error) {
	return u.base.Get()
}

func (u *UserRepository) Delete() error {
	return u.base.Remove()
}
