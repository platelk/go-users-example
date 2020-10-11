package pwdhasher

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultBcryptCost = 12

// Bcrypt provide an
type Bcrypt struct {
	cost int
}

// NewBcrypt will instantiate a bcrypt hasher
func NewBcrypt() *Bcrypt {
	return &Bcrypt{cost: defaultBcryptCost}
}

// Hash will generate a securely crypted password
func (b *Bcrypt) Hash(pwd string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(pwd), b.cost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}


