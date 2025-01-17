package domain

import "orbital/web/wasm/storage"

type Auth struct {
	SecretKey string `json:"secretKey"`
}

type AuthRepository struct {
	db storage.Storage
}

func (repo AuthRepository) Save(auth Auth) error {
	return repo.db.Set("auth", auth)
}

func (repo AuthRepository) Get() (*Auth, error) {
	auth := &Auth{}
	if err := repo.db.Get("auth", auth); err != nil {
		return nil, err
	}

	return auth, nil
}

func (repo AuthRepository) Remove() error {
	return repo.db.Del("auth")
}

func NewAuthRepository(db storage.Storage) AuthRepository {
	return AuthRepository{db: db}
}
