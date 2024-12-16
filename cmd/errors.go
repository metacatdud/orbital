package cmd

import "errors"

// Errors
var (
	ErrInvalidIP          = errors.New("invalid ip address")
	ErrInvalidEd25519Key  = errors.New("invalid ed25519 key")
	ErrInvalidEd25519Seed = errors.New("invalid ed25519 seed")
	ErrCannotCreateDir    = errors.New("cannot create dir")
	ErrInvalidFilepath    = errors.New("invalid filepath")
	ErrReadFile           = errors.New("error reading file")
	ErrCreateFile         = errors.New("error creating file")
	ErrWriteFile          = errors.New("error writing file")
)
