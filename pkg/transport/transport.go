package transport

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Code uint32

const (
	OK Code = iota
	Canceled
	Unknown
	NotFound
	Unimplemented
	Unauthenticated
	Internal
	Unavailable
)

type ErrorResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// Decode reads and unmarshals a JSON payload into the provided destination.
func Decode(r io.ReadCloser, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r, 1024*1024))
	if err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	if err = json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	return nil
}

// Encode writes a JSON response with optional gzip when the client accepts it.
func Encode(w http.ResponseWriter, r *http.Request, status int, data any) error {
	var writer io.Writer = w
	useGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if useGzip {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		writer = gzw
		defer func() { _ = gzw.Close() }()
	}

	w.WriteHeader(status)

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
