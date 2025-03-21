package orbital

import "errors"

var (
	ErrRegNotFound    = errors.New(`module not found`)
	ErrRegDuplicateID = errors.New(`module already exists`)
	ErrRegWrongType   = errors.New(`module type is wrong`)
)
