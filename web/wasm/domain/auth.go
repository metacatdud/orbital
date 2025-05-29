package domain

import "orbital/web/wasm/pkg/storage"

const (
	AuthStorageKey RepositoryKey = "auth"
)

type Auth struct {
	SecretKey string `json:"secretKey"`
}

type AuthRepository struct {
	base Repository[Auth]
}

func NewAuthRepository(db storage.Storage) *AuthRepository {
	return &AuthRepository{
		base: NewRepository[Auth](db, AuthStorageKey),
	}
}

func (u *AuthRepository) Save(auth Auth) error {
	return u.base.Save(auth)
}

func (u *AuthRepository) Get() (*Auth, error) {
	return u.base.Get()
}

func (u *AuthRepository) Delete() error {
	return u.base.Remove()
}
