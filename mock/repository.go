package mock

import (
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	UserRepository   UserRepository
	CourseRepository CourseRepository
}

// CourseRepository implements the eduboard.CourseRepository interface to mock functions and record successful invocations.
type CourseRepository struct {
	StoreFn        func(course *eduboard.Course) error
	StoreFnInvoked bool

	FindFn        func(id string) (error, eduboard.Course)
	FindFnInvoked bool

	FindManyFn        func() (error, []eduboard.Course)
	FindManyFnInvoked bool

	UpdateFn        func(id string, update bson.M) (error, eduboard.Course)
	UpdateFnInvoked bool
}

var (
	_ eduboard.CourseRepository = (*CourseRepository)(nil)
	_ eduboard.CourseInserter   = (*CourseRepository)(nil)
	_ eduboard.CourseOneFinder  = (*CourseRepository)(nil)
	_ eduboard.CourseManyFinder = (*CourseRepository)(nil)
	_ eduboard.CourseUpdater    = (*CourseRepository)(nil)
)

func (cRM *CourseRepository) Insert(course *eduboard.Course) error {
	cRM.StoreFnInvoked = true
	return cRM.StoreFn(course)
}

func (cRM *CourseRepository) FindOneByID(id string) (error, eduboard.Course) {
	cRM.FindFnInvoked = true
	return cRM.FindFn(id)
}

func (cRM *CourseRepository) FindMany(query bson.M) (error, []eduboard.Course) {
	cRM.FindManyFnInvoked = true
	return cRM.FindManyFn()
}

func (cRM *CourseRepository) Update(id string, update bson.M) (error, eduboard.Course) {
	cRM.UpdateFnInvoked = true
	return cRM.UpdateFn(id, update)
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

func (uRM *UserRepository) FindByEmail(email string) (error, *eduboard.User) {
	uRM.FindByEmailFnInvoked = true
	return uRM.FindByEmailFn(email)
}

func (uRM *UserRepository) FindBySessionID(sessionID string) (error, *eduboard.User) {
	uRM.FindBySessionIDFnInvoked = true
	return uRM.FindBySessionIDFn(sessionID)
}

func (uRM *UserRepository) UpdateSessionID(user *eduboard.User) (error, *eduboard.User) {
	uRM.UpdateSessionIDFnInvoked = true
	return uRM.UpdateSessionIDFn(user)
}

// CourseEntryRepository implements the eduboard.CourseEntryRepository interface to mock functions and record successful invocations.
type CourseEntryRepository struct {
	InsertFn        func(course eduboard.CourseEntry) error
	InsertFnInvoked bool

	FindOneFn        func(id string) (error, eduboard.CourseEntry)
	FindOneFnInvoked bool

	FindManyFn        func(query bson.M) (error, []eduboard.CourseEntry)
	FindManyFnInvoked bool

	UpdateFn        func(id string, update bson.M) (error, eduboard.CourseEntry)
	UpdateFnInvoked bool

	DeleteFn        func(id string) error
	DeleteFnInvoked bool
}

var _ eduboard.CourseEntryRepository = (*CourseEntryRepository)(nil)

func (cRM *CourseEntryRepository) Insert(course eduboard.CourseEntry) error {
	cRM.InsertFnInvoked = true
	return cRM.InsertFn(course)
}

func (cRM *CourseEntryRepository) FindOneByID(id string) (error, eduboard.CourseEntry) {
	cRM.FindOneFnInvoked = true
	return cRM.FindOneFn(id)
}

func (cRM *CourseEntryRepository) FindMany(query bson.M) (error, []eduboard.CourseEntry) {
	cRM.FindManyFnInvoked = true
	return cRM.FindManyFn(query)
}

func (cRM *CourseEntryRepository) Update(id string, update bson.M) (error, eduboard.CourseEntry) {
	cRM.UpdateFnInvoked = true
	return cRM.UpdateFn(id, update)
}

func (cRM *CourseEntryRepository) Delete(id string) error {
	cRM.DeleteFnInvoked = true
	return cRM.DeleteFn(id)
}
