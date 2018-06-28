package http

import (
	"encoding/json"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/url"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func (a *AppServer) PostCourseEntryHandler() httprouter.Handle {
	type request struct {
		Date      time.Time `json:"date"`
		Message   string    `json:"message"`
		Pictures  [][]byte  `json:"pictures"`
		Published bool      `json:"published"`
	}
	type response struct {
		ID        string    `json:"id"`
		Date      time.Time `json:"date"`
		Message   string    `json:"message"`
		Pictures  []string  `json:"pictures"`
		Published bool      `json:"published"`
	}
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			entryModel eduboard.CourseEntry
			request    request
		)
		id := p.ByName("courseID")
		if ok := bson.IsObjectIdHex(id); !ok {
			a.Logger.Printf("courseID %s is not a valid objectID", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			a.Logger.Printf("error decoding request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err, paths := a.CourseEntryService.StoreCourseEntryFiles(request.Pictures, id, request.Date)
		if err != nil {
			a.Logger.Printf("error storing files: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pURLs, err := url.URLifyStrings(paths)
		if err != nil {
			a.Logger.Printf("error parsing urls: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		entryModel.CourseID = bson.ObjectIdHex(id)
		entryModel.Date = request.Date
		entryModel.CreatedAt = time.Now()
		entryModel.Message = request.Message
		entryModel.Pictures = pURLs
		entryModel.Published = request.Published

		err, entry := a.CourseEntryService.StoreCourseEntry(&entryModel, a.CourseRepository)
		if err != nil {
			a.Logger.Printf("error storing courseEntry: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := response{
			ID:        entry.ID.Hex(),
			Date:      entry.Date,
			Message:   entry.Message,
			Pictures:  url.StringifyURLs(entry.Pictures),
			Published: entry.Published,
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) PutCourseEntryHandler() httprouter.Handle {
	type request struct {
		Date      time.Time `json:"date"`
		Message   string    `json:"message"`
		Pictures  []string  `json:"pictures"`
		Published bool      `json:"published"`
	}
	type response struct {
		ID        string    `json:"id"`
		Date      time.Time `json:"date"`
		Message   string    `json:"message"`
		Pictures  []string  `json:"pictures"`
		Published bool      `json:"published"`
	}
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (a *AppServer) DeleteCourseEntryHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		courseID := p.ByName("courseID")
		entryID := p.ByName("entryID")
		err := a.CourseEntryService.DeleteCourseEntry(entryID, courseID, a.CourseRepository)
		if err != nil {
			a.Logger.Printf("error deleting courseEntry: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
