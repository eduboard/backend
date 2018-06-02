package auth

import "github.com/satori/go.uuid"

func (Authenticator) SessionID() string {
	v4 := uuid.NewV4()
	return v4.String()
}
