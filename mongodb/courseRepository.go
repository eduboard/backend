package mongodb

import (
	"errors"
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (c *CourseRepository) Store(course *eduboard.Course) error {
	return nil
}

func (c *CourseRepository) Find(id string) (error, *eduboard.Course) {
	result := eduboard.Course{}

	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), &eduboard.Course{}
	}

	if err := c.c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return err, &eduboard.Course{}
	}
	return nil, &result
}

func (c *CourseRepository) FindAll() (error, []*eduboard.Course) {
	result := []*eduboard.Course{}

	if err := c.c.Find(bson.M{}).All(&result); err != nil {
		return err, []*eduboard.Course{}
	}

	return nil, result
}
