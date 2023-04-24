package mint

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"

	"golang.org/x/crypto/blake2b"
)

type Minter interface {
	Mint(v any) Token
}

type MinterFunc func(v any) Token

func (fn MinterFunc) Mint(v any) Token {
	return fn(v)
}

func NewMinter(secret []byte) Minter {
	return MinterFunc(func(v any) Token {
		b, _ := json.Marshal(v)
		return Token{
			header:    hash(secret),
			payload:   b,
			signature: sign(secret, b),
		}
	})
}

func hash(v []byte) []byte {
	h, _ := blake2b.New256(nil)
	return h.Sum(v)
}

func sign(secret, msg []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(msg)
	return mac.Sum(nil)
}

func validate(secret, message, signature []byte) bool {
	return hmac.Equal(signature, sign(secret, message))
}
