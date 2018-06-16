package http

import (
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAppServer_GetAllCoursesHandler(t *testing.T) {
	coursesList := []eduboard.Course{{
		ID:          bson.ObjectIdHex("5b23bbdc2bfa844c41a9f134"),
		Title:       "Course 1",
		Description: "Course 1 Description",
		Labels:      []string{"label 1"},
	}, {
		ID:          bson.ObjectIdHex("5b23bbdc2bfa844c41a9f135"),
		Title:       "Course 2",
		Description: "Course 2 Description",
		Labels:      []string{"label 2"},
	}}

	mockService := mock.CourseService{}
	appServer := AppServer{
		CourseService:    &mockService,
		CourseRepository: &mock.CourseRepository{},
		Logger:           log.New(os.Stdout, "", 0),
	}

	var testCases = []struct {
		name   string
		error  bool
		status int
	}{
		{"success", false, 200},
		{"error", true, 500},
	}

	for _, v := range testCases {

		mockService.CoursesFn = func() (err error, courses []eduboard.Course) {
			if v.error {
				return errors.New("error"), []eduboard.Course{}
			}
			return nil, coursesList
		}

		t.Run(v.name, func(t *testing.T) {
			mockService.CoursesFnInvoked = false
			r := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			appServer.GetAllCoursesHandler()(rr, r, httprouter.Params{})
			assert.True(t, mockService.CoursesFnInvoked, "GetAllCourses was not invoked")
			assert.Equal(t, v.status, rr.Code, "bad response code")
			assert.True(t, v.error || len(rr.Body.String()) > 2, "bad body length")
		})
	}
}

func TestAppServer_GetCourseHandler(t *testing.T) {
	mockService := mock.CourseService{}
	appServer := AppServer{
		CourseService:         &mockService,
		CourseRepository:      &mock.CourseRepository{},
		CourseEntryRepository: &mock.CourseEntryRepository{},
		Logger:                log.New(os.Stdout, "", 0),
	}

	var testCases = []struct {
		name        string
		id          string
		bodylentgth int
		status      int
	}{
		{"no entries", "1", 40, 200},
		{"success", "2", 226, 200},
		{"bad id", "3", 0, 404},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.CourseFnInvoked = false
			mockService.CourseFn = func(id string, cef eduboard.CourseEntryManyFinder) (err error, course eduboard.Course) {
				switch id {
				case "1":
					return nil, eduboard.Course{ID: "1"}
				case "2":
					return nil, eduboard.Course{
						ID: "2",
						Entries: []eduboard.CourseEntry{
							{ID: "1", CourseID: "2"},
							{ID: "2", CourseID: "2"}},
					}
				default:
					return errors.New("not found"), eduboard.Course{}

				}
			}

			r := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			appServer.GetCourseHandler()(rr, r, httprouter.Params{httprouter.Param{"id", v.id}})

			assert.True(t, mockService.CourseFnInvoked, "GetCourse was not invoked")
			assert.Equal(t, v.status, rr.Code, "status code unexpected")
			assert.Equal(t, v.bodylentgth, len(rr.Body.String()), "body length does not match")
		})
	}
}
