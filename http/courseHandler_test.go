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
	"strings"
	"testing"
	"time"
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
		name       string
		id         string
		bodylength int
		status     int
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
			appServer.GetCourseHandler()(rr, r, httprouter.Params{httprouter.Param{"courseID", v.id}})

			assert.True(t, mockService.CourseFnInvoked, "GetCourse was not invoked")
			assert.Equal(t, v.status, rr.Code, "status code unexpected")
			assert.Equal(t, v.bodylength, len(rr.Body.String()), "body length does not match")
		})
	}
}

func TestAppServer_GetMembersHandler(t *testing.T) {
	mockService := mock.CourseService{}
	appServer := AppServer{CourseService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name       string
		id         string
		status     int
		bodylength int
	}{
		{"success", "1", 200, 54},
		{"error", "2", 404, 0},
	}

	mockService.GetMembersFn = func(id string, uR eduboard.UserFinder) (error, []eduboard.User) {
		switch id {
		case "1":
			return nil, []eduboard.User{{ID: "1", Name: "User 1"}}
		default:
			return errors.New("Error fetching members"), []eduboard.User{}
		}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.GetMembersFnInvoked = false

			r := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			appServer.GetMembersHandler()(rr, r, httprouter.Params{httprouter.Param{"courseID", v.id}})

			assert.True(t, mockService.GetMembersFnInvoked, "GetMembers was not invoked")
			assert.Equal(t, v.status, rr.Code, "unexpected status code")
			assert.Equal(t, v.bodylength, len(rr.Body.String()), "body length does not equal expected value")
		})
	}
}

func TestAppServer_AddMembersHandler(t *testing.T) {
	mockService := mock.CourseService{}
	appServer := AppServer{CourseService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name   string
		id     string
		body   string
		status int
	}{
		{"success", "1", `[{"id": "1"},{"id": "2"}]`, 200},
		{"invalid body", "2", `lala <> ""`, 400},
		{"error subscribing user", "3", `[{"id": "1"},{"id": "2"}]`, 500},
	}

	mockService.AddMembersFn = func(id string, members []string) (error, eduboard.Course) {
		switch id {
		case "1":
			if len(members) != 2 {
				return errors.New("Members are not correct"), eduboard.Course{}
			}
			return nil, eduboard.Course{ID: "1"}
		default:
			return errors.New("Error fetching members"), eduboard.Course{}
		}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.AddMembersFnInvoked = false

			r := httptest.NewRequest("POST", "/", strings.NewReader(v.body))
			rr := httptest.NewRecorder()
			appServer.AddMembersHandler()(rr, r, httprouter.Params{httprouter.Param{"courseID", v.id}})

			if v.status != 400 {
				assert.True(t, mockService.AddMembersFnInvoked, "AddMembers was not invoked")
			}
			assert.Equal(t, v.status, rr.Code, "unexpected status code")
		})
	}
}

func TestAppServer_RemoveMembersHandler(t *testing.T) {
	mockService := mock.CourseService{}
	appServer := AppServer{CourseService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name   string
		id     string
		body   string
		status int
	}{
		{"success", "1", `[{"id": "1"},{"id": "2"}]`, 200},
		{"invalid body", "2", `lala <> ""`, 400},
		{"error subscribing user", "3", `[{"id": "1"},{"id": "2"}]`, 500},
	}

	mockService.RemoveMembersFn = func(id string, members []string) (error, eduboard.Course) {
		switch id {
		case "1":
			return nil, eduboard.Course{ID: "1"}
		default:
			return errors.New("Error fetching members"), eduboard.Course{}
		}
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.RemoveMembersFnInvoked = false

			r := httptest.NewRequest("POST", "/", strings.NewReader(v.body))
			rr := httptest.NewRecorder()
			appServer.RemoveMembersHandler()(rr, r, httprouter.Params{httprouter.Param{"courseID", v.id}})

			if v.status != 400 {
				assert.True(t, mockService.RemoveMembersFnInvoked, "RemoveMembers was not invoked")
			}
			assert.Equal(t, v.status, rr.Code, "unexpected status code")
		})
	}
}

func TestAppServer_CreateCourseHandler(t *testing.T) {
	mockService := mock.CourseService{}
	appServer := AppServer{
		CourseService:         &mockService,
		CourseRepository:      &mock.CourseRepository{},
		CourseEntryRepository: &mock.CourseEntryRepository{},
		Logger:                log.New(os.Stdout, "", 0),
	}
	var testCases = []struct {
		name    string
		invoked bool
		success bool
		body    string
		code    int
	}{
		{"good case", true, true, `{"title": "course 1", "description": "course 1"}`, 200},
		{"wrong body", false, true, `{"title": "course 1, "description": "course 1"`, 400},
		{"error", true, false, `{"title": "course 1", "description": "course 1"}`, 500},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.CreateCourseFnInvoked = false
			mockService.CreateCourseFn = func(c *eduboard.Course) (*eduboard.Course, error) {
				newCourse := eduboard.Course{}
				newCourse.ID = bson.NewObjectId()
				newCourse.Title = c.Title
				newCourse.Description = c.Description
				newCourse.Members = c.Members
				newCourse.Labels = c.Labels
				now := time.Now().UTC()
				newCourse.CreatedAt = now
				if !v.success {
					return nil, errors.New("invalid id")
				}
				return &newCourse, nil
			}

			r := httptest.NewRequest("POST", "/", strings.NewReader(v.body))
			rr := httptest.NewRecorder()
			appServer.CreateCourseHandler()(rr, r, httprouter.Params{})
			assert.Equal(t, v.invoked, mockService.CreateCourseFnInvoked, "CreateCourse was not invoked")
			assert.Equal(t, v.code, rr.Code, "status code unexpected")
		})
	}
}
