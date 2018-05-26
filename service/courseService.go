package service

import (
	"github.com/eduboard/backend"
)

type CourseService struct {
	r edubord.CourseRepository
}

func NewCourseService(repository edubord.CourseRepository) *CourseService {
	return &CourseService{
		r: repository,
	}
}

func (c *CourseService) GetAllCourses() (err error, courses []*edubord.Course) {
	return nil, []*edubord.Course{}
}

func (c *CourseService) GetCourse(id edubord.CourseId) (err error, course *edubord.Course) {
	return nil, &edubord.Course{}
}
