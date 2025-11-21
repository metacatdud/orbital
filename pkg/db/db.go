package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"atomika.io/atomika/atomika"
)

type DB struct {
	client  *sql.DB
	manager *atomika.DBManager
}

func (db *DB) Client() *sql.DB {
	return db.client
}

func NewDB(dbDirPath string) (*DB, error) {
	dbpath := filepath.Join(dbDirPath, "orbital.db")

	cfg := &atomika.CfgDatabase{
		Driver:         atomika.DBDriverSQLite,
		DSN:            dbpath,
		AutoMigrate:    true,
		MigrationsPath: filepath.Join(dbDirPath, "migrations"),
	}

	manager, err := atomika.NewDBManager(cfg)
	if err != nil {
		return nil, err
	}

	dbClient, err := manager.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w:[%s]", ErrDBConnect, err.Error())
	}

	return &DB{
		client:  dbClient,
		manager: manager,
	}, nil
}

func AutoMigrate(ctx context.Context, db *DB) error {
	if db == nil || db.manager == nil {
		return fmt.Errorf("db manager not initialized")
	}
	return db.manager.AutoMigrate(ctx)
}
