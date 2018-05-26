package mongodb

import (
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2"
)

type CourseRepository struct {
	c *mgo.Collection
}

func newCourseRepository(database *mgo.Database) *CourseRepository {
	collection := database.C("course")
	return &CourseRepository{
		c: collection,
	}
}

func (c *CourseRepository) Store(course *edubord.Course) error {
	return nil
}

func (c *CourseRepository) Find(id edubord.CourseId) (error, *edubord.Course) {
	return nil, &edubord.Course{}
}

func (c *CourseRepository) FindAll() (error, []*edubord.Course) {
	return nil, []*edubord.Course{}
}
