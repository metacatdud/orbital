package stringer

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

type StringToken string

func (st StringToken) String() string {

	return string(st)
}

const (
	RandAll       StringToken = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$^&*_-"
	RandLowercase StringToken = "abcdefghijklmnopqrstuvwxyz"
	RandUppercase StringToken = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	RandNumber    StringToken = "0123456789"
	RandSym       StringToken = "!@#$^&*_-"
)

func Random(n int, tokens ...StringToken) (string, error) {
	if n == 0 {
		return "", ErrRandTooShort
	}

	charset := RandAll
	var charsetBuf bytes.Buffer

	if len(tokens) > 0 {
		for _, t := range tokens {
			switch t {
			case RandLowercase:
				charsetBuf.WriteString(RandLowercase.String())
			case RandUppercase:
				charsetBuf.WriteString(RandUppercase.String())
			case RandNumber:
				charsetBuf.WriteString(RandNumber.String())
			case RandSym:
				charsetBuf.WriteString(RandSym.String())
			}
		}
	}

	if charsetBuf.Len() > 0 {
		charset = StringToken(charsetBuf.String())
	}

	charsetLen := big.NewInt(int64(len(charset)))
	charsetBytes := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		charsetBytes[i] = charset[num.Int64()]
	}

	return string(charsetBytes), nil
}
