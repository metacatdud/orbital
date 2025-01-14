package stringer

import "errors"

var (
	ErrRandTooShort = errors.New("string to short to randomize")
)
