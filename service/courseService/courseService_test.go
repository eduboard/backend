package courseService

import (
	"errors"
	"fmt"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	r := mock.CourseRepository{}
	cs := New(&r)
	assert.Equal(t, &r, cs.CR, "repository does not match")
}

func TestCourseService_GetAllCourses(t *testing.T) {
	t.Parallel()
	var mockRepo mock.CourseRepository

	service := CourseService{CR: &mockRepo}

	course1 := eduboard.Course{ID: "1", Title: "Course 1"}
	course2 := eduboard.Course{ID: "2", Title: "Course 2"}

	var testCases = []struct {
		name     string
		input    []eduboard.Course
		error    bool
		expected []eduboard.Course
	}{
		{"success", []eduboard.Course{course1, course2}, false, []eduboard.Course{course1, course2}},
		{"empty", []eduboard.Course{}, false, []eduboard.Course{}},
		{"error", []eduboard.Course{}, true, []eduboard.Course{}},
	}

	for _, v := range testCases {
		t.Run(fmt.Sprintf("GetAllCourses %s", v.name), func(t *testing.T) {
			mockRepo.FindManyFnInvoked = false
			mockRepo.FindManyFn = func() (error, []eduboard.Course) {
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
			assert.True(t, mockRepo.FindManyFnInvoked, "repository call was not invoked")
			assert.Equal(t, v.expected, courses, "courses do not equal expected values")
		})
	}
}

func TestCourseService_GetCourse(t *testing.T) {
	t.Parallel()
	var mockCourseRepo mock.CourseRepository
	var mockEntryRepo mock.CourseEntryRepository

	service := CourseService{CR: &mockCourseRepo}

	course := eduboard.Course{ID: "1", Title: "Course 1"}
	courseWithEntries := eduboard.Course{ID: "2", Title: "Course 2", EntryIDs: []bson.ObjectId{"1", "2"}}
	courseEntry1 := eduboard.CourseEntry{ID: "1", CourseID: "2", Message: "First Entry", Published: true}
	courseEntry2 := eduboard.CourseEntry{ID: "2", CourseID: "2", Message: "Second Entry", Published: true}

	expectedCourse := courseWithEntries
	expectedCourse.Entries = []eduboard.CourseEntry{courseEntry1, courseEntry2}

	var testCases = []struct {
		name            string
		input           string
		error           bool
		invokeEntryRepo bool
		expected        eduboard.Course
	}{
		{"success empty", "1", false, false, course},
		{"success", "2", false, true, expectedCourse},
		{"empty", "3", true, false, eduboard.Course{}},
		{"error", "", true, false, eduboard.Course{}},
	}

	mockCourseRepo.FindFn = func(id string) (error, eduboard.Course) {
		switch id {
		case "":
			return errors.New("no id"), eduboard.Course{}
		case "1":
			return nil, course
		case "2":
			return nil, courseWithEntries
		default:
			return errors.New("not found"), eduboard.Course{}
		}
	}

	mockEntryRepo.FindManyFn = func(query bson.M) (error, []eduboard.CourseEntry) {
		return nil, []eduboard.CourseEntry{courseEntry1, courseEntry2}
	}

	for _, v := range testCases {
		t.Run(fmt.Sprintf("GetAllCourses %s", v.name), func(t *testing.T) {
			mockCourseRepo.FindFnInvoked = false
			mockEntryRepo.FindManyFnInvoked = false

			err, course := service.GetCourse(v.input, &mockEntryRepo)
			assert.True(t, mockCourseRepo.FindFnInvoked, "courseRepository call was not invoked")
			assert.Equal(t, v.invokeEntryRepo, mockEntryRepo.FindManyFnInvoked, "entryRepository call was not invoked as expected")
			assert.Equal(t, v.expected, course, "courses do not equal expected values")

			if v.error {
				assert.Error(t, err, "did not return error when expected")
				return
			}
			assert.Nil(t, err, "returned error when it should not")
		})
	}
}
