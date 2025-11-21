package auth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"orbital/pkg/cryptographer"

	"atomika.io/atomika/atomika"
)

func MessageDecode() atomika.Interceptor {
	return interceptorFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var msg cryptographer.Message
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				http.Error(w, "bad JSON envelope", http.StatusBadRequest)
				return
			}
			valid, err := msg.Verify()
			if err != nil || !valid {
				http.Error(w, "bad JSON envelope", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, cryptographer.BodyCtxKey, msg.Body)
			ctx = context.WithValue(ctx, cryptographer.PublicKeyCtxKey, hex.EncodeToString(msg.PublicKey[:]))

			next(w, r.WithContext(ctx))
		}
	})
}

func ValidateRole() atomika.Interceptor {
	return interceptorFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("ValidateRole unhandled. Coming Soon!")
			next(w, r)
		}
	})
}

type interceptorFunc func(http.HandlerFunc) http.HandlerFunc

func (f interceptorFunc) Intercept(next http.HandlerFunc) http.HandlerFunc {
	return f(next)
}
