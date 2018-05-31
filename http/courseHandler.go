package http

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (a *AppServer) getAllCoursesHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err, courses := a.CourseService.GetAllCourses()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(courses); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) getCourseHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p[0].Value
		err, course := a.CourseService.GetCourse(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(course); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
