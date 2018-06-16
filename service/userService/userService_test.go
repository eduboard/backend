package userService

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
			return nil, &eduboard.User{ID: "0"}
		}
		return errors.New("not found"), &eduboard.User{}
	},
	FindByEmailFnInvoked: false,
	FindByEmailFn: func(email string) (error, *eduboard.User) {
		if email == "existing@mail.com" {
			return nil, &eduboard.User{ID: "0", PasswordHash: "password"}
		}
		return errors.New("not found"), &eduboard.User{}
	},
	FindBySessionIDFnInvoked: false,
	FindBySessionIDFn: func(sessionID string) (error, *eduboard.User) {
		if sessionID == "sessionID-0-0-0" {
			return nil, &eduboard.User{ID: "1"}
		}
		return errors.New("not found"), &eduboard.User{}
	},
	StoreFnInvoked: false,
	StoreFn: func(user *eduboard.User) error {
		return nil
	},
	UpdateSessionIDFnInvoked: false,
	UpdateSessionIDFn: func(user *eduboard.User) (error, *eduboard.User) {
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
	SessionIDFnInvoked: false,
	SessionIDFn: func() string {
		return "sessionID-0-0-0"
	},
}

var us = &UserService{&r, &a}

func TestNew(t *testing.T) {
	t.Parallel()
	u := New(&r)
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
			defer func() { a.SessionIDFnInvoked = false }()
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
			assert.True(t, a.SessionIDFnInvoked, "SessionID was not invoked")
			assert.Equal(t, "sessionID-0-0-0", user.SessionID, "sessionID was not set")
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
			defer func() { a.SessionIDFnInvoked = false }()
			defer func() { r.UpdateSessionIDFnInvoked = false }()

			err, user := us.Login(v.email, v.password)
			if v.error {
				assert.NotNil(t, err, "did not fail to log in user")
				assert.Equal(t, &eduboard.User{}, user, "did not return empty user")
				assert.False(t, r.UpdateSessionIDFnInvoked, "UpdateSessionID invoked")
				return
			}
			assert.Nil(t, err, "should not fail to login userService")
			assert.True(t, r.FindByEmailFnInvoked, "FindByEmail not invoked")
			assert.True(t, a.CompareHashFnInvoked, "CompareHash not invoked")
			assert.True(t, a.SessionIDFnInvoked, "SessionID not invoked")
			assert.True(t, r.UpdateSessionIDFnInvoked, "UpdateSessionID not invoked")
			assert.Equal(t, "sessionID-0-0-0", user.SessionID, "did not update sessionID")
		})
	}
}

func TestUserService_Logout(t *testing.T) {
	var testCases = []struct {
		name      string
		sessionID string
		error     bool
	}{
		{"unknown session", "someOtherSession", true},
		{"user exists", "sessionID-0-0-0", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindBySessionIDFnInvoked = false }()
			defer func() { r.UpdateSessionIDFnInvoked = false }()

			err := us.Logout(v.sessionID)
			if v.error {
				assert.NotNil(t, err, "did not fail")
				assert.True(t, r.FindBySessionIDFnInvoked, "FindBySessionID was not invoked")
				assert.False(t, r.UpdateSessionIDFnInvoked, "UpdateSessionIDFn was invoked")
				return
			}
			assert.Nil(t, err, "caused error logging out user")
			assert.True(t, r.FindBySessionIDFnInvoked, "FindBySessionID was not invoked")
			assert.True(t, r.UpdateSessionIDFnInvoked, "UpdateSessionID was not invoked")
		})
	}
}

func TestUserService_CheckAuthentication(t *testing.T) {
	var testCases = []struct {
		name      string
		sessionID string
		error     bool
	}{
		{"unknown session", "someOtherSession", true},
		{"user exists", "sessionID-0-0-0", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindBySessionIDFnInvoked = false }()

			err, id := us.CheckAuthentication(v.sessionID)
			if v.error {
				assert.NotNil(t, err, "did not fail")
				assert.Equal(t, "", id, "should not contain id")
				assert.True(t, r.FindBySessionIDFnInvoked, "UpdateSessionIDFn was invoked")
				return
			}
			assert.Nil(t, err, "caused error logging out user")
			assert.Equal(t, "31", id, "should contain id") // 31 is the Hex representation of 1.
			assert.True(t, r.FindBySessionIDFnInvoked, "FindBySessionID was not invoked")
		})
	}
}
