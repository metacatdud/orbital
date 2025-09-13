package domain

import (
	"database/sql"
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

type usersRow struct {
	ID     string
	Name   sql.NullString
	PubKey sql.NullString
	Access sql.NullString
}

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return UserRepository{db: db}
}

func (repo UserRepository) Save(u User) error {

	ur := mapUserToRow(u)

	args := []any{
		ur.ID,
		ur.Name,
		ur.PubKey,
		ur.Access,
	}

	query := `INSERT INTO users (id, name, pubkey, access) VALUES (?, ?, ?, ?)`
	_, err := repo.db.Client().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (repo UserRepository) GetByID(id string) (*User, error) {
	query := `SELECT id, name, pubkey, access FROM users WHERE id = ?`
	row := repo.db.Client().QueryRow(query, id)

	var userR usersRow
	if err := row.Scan(&userR.ID, &userR.Name, &userR.PubKey, userR.Access); err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	user := mapRowToUser(userR)
	return &user, nil
}

func (repo UserRepository) GetByPublicKey(pubKey string) (*User, error) {
	query := `SELECT id, name, pubkey, access FROM users WHERE pubkey = ?`
	row := repo.db.Client().QueryRow(query, pubKey)

	var userR usersRow
	if err := row.Scan(&userR.ID, &userR.Name, &userR.PubKey, &userR.Access); err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	user := mapRowToUser(userR)
	return &user, nil
}

func (repo UserRepository) ExistsByPublicKey(pubKey string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT pubkey FROM users WHERE pubkey = ? LIMIT 1);`
	err := repo.db.Client().QueryRow(query, pubKey).Scan(&exists)

	return exists, err
}

func (repo UserRepository) Find() (Users, error) {
	rows, err := repo.db.Client().Query(`SELECT id, name, pubkey, access FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users Users
	for rows.Next() {
		var userR usersRow
		if err = rows.Scan(&userR.ID, &userR.Name, &userR.PubKey); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		users = append(users, mapRowToUser(userR))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

func mapRowToUser(ur usersRow) User {
	return User{
		ID:     ur.ID,
		Name:   nullToString(ur.Name),
		PubKey: nullToString(ur.PubKey),
		Access: nullToString(ur.Access),
	}
}

func mapUserToRow(user User) usersRow {
	return usersRow{
		ID:     user.ID,
		Name:   stringToNull(user.Name),
		PubKey: stringToNull(user.PubKey),
		Access: stringToNull(user.Access),
	}
}
