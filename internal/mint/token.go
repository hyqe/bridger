package mint

import (
	"encoding/base64"
	"errors"
	"strings"
)

type Token struct {
	header    []byte
	payload   []byte
	signature []byte
}

func (t Token) String() string {
	return base64.RawURLEncoding.EncodeToString(t.header) + "." +
		base64.RawURLEncoding.EncodeToString(t.payload) + "." +
		base64.RawURLEncoding.EncodeToString(t.signature)
}

func (t Token) IsValid(key []byte) bool {
	return validate(key, append(t.header, t.payload...), t.signature)
}

var (
	errInvalidToken = errors.New("invalid token format")
)

func Parse(v string) (Token, error) {
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return Token{}, errInvalidToken
	}
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Token{}, errors.Join(errInvalidToken, err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Token{}, errors.Join(errInvalidToken, err)
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Token{}, errors.Join(errInvalidToken, err)
	}
	return Token{
		header:    header,
		payload:   payload,
		signature: signature,
	}, nil
}
