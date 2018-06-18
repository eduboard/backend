package courseService

import (
	"github.com/eduboard/backend"
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
		return err, eduboard.Course{}
	}

	if len(course.EntryIDs) == 0 {
		return nil, course
	}

	err, e := cef.FindMany(bson.M{"courseID": course.ID})
	if err != nil {
		return err, eduboard.Course{}
	}

	course.Entries = e
	return nil, course
}
