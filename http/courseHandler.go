package http

import (
	"encoding/json"
	"github.com/eduboard/backend/url"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (a *AppServer) GetAllCoursesHandler() httprouter.Handle {
	type response struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err, courses := a.CourseService.GetAllCourses()
		if err != nil {
			a.Logger.Printf("error getting courses: %v", err)
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

func (a *AppServer) GetCourseHandler() httprouter.Handle {
	type entryResponse struct {
		ID        string    `json:"id"`
		Date      time.Time `json:"date"`
		Message   string    `json:"message"`
		Pictures  []string  `json:"pictures"`
		Published bool      `json:"published"`
	}

	type courseResponse struct {
		ID          string          `json:"id"`
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Members     []string        `json:"members,omitempty"`
		Labels      []string        `json:"labels,omitempty"`
		Entries     []entryResponse `json:"entries,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("courseID")
		err, course := a.CourseService.GetCourse(id, a.CourseEntryRepository)
		if err != nil {
			a.Logger.Printf("error getting course: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res := courseResponse{
			ID:          course.ID.Hex(),
			Title:       course.Title,
			Description: course.Description,
			Members:     course.Members,
			Labels:      course.Labels,
			Entries:     make([]entryResponse, len(course.Entries)),
		}

		for k, v := range course.Entries {
			res.Entries[k] = entryResponse{
				ID:        v.ID.Hex(),
				Date:      v.Date,
				Message:   v.Message,
				Pictures:  url.StringifyURLs(v.Pictures),
				Published: v.Published,
			}
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) GetMembersHandler() httprouter.Handle {
	type response struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("courseID")
		err, members := a.CourseService.GetMembers(id, a.UserRepository)
		if err != nil {
			a.Logger.Printf("Error getting members: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var res = make([]response, len(members))
		for k, v := range members {
			res[k] = response{
				ID:      v.ID.Hex(),
				Name:    v.Name,
				Surname: v.Surname,
				Email:   v.Email,
			}
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) AddMembersHandler() httprouter.Handle {
	type request struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		request := []request{}
		id := p.ByName("courseID")
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			a.Logger.Printf("Error decoding body %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		members := []string{}
		for _, req := range request {
			members = append(members, req.ID)
		}
		err, _ = a.CourseService.AddMembers(id, members)

		if err != nil {
			a.Logger.Printf("Error while subscribing user to course %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func (a *AppServer) RemoveMembersHandler() httprouter.Handle {
	type request struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		request := []request{}
		id := p.ByName("courseID")
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			a.Logger.Printf("Error decoding body %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		members := []string{}
		for _, req := range request {
			members = append(members, req.ID)
		}
		err, _ = a.CourseService.RemoveMembers(id, members)
		if err != nil {
			a.Logger.Printf("Error while unsubscribing user from course %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
