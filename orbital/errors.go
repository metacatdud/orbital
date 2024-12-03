package orbital

import "errors"

var (
	ErrBadPayload       = errors.New("unable to decode payload")
	ErrUnmarshalPayload = errors.New("unable to unmarshal payload")
	ErrPathNotFound     = errors.New("path not found")
)

type Error struct {
	Code Code        `json:"code"`
	Msg  interface{} `json:"msg"`
}

func (e *Error) Error() string {
	switch msg := e.Msg.(type) {
	case string:
		return msg
	case error:
		return msg.Error()
	default:
		return "Unknown error"
	}
}
