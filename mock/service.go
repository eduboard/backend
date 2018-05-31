package mock

import "github.com/eduboard/backend"

// CourseService implements the eduboard.UserService interface to mock functions and record successful invocations.
type CourseService struct {
	CourseFn        func(id string) (err error, course *eduboard.Course)
	CourseFnInvoked bool

	CoursesFn          func() (err error, courses []*eduboard.Course)
	CoursesFuncInvoked bool
}

func (cSM *CourseService) GetCourse(id string) (err error, course *eduboard.Course) {
	cSM.CourseFnInvoked = true
	return cSM.CourseFn(id)
}

func (cSM *CourseService) GetAllCourses() (err error, courses []*eduboard.Course) {
	cSM.CoursesFuncInvoked = true
	return cSM.CoursesFn()
}

// UserService implements the eduboard.UserService interface to mock functions and record successful invocations.
type UserService struct {
	CreateUserFn        func(user *eduboard.User) (error, *eduboard.User)
	CreateUserFnInvoked bool

	GetUserFn        func(id string) (error, *eduboard.User)
	GetUserFnInvoked bool
}

func (uSM *UserService) CreateUser(user *eduboard.User) (error, *eduboard.User) {
	uSM.CreateUserFnInvoked = true
	return uSM.CreateUserFn(user)
}

func (uSM *UserService) GetUser(id string) (error, *eduboard.User) {
	uSM.GetUserFnInvoked = true
	return uSM.GetUserFn(id)
}

type UserAuthenticationProvider struct {
	LoginFn        func(username string, password string) (error, *eduboard.User)
	LoginFnInvoked bool

	LogoutFn        func(accessToken string) error
	LogoutFnInvoked bool

	CheckAuthenticationFn        func(sessionId string) (err error, ok bool)
	CheckAuthenticationFnInvoked bool
}

func (uAM *UserAuthenticationProvider) Login(username string, password string) (error, *eduboard.User) {
	uAM.LoginFnInvoked = true
	return uAM.LoginFn(username, password)
}

func (uAM *UserAuthenticationProvider) Logout(sessionId string) error {
	uAM.LogoutFnInvoked = true
	return uAM.LogoutFn(sessionId)
}
func (uAM *UserAuthenticationProvider) CheckAuthentication(sessionId string) (err error, ok bool) {
	uAM.CheckAuthenticationFnInvoked = true
	return uAM.CheckAuthenticationFn(sessionId)
}
