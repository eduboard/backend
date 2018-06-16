package mock

import (
	"github.com/eduboard/backend"
)

type CourseService struct {
	CourseFn        func(id string, cef eduboard.CourseEntryManyFinder) (err error, course eduboard.Course)
	CourseFnInvoked bool

	CoursesFn        func() (err error, courses []eduboard.Course)
	CoursesFnInvoked bool
}

var _ eduboard.CourseService = (*CourseService)(nil)

func (cSM *CourseService) GetCourse(id string, cef eduboard.CourseEntryManyFinder) (err error, course eduboard.Course) {
	cSM.CourseFnInvoked = true
	return cSM.CourseFn(id, cef)
}

func (cSM *CourseService) GetAllCourses() (err error, courses []eduboard.Course) {
	cSM.CoursesFnInvoked = true
	return cSM.CoursesFn()
}

type CourseEntryService struct {
	StoreCorseEntryFn        func(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater, ces eduboard.CourseEntryInserter) (err error, courseEntry *eduboard.CourseEntry)
	StoreCorseEntryFnInvoked bool
}

func (cSM *CourseEntryService) StoreCourseEntry(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater, ces eduboard.CourseEntryInserter) (err error, courseEntry *eduboard.CourseEntry) {
	cSM.StoreCorseEntryFnInvoked = true
	return cSM.StoreCorseEntryFn(entry, cfu, ces)
}

type UserService struct {
	CreateUserFn        func(u *eduboard.User, password string) (error, *eduboard.User)
	CreateUserFnInvoked bool

	GetUserFn        func(id string) (error, *eduboard.User)
	GetUserFnInvoked bool

	UserAuthenticationProvider
}

var _ eduboard.UserService = (*UserService)(nil)

func (uSM *UserService) CreateUser(u *eduboard.User, password string) (error, *eduboard.User) {
	uSM.CreateUserFnInvoked = true
	return uSM.CreateUserFn(u, password)
}

func (uSM *UserService) GetUser(id string) (error, *eduboard.User) {
	uSM.GetUserFnInvoked = true
	return uSM.GetUserFn(id)
}

type UserAuthenticationProvider struct {
	LoginFn        func(email string, password string) (error, *eduboard.User)
	LoginFnInvoked bool

	LogoutFn        func(sessionID string) error
	LogoutFnInvoked bool

	CheckAuthenticationFn        func(sessionID string) (err error, userID string)
	CheckAuthenticationFnInvoked bool
}

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
