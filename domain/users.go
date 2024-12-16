package domain

import (
	"fmt"
	database "orbital/pkg/db"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	PubKey string `json:"pubKey"`
	Access string `json:"access"`
}

type Users []User

type UserRepository struct {
	db *database.DB
}

func (repo UserRepository) Save(u User) error {
	query := `INSERT INTO users (id, name, pubkey, access) VALUES (?, ?, ?, ?)`
	_, err := repo.db.Client().Exec(query, u.ID, u.Name, u.PubKey, u.Access)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (repo UserRepository) FindByID(id string) (*User, error) {
	query := `SELECT id, name, pubkey, access FROM users WHERE id = ?`
	row := repo.db.Client().QueryRow(query, id)
	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.PubKey); err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (repo UserRepository) FindByPublicKey(pubKey string) (*User, error) {
	query := `SELECT id, name, pubkey, access FROM users WHERE pubkey = ?`
	row := repo.db.Client().QueryRow(query, pubKey)
	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.PubKey, &user.Access); err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (repo UserRepository) ExistsByPublicKey(pubKey string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT pubkey FROM users WHERE pubkey = ?);`
	err := repo.db.Client().QueryRow(query, pubKey).Scan(&exists)

	return exists, err
}

func NewUserRepository(db *database.DB) UserRepository {
	return UserRepository{db: db}
}
