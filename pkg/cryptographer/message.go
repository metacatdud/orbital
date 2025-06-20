package cryptographer

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// Payload cannot exceed 50MB
	maxSize = 50267340
)

type CtxKey string

const (
	BodyCtxKey      CtxKey = "body"
	PublicKeyCtxKey CtxKey = "publicKey"
)

const (
	TypeID        = 1
	TypePublicKey = 2
	TypeSignature = 3
	TypeV         = 4
	TypeTimestamp = 5
	TypeBody      = 6
	TypeMetadata  = 7
)

type Message struct {
	ID        [32]byte  `json:"id"`
	PublicKey [32]byte  `json:"publicKey"`
	V         int64     `json:"v"`
	Timestamp Timestamp `json:"timestamp"`
	Metadata  *Metadata `json:"metadata"`
	Body      []byte    `json:"body"`
	Signature [64]byte  `json:"sig"`
}

func (m *Message) ComputeID() ([32]byte, error) {
	serializedData, err := m.Serialize()
	if err != nil {
		return [32]byte{}, err
	}

	// Compute SHA-256 hash directly on the serialized data
	hash := sha256.Sum256(serializedData)

	return hash, nil
}

func (m *Message) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Serialize PublicKey
	if err := writeTLVtoBuffer(&buf, TypePublicKey, m.PublicKey[:]); err != nil {
		return nil, err
	}

	// Serialize Version
	verByt := intToByte(m.V)
	if verByt == nil {
		return nil, errors.New("cannot convert [V] to byte")
	}

	if err := writeTLVtoBuffer(&buf, TypeV, verByt); err != nil {
		return nil, err
	}

	if err := writeTLVtoBuffer(&buf, TypeTimestamp, m.Timestamp.Bytes()); err != nil {
		return nil, err
	}

	// Serialize Metadata
	metaBytes, err := m.Metadata.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize metadata: %w", err)
	}

	if err = writeTLVtoBuffer(&buf, TypeMetadata, metaBytes); err != nil {
		return nil, err
	}

	// Serialize Body
	if err = writeTLVtoBuffer(&buf, TypeBody, m.Body); err != nil {
		return nil, err
	}

	// Check total size limit
	if buf.Len() > maxSize {
		return nil, fmt.Errorf("%w:[%d]", ErrMessageSizeExceed, buf.Len())
	}

	return buf.Bytes(), nil
}

func (m *Message) Sign(privateKey32Byte []byte) error {
	if len(privateKey32Byte) != ed25519.SeedSize {
		return fmt.Errorf("%w:[%d]", ErrSeedSize, ed25519.SeedSize)
	}
	privateKey := ed25519.NewKeyFromSeed(privateKey32Byte)

	serial, err := m.Serialize()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(serial)
	sig := ed25519.Sign(privateKey, hash[:])

	m.ID, err = m.ComputeID()
	if err != nil {
		return err
	}

	m.Signature = [64]byte(sig)

	return nil
}

func (m *Message) Verify() (bool, error) {

	serial, err := m.Serialize()
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(serial)

	return ed25519.Verify(m.PublicKey[:], hash[:], m.Signature[:]), nil
}

func Encode(sk *PrivateKey, metadata *Metadata, body any) (*Message, error) {
	pubK := sk.PublicKey()
	var (
		b []byte
	)

	if body != nil {
		b, _ = json.Marshal(body)
	}

	if metadata == nil {
		metadata = &Metadata{}
	}

	msg := &Message{
		PublicKey: pubK.Compress(),
		V:         1,
		Timestamp: Now(),
		Metadata:  metadata,
		Body:      b,
	}

	if err := msg.Sign(sk.Bytes()); err != nil {
		return nil, errors.New("sign msg fail")
	}

	return msg, nil
}
