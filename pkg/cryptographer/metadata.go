package cryptographer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

const (
	// Metadata cannot exceed 1MB
	maxMetadata = 0x100000
)

const (
	TypeMetaDomain        = 1
	TypeMetaAction        = 2
	TypeMetaNonce         = 3
	TypeMetaCorrelationID = 4
	TypeMetaTags          = 5
)

type Metadata struct {
	Domain        string            `json:"domain"` // e.g., "auth", "recovery", "trust", "message"
	Action        string            `json:"action"` // e.g., "register", "rotate", "send", "request"
	Nonce         string            `json:"nonce"`  // optional nonce for replay protection
	CorrelationID string            `json:"cid"`    // track message flow/threads
	Tags          map[string]string `json:"tags"`   // custom defined tags
}

func (m *Metadata) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	if m.Nonce == "" {
		m.Nonce = uuid.NewString()
	}

	if m.Tags == nil {
		m.Tags = make(map[string]string)
	}

	if err := writeTLVtoBuffer(&buf, TypeMetaDomain, []byte(m.Domain)); err != nil {
		return nil, err
	}

	if err := writeTLVtoBuffer(&buf, TypeMetaAction, []byte(m.Action)); err != nil {
		return nil, err
	}

	if m.Nonce != "" {
		if err := writeTLVtoBuffer(&buf, TypeMetaNonce, []byte(m.Nonce)); err != nil {
			return nil, err
		}
	}

	if m.CorrelationID != "" {
		if err := writeTLVtoBuffer(&buf, TypeMetaCorrelationID, []byte(m.CorrelationID)); err != nil {
			return nil, err
		}
	}

	tagBytes, err := json.Marshal(m.Tags)
	if err != nil {
		return nil, fmt.Errorf("%w:[%v]", ErrMetadataTagMarshal, err)
	}

	if err = writeTLVtoBuffer(&buf, TypeMetaTags, tagBytes); err != nil {
		return nil, err
	}

	if buf.Len() > maxMetadata {
		return nil, fmt.Errorf("%w:[%d bytes]", ErrMetadataSizeExceed, buf.Len())
	}

	return buf.Bytes(), nil
}
