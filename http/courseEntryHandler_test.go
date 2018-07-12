package http

import (
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAppServer_PostCourseEntryHandler(t *testing.T) {
	var testCases = []struct {
		name        string
		input       string
		courseID    string
		invokeStore bool
		status      int
	}{
		{"success", `{"message": "success", "pictures": ["c3VwZXJzZWNyZXQ=\n"]}"`, "5b23bbdc2bfa844c41a9f135", true, 200},
		{"bad json", `{"message":`, "5b23bbdc2bfa844c41a9f135", false, 400},
		{"bad objectid", `{"message": "success"}"`, "5b23bbdc2bfa844c41a9f35", false, 400},
		{"unparsable url", `{"message": "success", "pictures": ["c3VwZXJzZWNyZXQ=\n", "c3VwZXJzZWNyZXQ=\n"]}"`, "5b23bbdc2bfa844c41a9f135", false, 500},
		{"error storing files", `{"message": "success", "pictures": ["c3VwZXJzZWNyZXQ=\n", "c3VwZXJzZWNyZXQ=\n", "c3VwZXJzZWNyZXQ=\n"]}"`, "5b23bbdc2bfa844c41a9f135", false, 500},
		{"error storing", `{"message": "success"}"`, "5b23bbdc2bfa844c41a9f136", true, 500},
	}

	service := mock.CourseEntryService{}
	service.StoreCourseEntryFn = func(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater) (err error, courseEntry *eduboard.CourseEntry) {
		if entry.CourseID.Hex() == "5b23bbdc2bfa844c41a9f136" {
			return errors.New("could not store"), &eduboard.CourseEntry{}
		}
		return nil, entry
	}

	service.StoreCourseEntryFilesFn = func(files [][]byte, id string, date time.Time) ([]string, error) {
		switch len(files) {
		case 0:
			return []string{}, nil
		case 1:
			return []string{"/test/test/1.jpg"}, nil
		case 2:
			return []string{":test.com/cant/parse/this//"}, nil
		default:
			return []string{}, errors.New("invalid something error while storing files")
		}
	}

	a := AppServer{CourseEntryService: &service, Logger: log.New(os.Stdout, "", 0)}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			service.StoreCourseEntryFnInvoked = false
			r := httptest.NewRequest("POST", "/", strings.NewReader(v.input))
			rr := httptest.NewRecorder()
			p := httprouter.Params{httprouter.Param{Key: "courseID", Value: v.courseID}}

			a.PostCourseEntryHandler()(rr, r, p)
			assert.Equal(t, v.status, rr.Code, "status code does not match")
			assert.Equal(t, v.invokeStore, service.StoreCourseEntryFnInvoked, "Store was not invoked unexpectedly")
		})
	}
}
