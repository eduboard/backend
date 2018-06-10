package mock

import (
	"github.com/eduboard/backend"
)

// CourseService implements the eduboard.UserService interface to mock functions and record successful invocations.
type CourseService struct {
	CourseFn        func(id string) (err error, course *eduboard.Course)
	CourseFnInvoked bool

	CoursesFn          func() (err error, courses []*eduboard.Course)
	CoursesFuncInvoked bool
}

// Statically check that CourseService actually implements eduboard.CourseService interface
var _ eduboard.CourseService = (*CourseService)(nil)

func (cSM *CourseService) GetCourse(id string) (err error, course *eduboard.Course) {
	cSM.CourseFnInvoked = true
	return cSM.CourseFn(id)
}

func (cSM *CourseService) GetAllCourses() (err error, courses []*eduboard.Course) {
	cSM.CoursesFuncInvoked = true
	return cSM.CoursesFn()
}

type UserAuthenticationProvider struct {
	LoginFn        func(email string, password string) (error, *eduboard.User)
	LoginFnInvoked bool

	LogoutFn        func(sessionID string) error
	LogoutFnInvoked bool

	CheckAuthenticationFn        func(sessionID string) (err error, userID string)
	CheckAuthenticationFnInvoked bool
}

// Statically check that UserAuthenticationProvider actually implements eduboard.UserAuthenticationProvider interface
var _ eduboard.UserAuthenticationProvider = (*UserAuthenticationProvider)(nil)

func (uAM *UserAuthenticationProvider) Login(email string, password string) (error, *eduboard.User) {
	uAM.LoginFnInvoked = true
	return uAM.LoginFn(email, password)
}

func (uAM *UserAuthenticationProvider) Logout(sessionID string) error {
	uAM.LogoutFnInvoked = true
	return uAM.LogoutFn(sessionID)
}
func (uAM *UserAuthenticationProvider) CheckAuthentication(sessionID string) (err error, userID string) {
	uAM.CheckAuthenticationFnInvoked = true
	return uAM.CheckAuthenticationFn(sessionID)
}

type Authenticator interface {
	Hash(password string) (string, error)
	CompareHash(hashedPassword string, plainPassword string) (bool, error)
	SessionID() string
}

type AuthenticatorMock struct {
	HashFn        func(password string) (string, error)
	HashFnInvoked bool

	CompareHashFn        func(hashedPassword string, plainPassword string) (bool, error)
	CompareHashFnInvoked bool

	SessionIDFn        func() string
	SessionIDFnInvoked bool
}

var _ Authenticator = (*AuthenticatorMock)(nil)

func (uAM *AuthenticatorMock) Hash(password string) (string, error) {
	uAM.HashFnInvoked = true
	return uAM.HashFn(password)
}

func (uAM *AuthenticatorMock) CompareHash(hashedPassword string, plainPassword string) (bool, error) {
	uAM.CompareHashFnInvoked = true
	return uAM.CompareHashFn(hashedPassword, plainPassword)
}
func (uAM *AuthenticatorMock) SessionID() string {
	uAM.SessionIDFnInvoked = true
	return uAM.SessionIDFn()
}
