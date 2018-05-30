package auth

import "github.com/satori/go.uuid"

type UserAuthenticator struct {}

func (aT *UserAuthenticator) CreateAccessToken() string {
	v4 := uuid.NewV4()
	return v4.String()
}
