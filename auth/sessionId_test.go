package auth

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticator_CreateAccessToken(t *testing.T) {
	t.Parallel()
	authenticator := Authenticator{}
	sessionId := authenticator.SessionId()

	id, err := uuid.FromString(sessionId)
	assert.Nil(t, err, "should not cause error")
	assert.Equal(t, byte(4), id.Version(), "should create UUID v4")
}
