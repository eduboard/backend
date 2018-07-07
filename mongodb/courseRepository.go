package mongodb

import (
	"errors"
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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

func (c *CourseRepository) Insert(course *eduboard.Course) error {
	if course.ID == "" {
		course.ID = bson.NewObjectId()
	}

	now := time.Now().UTC()
	course.CreatedAt = now

	return c.c.Insert(course)
}

func (c *CourseRepository) FindOneByID(id string) (error, eduboard.Course) {
	result := eduboard.Course{}

	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), eduboard.Course{}
	}

	if err := c.c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return err, eduboard.Course{}
	}
	return nil, result
}

func (c *CourseRepository) FindMany(query bson.M) (error, []eduboard.Course) {
	result := []eduboard.Course{}

	if err := c.c.Find(bson.M{}).All(&result); err != nil {
		return err, []eduboard.Course{}
	}

	return nil, result
}

func (c *CourseRepository) Update(id string, update bson.M) (error, eduboard.Course) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), eduboard.Course{}
	}

	if err := c.c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, update); err != nil {
		return err, eduboard.Course{}
	}

	return nil, eduboard.Course{}
}
