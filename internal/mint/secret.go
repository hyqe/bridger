package mint

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

func NewSecret(length int) Secret {
	out := make(Secret, length)
	rand.Read(out)
	return out
}

type Secret []byte

func (s Secret) String() string {
	return s.B64()
}

func (s Secret) B64() string {
	return base64.RawURLEncoding.EncodeToString(s)
}

func (s Secret) Hex() string {
	return hex.EncodeToString(s)
}
