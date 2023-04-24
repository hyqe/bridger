package mint

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
)

type Minter interface {
	Mint(v any) Token
}

type MinterFunc func(v any) Token

func (fn MinterFunc) Mint(v any) Token {
	return fn(v)
}

func NewMinter(k []byte) Minter {
	return MinterFunc(func(v any) Token {
		b, _ := json.Marshal(v)
		return Token{
			header:    hash(k),
			payload:   b,
			signature: sign(k, b),
		}
	})
}

func hash(v []byte) []byte {
	return sha256.New().Sum(v)
}

func sign(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

func validate(key, message, signature []byte) bool {
	return hmac.Equal(signature, sign(key, message))
}
