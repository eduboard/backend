package auth

import "github.com/satori/go.uuid"

func (Authenticator) SessionID() string {
	return uuid.NewV4().String()
}
