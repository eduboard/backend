package auth

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticator_SessionID(t *testing.T) {
	t.Parallel()
	authenticator := Authenticator{}
	sessionID := authenticator.SessionID()

	id, err := uuid.FromString(sessionID)
	assert.Nil(t, err, "should not cause error")
	assert.Equal(t, byte(4), id.Version(), "should create UUID v4")
}
