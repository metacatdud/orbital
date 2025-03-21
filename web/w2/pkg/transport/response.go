package transport

type (
	ErrorResponse struct {
		Type string `json:"type"`
		Msg  string `json:"msg"`
	}
)
