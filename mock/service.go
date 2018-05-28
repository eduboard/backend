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

	LoginFn        func(username string, password string) (error, *eduboard.User)
	LoginFnInvoked bool
}

func (uSM *UserService) CreateUser(user *eduboard.User) (error, *eduboard.User) {
	uSM.CreateUserFnInvoked = true
	return uSM.CreateUserFn(user)
}

func (uSM *UserService) Login(username string, password string) (error, *eduboard.User) {
	uSM.LoginFnInvoked = true
	return uSM.LoginFn(username, password)
}
