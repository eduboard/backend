package courseService

import (
	"github.com/eduboard/backend"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type CourseService struct {
	CR eduboard.CourseRepository
}

func New(repository eduboard.CourseRepository) CourseService {
	return CourseService{
		CR: repository,
	}
}

func (cS CourseService) GetAllCourses() (error, []eduboard.Course) {
	return cS.CR.FindMany(bson.M{})
}

func (cS CourseService) GetCourse(id string, cef eduboard.CourseEntryManyFinder) (error, eduboard.Course) {
	err, course := cS.CR.FindOneByID(id)
	if err != nil {
		return errors.Wrapf(err, "error finding course %s", id), eduboard.Course{}
	}

	if len(course.EntryIDs) == 0 {
		return nil, course
	}

	err, e := cef.FindMany(bson.M{"courseID": course.ID})
	if err != nil {
		return errors.Wrapf(err, "error finding courseEntries from %d", course.ID), eduboard.Course{}
	}

	course.Entries = e
	return nil, course
}

func (cS CourseService) GetCoursesByMember(id string, cef eduboard.CourseEntryManyFinder) (error, []eduboard.Course) {
	err, courses := cS.CR.FindMany(bson.M{"members": id})
	if err != nil {
		return errors.Wrapf(err, "error finding courses %s", id), []eduboard.Course{}
	}

	for _, course := range courses {
		if len(course.EntryIDs) > 0 {
			err, e := cef.FindMany(bson.M{"courseID": course.ID})
			if err != nil {
				return errors.Wrapf(err, "error finding courseEntries from %s", course.ID), []eduboard.Course{}
			}
			course.Entries = e
		}
	}

	return nil, courses
}

func (cS CourseService) GetMembers(id string, uF eduboard.UserFinder) (error, []eduboard.User) {
	err, course := cS.CR.FindOneByID(id)
	if err != nil {
		return errors.Wrapf(err, "error finding course %s", id), []eduboard.User{}
	}

	err, users := uF.FindMembers(course.Members)
	if err != nil {
		return errors.Wrapf(err, "error finding members from course %s", id), []eduboard.User{}
	}

	return nil, users
}

func (cS CourseService) AddMembers(course string, members []string) (error, eduboard.Course) {
	return cS.CR.Update(course, bson.M{"$push": bson.M{"members": bson.M{"$each": members}}})
}

func (cS CourseService) RemoveMembers(course string, members []string) (error, eduboard.Course) {
	return cS.CR.Update(course, bson.M{"$pull": bson.M{"members": bson.M{"$in": members}}})
}

func (cS CourseService) CreateCourse(c *eduboard.Course) (*eduboard.Course, error) {
	err := cS.CR.Insert(c)
	if err != nil {
		return &eduboard.Course{}, errors.Wrap(err, "error storing course")
	}
	return c, nil
}
