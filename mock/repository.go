package mock

import "github.com/eduboard/backend"

type Repository struct {
	UserRepository   UserRepository
	CourseRepository CourseRepository
}

// CourseRepository implements the eduboard.CourseRepository interface to mock functions and record successful invocations.
type CourseRepository struct {
	StoreFn        func(course *eduboard.Course) error
	StoreFnInvoked bool

	FindFn        func(id string) (error, *eduboard.Course)
	FindFnInvoked bool

	FindAllFn        func() (error, []*eduboard.Course)
	FindAllFnInvoked bool
}

var _ eduboard.CourseRepository = (*CourseRepository)(nil)

func (cRM *CourseRepository) Store(course *eduboard.Course) error {
	cRM.StoreFnInvoked = true
	return cRM.StoreFn(course)
}

func (cRM *CourseRepository) Find(id string) (error, *eduboard.Course) {
	cRM.FindFnInvoked = true
	return cRM.FindFn(id)
}

func (cRM *CourseRepository) FindAll() (error, []*eduboard.Course) {
	cRM.FindAllFnInvoked = true
	return cRM.FindAllFn()
}

// Course implements the eduboard.CourseRepository interface to mock functions and record successful invocations.
type UserRepository struct {
	StoreFn        func(user *eduboard.User) error
	StoreFnInvoked bool

	FindFn        func(id string) (error, *eduboard.User)
	FindFnInvoked bool

	FindByEmailFn        func(email string) (error, *eduboard.User)
	FindByEmailFnInvoked bool

	FindBySessionIDFn        func(sessionID string) (error, *eduboard.User)
	FindBySessionIDFnInvoked bool

	UpdateSessionIDFn        func(user *eduboard.User) (error, *eduboard.User)
	UpdateSessionIDFnInvoked bool
}

var _ eduboard.UserRepository = (*UserRepository)(nil)

func (uRM *UserRepository) Store(user *eduboard.User) error {
	uRM.StoreFnInvoked = true
	return uRM.StoreFn(user)
}

func (uRM *UserRepository) Find(id string) (error, *eduboard.User) {
	uRM.FindFnInvoked = true
	return uRM.FindFn(id)
}

func (uRM *UserRepository)FindByEmail(email string) (error, *eduboard.User) {
	uRM.FindByEmailFnInvoked = true
	return uRM.FindByEmailFn(email)
}

func (uRM *UserRepository)FindBySessionID(sessionID string) (error, *eduboard.User) {
	uRM.FindBySessionIDFnInvoked = true
	return uRM.FindBySessionIDFn(sessionID)
}

func (uRM *UserRepository)UpdateSessionID(user *eduboard.User) (error, *eduboard.User) {
	uRM.UpdateSessionIDFnInvoked = true
	return uRM.UpdateSessionIDFn(user)
}
