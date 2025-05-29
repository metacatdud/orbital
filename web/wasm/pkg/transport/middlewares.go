package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"orbital/pkg/proto"
)

func VerifyAndUnwrap(raw []byte) ([]byte, error) {
	var msg proto.Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, fmt.Errorf("invalid envelope JSON: %w", err)
	}

	ok, err := msg.Verify()
	if err != nil {
		return nil, fmt.Errorf("message verification error: %w", err)
	}
	if !ok {
		return nil, errors.New("message signature invalid")
	}

	return msg.Body, nil
}
