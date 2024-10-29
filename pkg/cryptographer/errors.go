package cryptographer

import "errors"

var (
	ErrTLVLenExceed       = errors.New("tlv length size exceed")
	ErrMetadataSizeExceed = errors.New("metadata exceed allowed limit")
	ErrMessageSizeExceed  = errors.New("message exceed allowed limit")
	ErrSeedSize           = errors.New("invalid seed size")
	ErrPublickeySize      = errors.New("invalid public key size")
)
