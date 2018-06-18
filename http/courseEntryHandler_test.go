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
)

func TestAppServer_PostCourseEntryHandler(t *testing.T) {
	var testCases = []struct {
		name        string
		input       string
		courseID    string
		invokeStore bool
		status      int
	}{
		{"success", `{"message": "success"}"`, "5b23bbdc2bfa844c41a9f135", true, 200},
		{"bad json", `{"message":`, "5b23bbdc2bfa844c41a9f135", false, 400},
		{"bad objectid", `{"message": "success"}"`, "5b23bbdc2bfa844c41a9f35", false, 400},
		{"bad urls", `{"message": "success", "pictures": ["htttp\\:.orgcom"]}"`, "5b23bbdc2bfa844c41a9f135", false, 400},
		{"error storing", `{"message": "success"}"`, "5b23bbdc2bfa844c41a9f136", true, 500},
	}

	service := mock.CourseEntryService{}
	service.StoreCourseEntryFn = func(entry *eduboard.CourseEntry, cfu eduboard.CourseFindUpdater) (err error, courseEntry *eduboard.CourseEntry) {
		if entry.CourseID.Hex() == "5b23bbdc2bfa844c41a9f136" {
			return errors.New("could not store"), &eduboard.CourseEntry{}
		}
		return nil, entry
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
