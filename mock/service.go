package mock

import (
	"github.com/eduboard/backend"
	"time"
)

type CourseService struct {
	CourseFn        func(id string, cef eduboard.CourseEntryManyFinder) (err error, course eduboard.Course)
	CourseFnInvoked bool

	CoursesFn        func() (err error, courses []eduboard.Course)
	CoursesFnInvoked bool

	GetCoursesByMemberFn        func(id string) (error, []eduboard.Course)
	GetCoursesByMemberFnInvoked bool

	GetMembersFn        func(course string, uF eduboard.UserFinder) (error, []eduboard.User)
	GetMembersFnInvoked bool

	AddMembersFn        func(course string, members []string) (error, eduboard.Course)
	AddMembersFnInvoked bool

	RemoveMembersFn        func(course string, members []string) (error, eduboard.Course)
	RemoveMembersFnInvoked bool
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

func (cSM *CourseService) GetCoursesByMember(id string) (error, []eduboard.Course) {
	cSM.GetCoursesByMemberFnInvoked = true
	return cSM.GetCoursesByMemberFn(id)
}

func (cSM *CourseService) GetMembers(course string, uF eduboard.UserFinder) (error, []eduboard.User) {
	cSM.GetMembersFnInvoked = true
	return cSM.GetMembersFn(course, uF)
}

func (cSM *CourseService) AddMembers(course string, members []string) (error, eduboard.Course) {
	cSM.AddMembersFnInvoked = true
	return cSM.AddMembersFn(course, members)
}

func (cSM *CourseService) RemoveMembers(course string, members []string) (error, eduboard.Course) {
	cSM.RemoveMembersFnInvoked = true
	return cSM.RemoveMembersFn(course, members)
}

type CourseEntryService struct {
	StoreCourseEntryFn        func(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater) (err error, courseEntry *eduboard.CourseEntry)
	StoreCourseEntryFnInvoked bool

	StoreCourseEntryFilesFn        func(files [][]byte, id string, date time.Time) ([]string, error)
	StoreCourseEntryFilesFnInvoked bool

	UpdateCourseEntryFn        func(*eduboard.CourseEntry) (*eduboard.CourseEntry, error)
	UpdateCourseEntryFnInvoked bool

	DeleteCourseEntryFn        func(entryID string, courseID string, updater eduboard.CourseUpdater) error
	DeleteCourseEntryFnInvoked bool
}

func (cSM *CourseEntryService) StoreCourseEntry(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater) (err error, courseEntry *eduboard.CourseEntry) {
	cSM.StoreCourseEntryFnInvoked = true
	return cSM.StoreCourseEntryFn(entry, cfu)
}

func (cSM *CourseEntryService) StoreCourseEntryFiles(files [][]byte, id string, date time.Time) ([]string, error) {
	cSM.StoreCourseEntryFilesFnInvoked = true
	return cSM.StoreCourseEntryFilesFn(files, id, date)
}

func (cSM *CourseEntryService) UpdateCourseEntry(e *eduboard.CourseEntry) (*eduboard.CourseEntry, error) {
	cSM.UpdateCourseEntryFnInvoked = true
	return cSM.UpdateCourseEntryFn(e)
}

func (cSM *CourseEntryService) DeleteCourseEntry(entryID string, courseID string, updater eduboard.CourseUpdater) error {
	cSM.DeleteCourseEntryFnInvoked = true
	return cSM.DeleteCourseEntryFn(entryID, courseID, updater)
}

type UserService struct {
	CreateUserFn        func(u *eduboard.User, password string) (error, eduboard.User)
	CreateUserFnInvoked bool

	GetUserFn        func(id string) (error, eduboard.User)
	GetUserFnInvoked bool

	GetMyCoursesFn        func(id string, cBMF eduboard.CourseManyFinder) (error, []eduboard.Course)
	GetMyCoursesFnInvoked bool

	UserAuthenticationProvider
}

var _ eduboard.UserService = (*UserService)(nil)

func (uSM *UserService) CreateUser(u *eduboard.User, password string) (error, eduboard.User) {
	uSM.CreateUserFnInvoked = true
	return uSM.CreateUserFn(u, password)
}

func (uSM *UserService) GetUser(id string) (error, eduboard.User) {
	uSM.GetUserFnInvoked = true
	return uSM.GetUserFn(id)
}

func (uSM *UserService) GetMyCourses(id string, cBMF eduboard.CourseManyFinder) (error, []eduboard.Course) {
	uSM.GetMyCoursesFnInvoked = true
	return uSM.GetMyCoursesFn(id, cBMF)
}

type UserAuthenticationProvider struct {
	LoginFn        func(email string, password string) (error, eduboard.User)
	LoginFnInvoked bool

	LogoutFn        func(sessionID string) error
	LogoutFnInvoked bool

	CheckAuthenticationFn        func(sessionID string) (err error, userID string)
	CheckAuthenticationFnInvoked bool
}

var _ eduboard.UserAuthenticationProvider = (*UserAuthenticationProvider)(nil)

func (uAM *UserAuthenticationProvider) Login(email string, password string) (error, eduboard.User) {
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
