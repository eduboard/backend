package service

import (
	"errors"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var r = mock.UserRepository{
	FindFnInvoked: false,
	FindFn: func(id string) (error, *eduboard.User) {
		if id == "0" {
			return nil, &eduboard.User{Id: "0"}
		}
		return errors.New("not found"), &eduboard.User{}
	},
	FindByEmailFnInvoked: false,
	FindByEmailFn: func(email string) (error, *eduboard.User) {
		if email == "existing@mail.com" {
			return nil, &eduboard.User{Id: "0", PasswordHash: "password"}
		}
		return errors.New("not found"), &eduboard.User{}
	},
	FindBySessionIdFnInvoked: false,
	FindBySessionIdFn: func(sessionId string) (error, *eduboard.User) {
		if sessionId == "sessionID-0-0-0" {
			return nil, &eduboard.User{Id: "1"}
		}
		return errors.New("not found"), &eduboard.User{}
	},
	StoreFnInvoked: false,
	StoreFn: func(user *eduboard.User) error {
		return nil
	},
	UpdateSessionIdFnInvoked: false,
	UpdateSessionIdFn: func(user *eduboard.User) (error, *eduboard.User) {
		return nil, user
	},
}
var a = mock.AuthenticatorMock{
	HashFnInvoked: false,
	HashFn: func(password string) (string, error) {
		return password, nil
	},
	CompareHashFnInvoked: false,
	CompareHashFn: func(hashedPassword string, plainPassword string) (bool, error) {
		if hashedPassword != plainPassword {
			return false, errors.New("incorrect password")
		}
		return true, nil
	},
	SessionIdFnInvoked: false,
	SessionIdFn: func() string {
		return "sessionID-0-0-0"
	},
}

var us = &UserService{&r, &a}

func TestNewUserService(t *testing.T) {
	t.Parallel()
	u := NewUserService(&r)
	assert.Equal(t, &r, u.r, "repository does not match")
	assert.NotNil(t, u.a, "no authenticator")
}

func TestUserService_CreateUser(t *testing.T) {
	var testCases = []struct {
		name     string
		email    string
		password string
		error    bool
	}{
		{"user exists", "existing@mail.com", "password", true},
		{"password too short", "new@mail.com", "pass", true},
		{"new user", "new@mail.com", "longpassword", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindByEmailFnInvoked = false }()
			defer func() { a.SessionIdFnInvoked = false }()
			defer func() { r.StoreFnInvoked = false }()
			defer func() { a.HashFnInvoked = false }()

			err, user := us.CreateUser(&eduboard.User{Email: v.email}, v.password)
			if v.error {
				assert.NotNil(t, err, "did not fail to create existing user")
				assert.False(t, r.StoreFnInvoked, "Store was invoked")
				assert.Equal(t, &eduboard.User{}, user, "did not return empty user")
				return
			}

			assert.Nil(t, err, "should not fail to create user")
			assert.True(t, r.FindByEmailFnInvoked, "FindByEmail was not invoked")
			assert.True(t, a.HashFnInvoked, "Hash was not invoked")
			assert.Equal(t, v.password, user.PasswordHash, "did not hash password")
			assert.True(t, a.SessionIdFnInvoked, "SessionId was not invoked")
			assert.Equal(t, "sessionID-0-0-0", user.SessionId, "sessionId was not set")
			assert.True(t, r.StoreFnInvoked, "Store was not invoked")
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	var testCases = []struct {
		name  string
		id    string
		error bool
	}{
		{"user exists", "0", false},
		{"new user", "1", true},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindFnInvoked = false }()

			err, user := us.GetUser(v.id)
			if v.error {
				assert.NotNil(t, err, "did not fail to get non-existing user")
				assert.Equal(t, &eduboard.User{}, user, "did not return empty user")
				assert.True(t, r.FindFnInvoked, "Find was not invoked")
				return
			}
			assert.Nil(t, err, "failed to get existing user")
			assert.NotEqual(t, &eduboard.User{}, user, "returned non-empty user")
			assert.True(t, r.FindFnInvoked, "Find was not invoked")
		})
	}
}

func TestUserService_Login(t *testing.T) {
	var testCases = []struct {
		name     string
		email    string
		password string
		error    bool
	}{
		{"new user", "new@mail.com", "password", true},
		{"incorrect password", "exiisting@mail.com", "pass2", true},
		{"user exists", "existing@mail.com", "password", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindByEmailFnInvoked = false }()
			defer func() { a.CompareHashFnInvoked = false }()
			defer func() { a.SessionIdFnInvoked = false }()
			defer func() { r.UpdateSessionIdFnInvoked = false }()

			err, user := us.Login(v.email, v.password)
			if v.error {
				assert.NotNil(t, err, "did not fail to log in user")
				assert.Equal(t, &eduboard.User{}, user, "did not return empty user")
				assert.False(t, r.UpdateSessionIdFnInvoked, "UpdateSessionId invoked")
				return
			}
			assert.Nil(t, err, "should not fail to login user")
			assert.True(t, r.FindByEmailFnInvoked, "FindByEmail not invoked")
			assert.True(t, a.CompareHashFnInvoked, "CompareHash not invoked")
			assert.True(t, a.SessionIdFnInvoked, "SessionId not invoked")
			assert.True(t, r.UpdateSessionIdFnInvoked, "UpdateSessionId not invoked")
			assert.Equal(t, "sessionID-0-0-0", user.SessionId, "did not update sessionId")
		})
	}
}

func TestUserService_Logout(t *testing.T) {
	var testCases = []struct {
		name      string
		sessionId string
		error     bool
	}{
		{"unknown session", "someOtherSession", true},
		{"user exists", "sessionID-0-0-0", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindBySessionIdFnInvoked = false }()
			defer func() { r.UpdateSessionIdFnInvoked = false }()

			err := us.Logout(v.sessionId)
			if v.error {
				assert.NotNil(t, err, "did not fail")
				assert.True(t, r.FindBySessionIdFnInvoked, "FindBySessionId was not invoked")
				assert.False(t, r.UpdateSessionIdFnInvoked, "UpdateSessionIdFn was invoked")
				return
			}
			assert.Nil(t, err, "caused error logging out user")
			assert.True(t, r.FindBySessionIdFnInvoked, "FindBySessionId was not invoked")
			assert.True(t, r.UpdateSessionIdFnInvoked, "UpdateSessionId was not invoked")
		})
	}
}

func TestUserService_CheckAuthentication(t *testing.T) {
	var testCases = []struct {
		name      string
		sessionId string
		error     bool
	}{
		{"unknown session", "someOtherSession", true},
		{"user exists", "sessionID-0-0-0", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindBySessionIdFnInvoked = false }()

			err, ok := us.CheckAuthentication(v.sessionId)
			if v.error {
				assert.NotNil(t, err, "did not fail")
				assert.False(t, ok, "should not be ok")
				assert.True(t, r.FindBySessionIdFnInvoked, "UpdateSessionIdFn was invoked")
				return
			}
			assert.Nil(t, err, "caused error logging out user")
			assert.True(t, ok, "should be ok")
			assert.True(t, r.FindBySessionIdFnInvoked, "FindBySessionId was not invoked")
		})
	}
}
