package auth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
)

func MessageDecode() orbital.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var msg cryptographer.Message
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				http.Error(w, "bad JSON envelope", 400)
				return
			}
			valid, err := msg.Verify()
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			if !valid {
				http.Error(w, "bad JSON envelope", 400)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, cryptographer.BodyCtxKey, msg.Body)
			ctx = context.WithValue(ctx, cryptographer.PublicKeyCtxKey, hex.EncodeToString(msg.PublicKey[:]))

			next(w, r.WithContext(ctx))
		}
	}
}

func ValidateRole() orbital.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("ValidateRole unhandled. Coming Soon!")
			next(w, r)
		}
	}
}
