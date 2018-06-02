package http

import (
	"encoding/json"
	"fmt"
	"github.com/eduboard/backend"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (a *AppServer) registerUserHandler() httprouter.Handle {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID    string `json:"id"`
		Email string `json:"email"`
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
		err, user := a.UserService.CreateUser(&userModel, request.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := response{ID: user.ID.Hex(), Email: user.Email}

		cookie := createCookie(user.SessionID)
		http.SetCookie(w, &cookie)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) loginUserHandler() httprouter.Handle {
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
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err, user := a.UserService.Login(request.Email, request.Password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := createCookie(user.SessionID)
		response := response{user.Name, user.Surname, user.Email}
		http.SetCookie(w, &cookie)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (a *AppServer) logoutUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sessionID, err := r.Cookie("sessionID")
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		err = a.UserService.Logout(sessionID.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := createCookie("")
		http.SetCookie(w, &cookie)
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

func createCookie(value string) http.Cookie {
	expire := time.Now().Add(24 * time.Hour)
	return http.Cookie{Name: "sessionID", Value: value, Path: "/", Expires: expire, MaxAge: 86400}
}
