package http

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (a *AppServer) getAllCoursesHandler() httprouter.Handle {
	type response struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err, courses := a.CourseService.GetAllCourses()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var res = make([]response, len(courses))
		for k, v := range courses {
			res[k] = response{ID: v.ID.Hex(), Title: v.Title, Description: v.Description}
		}

		if err = json.NewEncoder(w).Encode(&res); err != nil {
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
