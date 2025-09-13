package cryptographer

import "errors"

var (
	ErrTLVLenExceed       = errors.New("tlv length size exceed")
	ErrMetadataSizeExceed = errors.New("metadata exceed allowed limit")
	ErrMetadataTagMarshal = errors.New("metadata tags marshal error")
	ErrMessageSizeExceed  = errors.New("message exceed allowed limit")
	ErrSeedSize           = errors.New("invalid seed size")
	ErrInvalidKeySize     = errors.New("invalid key size")
	ErrSignMessage        = errors.New("sign message failed")
	ErrPubKeyMessage      = errors.New("cannot create public key bytes")
)
