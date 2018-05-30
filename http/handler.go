package http

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/eduboard/backend"
	"net/http"
	"time"
)

// USER

func (a *AppServer) registerUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var userModel eduboard.User
		err := json.NewDecoder(r.Body).Decode(&userModel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err, user := a.UserService.CreateUser(&userModel)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) loginUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var userModel eduboard.User
		err := json.NewDecoder(r.Body).Decode(&userModel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err, user := a.UserService.Login(userModel.Username, userModel.PasswordHash)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		expire := time.Now().AddDate(0, 0, 1)
		cookie := http.Cookie{Name: "accessToken", Value: user.AccessToken, Path: "/", Expires: expire, MaxAge: 86400}
		http.SetCookie(w, &cookie)
		if err = json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) logoutUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		accessToken, err := r.Cookie("accessToken")
		if err != nil {
			// "logout succeded" since we werent logged in in the first place
			w.WriteHeader(http.StatusOK)
			return
		}
		err = a.UserService.Logout(accessToken.Value)
		expire := time.Now().AddDate(0, 0, -1)
		cookie := http.Cookie{Name: "accessToken", Value: accessToken.Value, Path: "/", Expires: expire, MaxAge: 86400}
		http.SetCookie(w, &cookie)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (a *AppServer) getUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		err, user := a.UserService.GetUser(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// Course

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
