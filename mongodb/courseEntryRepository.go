package mongodb

import (
	"errors"
	"github.com/eduboard/backend"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CourseEntryRepository struct {
	c *mgo.Collection
}

func newCourseEntryRepository(database *mgo.Database) *CourseEntryRepository {
	collection := database.C("courseEntry")
	return &CourseEntryRepository{
		c: collection,
	}
}

func (c *CourseEntryRepository) Insert(course eduboard.CourseEntry) error {
	return c.c.Insert(course)
}

func (c *CourseEntryRepository) FindOneByID(id string) (error, eduboard.CourseEntry) {
	result := eduboard.CourseEntry{}

	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id"), eduboard.CourseEntry{}
	}

	if err := c.c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		return err, eduboard.CourseEntry{}
	}
	return nil, result
}

func (c *CourseEntryRepository) FindMany(query bson.M) (error, []eduboard.CourseEntry) {
	result := []eduboard.CourseEntry{}

	if err := c.c.Find(query).All(&result); err != nil {
		return err, []eduboard.CourseEntry{}
	}
	return nil, result
}

func (c *CourseEntryRepository) Update(id string, update bson.M) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id")
	}

	entryID := bson.ObjectIdHex(id)

	if err := c.c.UpdateId(entryID, update); err != nil {
		return err
	}

	return nil
}

func (c *CourseEntryRepository) Delete(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id")
	}

	return c.c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
