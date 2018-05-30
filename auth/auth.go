package auth

import (
	"log"
	"golang.org/x/crypto/bcrypt"
)

func (a *UserAuthenticator)HashAndSalt(pwd string) (error string) {
	bytePwd := []byte(pwd)

	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func (a *UserAuthenticator)CompareWithHash(plainPwd string, hashedPwd string) (error bool) {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
