package http

import (
	"encoding/json"
	"fmt"
	"github.com/eduboard/backend"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (a *AppServer) RegisterUserHandler() httprouter.Handle {
	type request struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			userModel eduboard.User
			request   request
		)

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil || request.Email == "" || request.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userModel.Email = request.Email
		userModel.Name = request.Name
		userModel.Surname = request.Surname

		err, user := a.UserService.CreateUser(&userModel, request.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := response{
			ID:      user.ID.Hex(),
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
		}

		cookie := createCookie(user.SessionID, user.SessionExpires)
		http.SetCookie(w, &cookie)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) LoginUserHandler() httprouter.Handle {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var request request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil || request.Email == "" || request.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err, user := a.UserService.Login(request.Email, request.Password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := createCookie(user.SessionID, user.SessionExpires)
		response := response{user.Name, user.Surname, user.Email}
		http.SetCookie(w, &cookie)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) LogoutUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sessionID, err := r.Cookie("sessionID")
		if err != nil || sessionID.Value == "" {
			w.WriteHeader(http.StatusOK)
			return
		}

		err = a.UserService.Logout(sessionID.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := createCookie("", time.Time{})
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}

func (a *AppServer) GetUserHandler() httprouter.Handle {
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

func (a *AppServer) GetMeHandler() httprouter.Handle {
	type response struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := r.Header.Get("userID")
		err, user := a.UserService.GetUser(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response := response{
			ID:      user.ID.Hex(),
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func createCookie(value string, expires time.Time) http.Cookie {
	return http.Cookie{Name: "sessionID", Value: value, Path: "/", Expires: expires, MaxAge: 86400}
}
