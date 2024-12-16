package db

import "errors"

var (
	ErrDBOpen    = errors.New("failed to open database")
	ErrDBConnect = errors.New("cannot connect to database")
)
