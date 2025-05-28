package proto

import (
	"encoding/json"
	"errors"
	"orbital/pkg/cryptographer"
)

type Message = cryptographer.Message

var TimestampNow = cryptographer.Now

func Encode(sk *cryptographer.PrivateKey, metadata *cryptographer.Metadata, body any) (*Message, error) {
	pubK := sk.PublicKey()
	var (
		b []byte
	)

	if body != nil {
		b, _ = json.Marshal(body)
	}

	if metadata == nil {
		metadata = &cryptographer.Metadata{}
	}

	msg := &Message{
		PublicKey: pubK.Compress(),
		V:         1,
		Timestamp: TimestampNow(),
		Metadata:  metadata,
		Body:      b,
	}

	if err := msg.Sign(sk.Bytes()); err != nil {
		return nil, errors.New("sign msg fail")
	}

	return msg, nil
}

// Decode validate message signature and decide the body
func Decode(msg Message, body interface{}) error {
	valid, err := msg.Verify()
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid message signature")
	}

	if body != nil {
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	}

	return nil
}
