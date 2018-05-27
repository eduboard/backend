package eduboard

import "gopkg.in/mgo.v2/bson"

type Course struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
	Members     []string      `json:"members,omitempty" bson:"members,omitempty"`
}

type CourseRepository interface {
	Store(course *Course) error
	Find(id string) (error, *Course)
	FindAll() (error, []*Course)
}

type CourseService interface {
	GetAllCourses() (err error, courses []*Course)
	GetCourse(id string) (err error, course *Course)
}
