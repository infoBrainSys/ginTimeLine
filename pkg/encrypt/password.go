package encrypt

import "golang.org/x/crypto/bcrypt"

type Password struct {
}

func NewEncryptPassword() *Password {
	return new(Password)
}

func (*Password) Encrypt(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

func (*Password) ComparePassword(hash []byte, pass string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pass))
}
