package auth

import "golang.org/x/crypto/bcrypt"

type Authenticator struct{}

func (a *Authenticator) Hash(password string) (string, error) {
	bytePwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (a *Authenticator) CompareHash(hashedPassword string, plainPassword string) (bool, error) {
	byteHash := []byte(hashedPassword)
	bytePlain := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false, err
	}

	return true, nil
}
