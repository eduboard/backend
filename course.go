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
	Schedules   []Schedule      `json:"schedules" bson:"schedules"`
}

type Schedule struct {
	Day      time.Weekday  `json:"day" bson:"day"`
	Start    time.Time     `json:"startsAt" bson:"startsAt"`
	Duration time.Duration `json:"duration,omitempty" bson:"duration"`
	Room     string        `json:"room,omitempty" bson:"room"`
	Title    string        `json:"title,omitempty" bson:"title"`
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
	CourseInserter
	CourseOneFinder
	CourseManyFinder
	CourseUpdater
}

type CourseService interface {
	CreateCourse(c *Course) (*Course, error)
	GetAllCourses() (err error, courses []Course)
	GetCourse(id string, cef CourseEntryManyFinder) (err error, course Course)
	GetCoursesByMember(id string, cef CourseEntryManyFinder) (err error, courses []Course)
	GetMembers(id string, uF UserFinder) (error, []User)
	AddMembers(id string, members []string) (error, Course)
	RemoveMembers(id string, members []string) (error, Course)
}
