package proto

import (
	"encoding/json"
	"errors"
	"orbital/pkg/cryptographer"
)

type Message = cryptographer.Message

type Topic struct {
	Name string `json:"name"`
}

var TimestampNow = cryptographer.Now

func Encode(sk cryptographer.PrivateKey, metadata, body any) (*Message, error) {
	pubK := sk.PublicKey()
	var (
		m []byte
		b []byte
	)

	if metadata != nil {
		m, _ = json.Marshal(metadata)
	}

	if body != nil {
		b, _ = json.Marshal(body)
	}

	msg := &Message{
		PublicKey: pubK.Compress(),
		V:         1,
		Timestamp: TimestampNow(),
		Metadata:  m,
		Body:      b,
	}

	if err := msg.Sign(sk.Bytes()); err != nil {
		return nil, errors.New("sign msg fail")
	}

	return msg, nil
}

// Decode validate message signature and decide the body
func Decode(msg Message, body, metadata any) error {
	valid, err := msg.Verify()
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid message signature")
	}

	if metadata != nil {
		if err = json.Unmarshal(msg.Body, &metadata); err != nil {
			return err
		}
	}

	if body != nil {
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	}

	return nil
}
