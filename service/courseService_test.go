package service

import (
	"errors"
	"fmt"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCourseService(t *testing.T) {
	t.Parallel()
	var mockRepo mock.CourseRepository

	service := NewCourseService(&mockRepo)
	assert.Equal(t, service.r, &mockRepo, "did not correctly attach repo to service")
}

func TestCourseService_GetAllCourses(t *testing.T) {
	t.Parallel()
	var mockRepo mock.CourseRepository
	var service CourseService

	service.r = &mockRepo

	course1 := &eduboard.Course{ID: "1", Name: "Course 1"}
	course2 := &eduboard.Course{ID: "2", Name: "Course 2"}

	var testCases = []struct {
		name     string
		input    []*eduboard.Course
		error    bool
		expected []*eduboard.Course
	}{
		{"success", []*eduboard.Course{course1, course2}, false, []*eduboard.Course{course1, course2}},
		{"empty", []*eduboard.Course{}, false, []*eduboard.Course{}},
		{"error", []*eduboard.Course{}, true, []*eduboard.Course{}},
	}

	for _, v := range testCases {
		t.Run(fmt.Sprintf("GetAllCourses %s", v.name), func(t *testing.T) {
			mockRepo.FindAllFnInvoked = false
			mockRepo.FindAllFn = func() (error, []*eduboard.Course) {
				var err error
				if v.error {
					err = errors.New("error")
				}
				return err, v.input
			}

			err, courses := service.GetAllCourses()
			if !v.error {
				assert.Nil(t, err, "failed when it should not")
			} else {
				assert.Error(t, err, "did not return error when expected")
			}
			assert.True(t, mockRepo.FindAllFnInvoked, "repository call was not invoked")
			assert.Equal(t, v.expected, courses, "courses do not equal expected values")
		})
	}
}

func TestCourseService_GetCourse(t *testing.T) {
	t.Parallel()
	var mockRepo mock.CourseRepository
	var service CourseService

	service.r = &mockRepo

	course1 := &eduboard.Course{ID: "1", Name: "Course 1"}

	var testCases = []struct {
		name     string
		input    string
		error    bool
		expected *eduboard.Course
	}{
		{"success", "1", false, course1},
		{"empty", "3", true, &eduboard.Course{}},
		{"error", "", true, &eduboard.Course{}},
	}

	for _, v := range testCases {
		t.Run(fmt.Sprintf("GetAllCourses %s", v.name), func(t *testing.T) {
			mockRepo.FindFnInvoked = false
			mockRepo.FindFn = func(id string) (error, *eduboard.Course) {
				if id == "" {
					return errors.New("no id"), &eduboard.Course{}
				}
				if id == "3" {
					return errors.New("not found"), &eduboard.Course{}
				}
				return nil, course1
			}

			err, course := service.GetCourse(v.input)
			if !v.error {
				assert.Nil(t, err, "returned error when it should not")
			} else {
				assert.Error(t, err, "did not return error when expected")
			}
			assert.True(t, mockRepo.FindFnInvoked, "repository call was not invoked")
			assert.Equal(t, v.expected, course, "courses do not equal expected values")
		})
	}
}
