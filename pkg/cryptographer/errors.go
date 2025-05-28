package cryptographer

import "errors"

var (
	ErrTLVLenExceed       = errors.New("tlv length size exceed")
	ErrMetadataSizeExceed = errors.New("metadata exceed allowed limit")
	ErrMetadataTagMarshal = errors.New("metadata tags marshal error")
	ErrMessageSizeExceed  = errors.New("message exceed allowed limit")
	ErrSeedSize           = errors.New("invalid seed size")
	ErrPublicKeySize      = errors.New("invalid public key size")
	ErrInvalidKeySize     = errors.New("invalid key size")
)
