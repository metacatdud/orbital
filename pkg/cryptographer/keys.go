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

// Compress compresses the public key to a [32]byte format
func (pk *PublicKey) Compress() [32]byte {
	var compressed [32]byte
	copy(compressed[:], pk.key)
	return compressed
}

func (pk *PublicKey) String() string {
	return hex.EncodeToString(pk.key[:])
}

type PrivateKey struct {
	seed  []byte
	nonce [32]byte // Not in use for now
}

// PublicKey returns the public key corresponding to the secret key
func (privateKey *PrivateKey) PublicKey() *PublicKey {
	privateKeyGen := ed25519.NewKeyFromSeed(privateKey.seed)
	pubKey := privateKeyGen.Public().(ed25519.PublicKey)

	return &PublicKey{key: pubKey}
}

func (privateKey *PrivateKey) Bytes() []byte {
	return privateKey.seed
}

func (privateKey *PrivateKey) String() string {
	return hex.EncodeToString(privateKey.seed)
}

// GenerateKeysPair generates a new public and private key pair
func GenerateKeysPair() (*PublicKey, *PrivateKey, error) {
	privateKey, err := NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	return privateKey.PublicKey(), privateKey, nil
}

// NewPrivateKeyFromSeed creat a private key from a seed
func NewPrivateKeyFromSeed(seed []byte) (*PrivateKey, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("%w:[%d]", ErrSeedSize, len(seed))
	}

	ed25519Sk := ed25519.NewKeyFromSeed(seed)

	// Create nonce for PrivateKey
	var nonce [32]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	return &PrivateKey{seed: ed25519Sk.Seed(), nonce: nonce}, nil
}

// NewPrivateKeyFromString creat a private key from a string
func NewPrivateKeyFromString(privateKey string) (*PrivateKey, error) {
	seed, err := hex.DecodeString(privateKey)
	if err != nil || len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("%w:[secret: %s]", ErrInvalidKeySize, privateKey)
	}

	return NewPrivateKeyFromSeed(seed)
}

// NewPrivateKey generates a random seed key
func NewPrivateKey() (*PrivateKey, error) {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}

	return NewPrivateKeyFromSeed(seed)
}
