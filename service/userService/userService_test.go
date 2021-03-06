package userService

import (
	"errors"
	"fmt"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var r = mock.UserRepository{
	FindFnInvoked: false,
	FindFn: func(id string) (error, eduboard.User) {
		if id == "0" {
			return nil, eduboard.User{ID: "0"}
		}
		return errors.New("not found"), eduboard.User{}
	},
	FindManyFnInvoked: false,
	FindManyFn: func(query bson.M) ([]eduboard.User, error) {
		return []eduboard.User{{ID: "0"}}, nil
	},
	FindByEmailFnInvoked: false,
	FindByEmailFn: func(email string) (error, eduboard.User) {
		if email == "existing@mail.com" {
			return nil, eduboard.User{ID: "0", PasswordHash: "password"}
		}
		return errors.New("not found"), eduboard.User{}
	},
	FindBySessionIDFnInvoked: false,
	FindBySessionIDFn: func(sessionID string) (error, eduboard.User) {
		if sessionID == "sessionID-0-0-0" {
			return nil, eduboard.User{ID: "1"}
		}
		return errors.New("not found"), eduboard.User{}
	},
	StoreFnInvoked: false,
	StoreFn: func(user *eduboard.User) error {
		if user.Email == "fail@mail.com" {
			return errors.New("error storing user")
		}
		return nil
	},
	UpdateSessionIDFnInvoked: false,
	UpdateSessionIDFn: func(user eduboard.User) (error, eduboard.User) {
		return nil, user
	},
}
var a = mock.AuthenticatorMock{
	HashFnInvoked: false,
	HashFn: func(password string) (string, error) {
		if password == "longpasswordbuthashfailed" {
			return "", errors.New("Error hashing password")
		}
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
		{"error hashing password", "new@mail.com", "longpasswordbuthashfailed", true},
		{"error storing user", "fail@mail.com", "longpassword", true},
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
				if err.Error() != "error storing user: error storing user" {
					assert.False(t, r.StoreFnInvoked, "Store was invoked")
				}
				assert.Equal(t, eduboard.User{}, user, "did not return empty user")
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
				assert.Equal(t, eduboard.User{}, user, "did not return empty user")
				assert.True(t, r.FindFnInvoked, "Find was not invoked")
				return
			}
			assert.Nil(t, err, "failed to get existing user")
			assert.NotEqual(t, eduboard.User{}, user, "returned non-empty user")
			assert.True(t, r.FindFnInvoked, "Find was not invoked")
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		name  string
		error bool
	}{
		{"success", false},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			defer func() { r.FindManyFnInvoked = false }()

			users, err := us.GetAllUsers()

			assert.Nil(t, err, "failed to get users")
			assert.NotEqual(t, []eduboard.User{}, users, "returned non empty users")
			assert.True(t, r.FindManyFnInvoked, "FindMany was not invoked")
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
				assert.Equal(t, eduboard.User{}, user, "did not return empty user")
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

func TestUserService_GetMyCourses(t *testing.T) {
	t.Parallel()
	var mockUserRepo mock.UserRepository
	var mockCourseRepo mock.CourseRepository
	var mockCourseEntryRepo mock.CourseEntryRepository

	service := UserService{r: &mockUserRepo}

	course1 := eduboard.Course{ID: "1", Title: "Course 1", Members: []string{"1"}}
	courseWithEntries := eduboard.Course{ID: "2", Title: "Course 2", Members: []string{"2"}, EntryIDs: []bson.ObjectId{"1", "2"}}
	brokenCourseWithEntries := eduboard.Course{ID: "4", Title: "Course 4", EntryIDs: []bson.ObjectId{"1"}}
	courseEntry1 := eduboard.CourseEntry{ID: "1", CourseID: "2", Message: "First Entry", Published: true}
	courseEntry2 := eduboard.CourseEntry{ID: "2", CourseID: "2", Message: "Second Entry", Published: true}

	courseWithEntries.Entries = []eduboard.CourseEntry{courseEntry1, courseEntry2}

	var testCases = []struct {
		name     string
		input    string
		error    bool
		expected []eduboard.Course
	}{
		{"success", "1", false, []eduboard.Course{course1}},
		{"invalid id", "", true, []eduboard.Course{}},
		{"valid id but user not found", "", true, []eduboard.Course{}},
		{"success not empty", "2", false, []eduboard.Course{courseWithEntries}},
		{"unexpected repository error", "4", true, []eduboard.Course{}},
	}

	mockUserRepo.IsIDValidFn = func(id string) bool {
		switch id {
		case "":
			return false
		case "1", "2", "3", "4":
			return true
		default:
			return false
		}
	}

	mockCourseRepo.FindManyFn = func(query bson.M) (error, []eduboard.Course) {
		switch query["members"] {
		case "1":
			return nil, []eduboard.Course{course1}
		case "2":
			return nil, []eduboard.Course{courseWithEntries}
		case "4":
			return nil, []eduboard.Course{brokenCourseWithEntries}
		default:
			return errors.New("not found"), []eduboard.Course{}
		}
	}

	mockCourseEntryRepo.FindManyFn = func(query bson.M) (error, []eduboard.CourseEntry) {
		if string(query["courseID"].(bson.ObjectId)) == "4" {
			return errors.New("error"), nil
		}
		return nil, []eduboard.CourseEntry{courseEntry1, courseEntry2}
	}

	for _, v := range testCases {
		t.Run(fmt.Sprintf("FindByMember: %s", v.name), func(t *testing.T) {
			mockUserRepo.IsIDValidFnInvoked = false
			mockCourseRepo.FindManyFnInvoked = false

			err, courses := service.GetMyCourses(v.input, &mockCourseRepo, &mockCourseEntryRepo)
			assert.True(t, mockUserRepo.IsIDValidFnInvoked, "user repository call was not invoked")

			assert.Equal(t, v.expected, courses, "courses do not equal expected value")

			if v.error {
				assert.Error(t, err, "did not return error when expected")
				return
			}
			assert.Nil(t, err, "returned error when it shouldn't")
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
