package courseEntryService

import (
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()
	r := mock.CourseEntryRepository{}
	ces := New(&r)
	assert.Equal(t, &r, ces.ER, "repository does not match")
}

func TestCourseEntryService_StoreCourseEntry(t *testing.T) {
	successEntry := eduboard.CourseEntry{CourseID: bson.ObjectIdHex("5b23c8d5382d33000150681e")}
	failureEntry := eduboard.CourseEntry{CourseID: bson.ObjectIdHex("5b23c8d5382d33000150681f")}

	var testCases = []struct {
		name  string
		error bool
		entry eduboard.CourseEntry
	}{
		{"success", false, successEntry},
		{"no course", true, failureEntry},
	}

	mockEntryRepo := mock.CourseEntryRepository{}
	service := CourseEntryService{ER: &mockEntryRepo}

	mockEntryRepo.InsertFn = func(course eduboard.CourseEntry) error {
		return nil
	}

	mockCourseRepo := mock.CourseRepository{}
	mockCourseRepo.FindFn = func(id string) (error, eduboard.Course) {
		if id == "5b23c8d5382d33000150681e" {
			return nil, eduboard.Course{ID: bson.ObjectIdHex("5b23c8d5382d33000150681e")}
		}
		return errors.New("not found"), eduboard.Course{}
	}
	mockCourseRepo.UpdateFn = func(id string, update bson.M) (error, eduboard.Course) {
		return nil, eduboard.Course{}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockEntryRepo.InsertFnInvoked = false
			mockCourseRepo.FindFnInvoked = false
			mockCourseRepo.UpdateFnInvoked = false

			err, e := service.StoreCourseEntry(&v.entry, &mockCourseRepo)
			if v.error {
				assert.Equal(t, &eduboard.CourseEntry{}, e, "courseEntry unexpected")
				assert.NotNil(t, err, "error nil")
				assert.True(t, mockCourseRepo.FindFnInvoked, "FindOne was not invoked")
				assert.False(t, mockCourseRepo.UpdateFnInvoked, "update was invoked")
				assert.False(t, mockEntryRepo.InsertFnInvoked, "insert was invoked")
				return
			}
			assert.Nil(t, err, "error not nil")
			assert.True(t, mockCourseRepo.FindFnInvoked, "FindOne was not invoked")
			assert.True(t, mockCourseRepo.UpdateFnInvoked, "update was not invoked")
			assert.True(t, mockEntryRepo.InsertFnInvoked, "insert was not invoked")
			assert.True(t, bson.IsObjectIdHex(e.ID.Hex()), "new entry does not have valid ID")
			return
		})
	}
}

func TestCourseEntryService_StoreCourseEntryFiles(t *testing.T) {
	var testCases = []struct {
		name   string
		error  bool
		files  [][]byte
		course string
	}{
		{"success", false, [][]byte{[]byte("test")}, "success"},
		{"error", true, [][]byte{[]byte("moep")}, "failure"},
	}
	uploader := mock.Uploader{}
	uploader.UploadFileFn = func(file []byte, course string, filename string) (error, string) {
		if course == "success" {
			return nil, filename
		}
		return errors.New("error uploading files"), ""
	}
	service := CourseEntryService{u: &uploader}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			uploader.UploadFileFnInvoked = false

			err, _ := service.StoreCourseEntryFiles(v.files, v.course, time.Now())
			if v.error {
				assert.NotNil(t, err, "error nil")
				return
			}
			assert.Nil(t, err, "error not nil")

		})
	}
}

func TestCourseEntryService_DeleteCourseEntry(t *testing.T) {
	successEntry := "5b23c8d5382d33000150681e"
	failureEntry := "5b23c8d5382d33000150681f"
	successCourse := "5b23c8d5382d33000150681a"
	failureCourse := "5b23c8d5382d33000150681b"

	var testCases = []struct {
		name   string
		error  bool
		entry  string
		course string
	}{
		{"success", false, successEntry, successCourse},
		{"no entry", true, failureEntry, successCourse},
		{"no course", true, successEntry, failureCourse},
	}

	mockEntryRepo := mock.CourseEntryRepository{}
	service := CourseEntryService{ER: &mockEntryRepo}
	mockEntryRepo.FindOneFn = func(id string) (error, eduboard.CourseEntry) {
		if id == successEntry {
			return nil, eduboard.CourseEntry{ID: bson.ObjectIdHex(successEntry), CourseID: bson.ObjectIdHex(successCourse)}
		}
		return errors.New("not found"), eduboard.CourseEntry{}
	}
	mockEntryRepo.DeleteFn = func(id string) error { return nil }

	mockCourseRepo := mock.CourseRepository{}
	mockCourseRepo.UpdateFn = func(id string, update bson.M) (error, eduboard.Course) { return nil, eduboard.Course{} }

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockEntryRepo.FindOneFnInvoked = false
			mockEntryRepo.DeleteFnInvoked = false
			mockCourseRepo.UpdateFnInvoked = false

			err := service.DeleteCourseEntry(v.entry, v.course, &mockCourseRepo)
			assert.True(t, mockEntryRepo.FindOneFnInvoked, "FindeOne was not invoked")
			if v.error {
				assert.Errorf(t, err, "error is nil")
				assert.False(t, mockCourseRepo.UpdateFnInvoked, "Update was invoked")
				assert.False(t, mockEntryRepo.DeleteFnInvoked, "Delete was invoked")
				return
			}
			assert.Nil(t, err, "error not nil")
			assert.True(t, mockEntryRepo.DeleteFnInvoked, "Delete was not invoked")
			assert.True(t, mockCourseRepo.UpdateFnInvoked, "Update was not invoked")
		})
	}
}

func TestCourseEntryService_UpdateCourseEntry(t *testing.T) {

}
