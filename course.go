package eduboard

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Course struct {
	ID          bson.ObjectId   `json:"id,omitempty" bson:"_id"`
	Title       string          `json:"title,omitempty" bson:"title,omitempty"`
	Description string          `json:"description,omitempty" bson:"description,omitempty"`
	Members     []string        `json:"members,omitempty" bson:"members,omitempty"`
	CreatedAt   time.Time       `json:"createdAt" bson:"createdAt"`
	Labels      []string        `json:"labels" bson:"labels"`
	EntryIDs    []bson.ObjectId `json:"entryIDs" bson:"entryIDs"`
	Entries     []CourseEntry   `json:"entries" bson:"entries"`
}

type CourseInserter interface {
	Insert(course *Course) error
}

type CourseOneFinder interface {
	FindOneByID(id string) (error, Course)
}

type CourseManyFinder interface {
	FindMany(query bson.M) (error, []Course)
}

type CourseUpdater interface {
	Update(id string, update bson.M) (error, Course)
}

type CourseFindUpdater interface {
	CourseOneFinder
	CourseUpdater
}

type CourseRepository interface {
	CourseOneFinder
	CourseManyFinder
	CourseUpdater
}

type CourseService interface {
	GetAllCourses() (err error, courses []Course)
	GetCourse(id string, cef CourseEntryManyFinder) (err error, course Course)
	GetCoursesByMember(id string, cef CourseEntryManyFinder) (err error, courses []Course)
	GetMembers(id string, uF UserFinder) (error, []User)
	AddMembers(id string, members []string) (error, Course)
	RemoveMembers(id string, members []string) (error, Course)
}
