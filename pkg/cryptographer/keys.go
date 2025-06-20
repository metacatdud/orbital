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

func NewPublicKeyFromString(publicKeyHex string) (*PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, err
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("%w:[%d]", ErrPublicKeySize, len(publicKeyBytes))
	}

	pubKey := ed25519.PublicKey(publicKeyBytes)

	return &PublicKey{key: pubKey}, nil
}

type PrivateKey struct {
	seed  []byte
	nonce [32]byte
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

	// Create nonce for PrivateKey
	var nonce [32]byte
	_, err = rand.Read(nonce[:])
	if err != nil {
		return nil, nil, err
	}

	publicKey := privateKey.PublicKey()

	return publicKey, &PrivateKey{seed: privateKey.Bytes(), nonce: nonce}, nil
}

// NewPrivateKeyFromSeed creat a private key from a seed
func NewPrivateKeyFromSeed(seed []byte) (*PrivateKey, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("%w:[%d]", ErrSeedSize, len(seed))
	}

	// Create nonce for PrivateKey
	var nonce [32]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, err
	}

	return &PrivateKey{seed: seed, nonce: nonce}, nil
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
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}

	return NewPrivateKeyFromSeed(seed)
}
