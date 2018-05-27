package service

import (
	"github.com/eduboard/backend"
)

type CourseService struct {
	r eduboard.CourseRepository
}

func NewCourseService(repository eduboard.CourseRepository) *CourseService {
	return &CourseService{
		r: repository,
	}
}

func (cS *CourseService) GetAllCourses() (err error, courses []*eduboard.Course) {
	return cS.r.FindAll()
}

func (cS *CourseService) GetCourse(id string) (err error, course *eduboard.Course) {
	return cS.r.Find(id)
}
