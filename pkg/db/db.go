package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
	"path/filepath"
)

type DB struct {
	client *sql.DB
}

func (db *DB) Client() *sql.DB {
	return db.client
}

func NewDB(dbDirPath string) (*DB, error) {
	dbpath := filepath.Join(dbDirPath, "orbital.db")
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		return nil, fmt.Errorf("%w:[%s]", ErrDBOpen, err.Error())
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%w:[%s]", ErrDBConnect, err.Error())
	}

	return &DB{
		client: db,
	}, nil
}

func AutoMigrate(db *DB, orbitalDir string) error {
	migrationsPath := fmt.Sprintf("file://%s/data/migrations", orbitalDir)
	driver, err := sqlite.WithInstance(db.Client(), &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite", driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
