package cryptographer

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type PublicKey struct {
	key ed25519.PublicKey
}

func (pk PublicKey) Bytes() []byte {
	if pk.key == nil {
		return nil
	}
	cp := make([]byte, ed25519.PublicKeySize)
	copy(cp, pk.key)
	return cp
}

func (pk PublicKey) ToHex() string {
	return hex.EncodeToString(pk.key)
}

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (sk PrivateKey) Bytes() []byte {
	if sk.key == nil {
		return nil
	}
	cp := make([]byte, ed25519.PrivateKeySize)
	copy(cp, sk.key)
	return cp
}

func (sk PrivateKey) Seed() []byte {
	if sk.key == nil || len(sk.key) != ed25519.PrivateKeySize {
		return nil
	}
	cp := make([]byte, ed25519.SeedSize)
	copy(cp, sk.key.Seed())
	return cp
}

func (sk PrivateKey) PublicKey() PublicKey {
	if sk.key == nil {
		return PublicKey{}
	}

	pk := sk.key.Public().(ed25519.PublicKey)

	return PublicKey{key: pk}
}

// GenerateKeysPair generates a new public and private key pair
func GenerateKeysPair() (PublicKey, PrivateKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return PublicKey{}, PrivateKey{}, err
	}

	return PublicKey{key: pk}, PrivateKey{key: sk}, nil
}

// NewPrivateKeyFromSeed creat a private key from a seed
func NewPrivateKeyFromSeed(seed []byte) (PrivateKey, error) {
	if len(seed) != ed25519.SeedSize {
		return PrivateKey{}, fmt.Errorf("%w:[%d]", ErrSeedSize, len(seed))
	}

	ed25519Sk := ed25519.NewKeyFromSeed(seed)

	return PrivateKey{key: ed25519Sk}, nil
}

// NewPrivateKeyFromHex creat a private key from a string
func NewPrivateKeyFromHex(skStr string) (PrivateKey, error) {
	skBytes, err := hex.DecodeString(skStr)
	if err != nil {
		return PrivateKey{}, fmt.Errorf("%w:[secret: %s]", ErrInvalidKeySize, err.Error())
	}

	switch len(skBytes) {
	case ed25519.PrivateKeySize:
		keyBytes := make([]byte, ed25519.PrivateKeySize)
		copy(keyBytes, skBytes)
		return PrivateKey{key: keyBytes}, nil
	case ed25519.SeedSize:
		return NewPrivateKeyFromSeed(skBytes)
	default:
		return PrivateKey{}, ErrInvalidKeySize
	}
}
