package mint

import (
	"testing"
)

func Test_sign_val(t *testing.T) {
	secret := []byte("adsflkajhdflkjahdlkfjasdfjkad")
	message := []byte("hello")
	signed := sign(secret, message)
	if !validate(secret, message, signed) {
		t.Fail()
	}
}
