package courseService

import (
	"errors"
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
		t.Run(v.name, func(t *testing.T) {
			mockRepo.FindManyFnInvoked = false
			mockRepo.FindManyFn = func(query bson.M) (error, []eduboard.Course) {
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
		t.Run(v.name, func(t *testing.T) {
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

func TestCourseService_GetCoursesByMember(t *testing.T) {
	t.Parallel()

	var mockCourseRepo mock.CourseRepository
	service := CourseService{CR: &mockCourseRepo}

	course1 := eduboard.Course{ID: "1", Title: "Course 1", Members: []string{"1"}}

	testCases := []struct {
		name     string
		input    string
		error    bool
		expected []eduboard.Course
	}{
		{"success", "1", false, []eduboard.Course{course1}},
		{"error", "", true, []eduboard.Course{}},
	}

	mockCourseRepo.FindManyFn = func(update bson.M) (error, []eduboard.Course) {
		switch update["members"] {
		case "1":
			return nil, []eduboard.Course{course1}
		default:
			return errors.New("Error fetching courses"), []eduboard.Course{}
		}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockCourseRepo.FindManyFnInvoked = false

			err, courses := service.GetCoursesByMember(v.input)
			assert.True(t, mockCourseRepo.FindManyFnInvoked, "courseRepository call was not invoked")

			assert.Equal(t, v.expected, courses, "courses do not equal expected values")

			if v.error {
				assert.Error(t, err, "did not return error when expected")
				return
			}
			assert.Nil(t, err, "returned error when it shouldn't")
		})
	}
}

func TestCourseService_GetMembers(t *testing.T) {
	t.Parallel()

	var mockCourseRepo mock.CourseRepository
	var mockUserRepo mock.UserRepository

	service := CourseService{CR: &mockCourseRepo}

	course1 := eduboard.Course{ID: "1", Title: "Course 1", Members: []string{"1"}}
	course2 := eduboard.Course{ID: "2", Title: "Course 1", Members: []string{"1", "2"}}

	user1 := eduboard.User{ID: "1", Name: "User 1"}

	testCases := []struct {
		name     string
		input    string
		error    bool
		expected []eduboard.User
	}{
		{"success", "1", false, []eduboard.User{user1}},
		{"course not found", "", true, []eduboard.User{}},
		{"error fetching members", "2", true, []eduboard.User{}},
	}

	mockCourseRepo.FindFn = func(id string) (error, eduboard.Course) {
		switch id {
		case "1":
			return nil, course1
		case "2":
			return nil, course2
		default:
			return errors.New("Error fetching courses"), eduboard.Course{}
		}
	}

	mockUserRepo.FindMembersFn = func(members []string) (error, []eduboard.User) {
		// Just looking for array length, since there is no easy way to compare slices
		if len(members) == 1 {
			return nil, []eduboard.User{user1}
		}
		return errors.New("error fetching members"), []eduboard.User{}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockCourseRepo.FindFnInvoked = false

			err, users := service.GetMembers(v.input, &mockUserRepo)
			assert.True(t, mockCourseRepo.FindFnInvoked, "courseRepository call was not invoked")

			assert.Equal(t, v.expected, users, "courses do not equal expected values")

			if v.error {
				assert.Error(t, err, "did not return error when expected")
				return
			}
			assert.Nil(t, err, "returned error when it shouldn't")
		})
	}
}

func TestCourseService_AddMembers(t *testing.T) {
	t.Parallel()

	var mockCourseRepo mock.CourseRepository
	service := CourseService{CR: &mockCourseRepo}

	course1 := eduboard.Course{ID: "1", Title: "Course 1", Members: []string{"1"}}

	testCases := []struct {
		name         string
		courseInput  string
		membersInput []string
		error        bool
		expected     eduboard.Course
	}{
		{"success", "1", []string{"1", "2"}, false, course1},
		{"error", "", []string{""}, true, eduboard.Course{}},
	}

	mockCourseRepo.UpdateFn = func(course string, query bson.M) (error, eduboard.Course) {
		switch len(query["$push"].(bson.M)["members"].(bson.M)["$each"].([]string)) {
		case 2:
			return nil, course1
		default:

			return errors.New("Error fetching courses"), eduboard.Course{}
		}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockCourseRepo.UpdateFnInvoked = false

			err, courses := service.AddMembers(v.courseInput, v.membersInput)
			assert.True(t, mockCourseRepo.UpdateFnInvoked, "courseRepository call was not invoked")

			assert.Equal(t, v.expected, courses, "courses do not equal expected values")

			if v.error {
				assert.Error(t, err, "did not return error when expected")
				return
			}
			assert.Nil(t, err, "returned error when it shouldn't")
		})
	}
}

func TestCoursesService_RemoveMembers(t *testing.T) {
	t.Parallel()

	var mockCourseRepo mock.CourseRepository
	service := CourseService{CR: &mockCourseRepo}

	course1 := eduboard.Course{ID: "1", Title: "Course 1", Members: []string{"1"}}

	testCases := []struct {
		name         string
		courseInput  string
		membersInput []string
		error        bool
		expected     eduboard.Course
	}{
		{"success", "1", []string{"1", "2"}, false, course1},
		{"error", "", []string{""}, true, eduboard.Course{}},
	}

	mockCourseRepo.UpdateFn = func(course string, query bson.M) (error, eduboard.Course) {
		switch len(query["$pull"].(bson.M)["members"].(bson.M)["$in"].([]string)) {
		case 2:
			return nil, course1
		default:
			return errors.New("Error fetching courses"), eduboard.Course{}
		}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockCourseRepo.UpdateFnInvoked = false

			err, courses := service.RemoveMembers(v.courseInput, v.membersInput)
			assert.True(t, mockCourseRepo.UpdateFnInvoked, "courseRepository call was not invoked")

			assert.Equal(t, v.expected, courses, "courses do not equal expected values")

			if v.error {
				assert.Error(t, err, "did not return error when expected")
				return
			}
			assert.Nil(t, err, "returned error when it shouldn't")
		})
	}
}
