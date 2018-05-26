package user

import (
	"encoding/json"
	"fmt"
	"github.com/eduboard/backend/repository/user"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "<html><body>Static</bod></html>")
		register("1", "2")
	}
}

func register(user string, password string) {}

func GetUserHandler() httprouter.Handle {

	repo := user.New()
	repo.InitializeDB()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_, u := repo.FindUser(p.ByName("id"))
		json.NewEncoder(w).Encode(u)
	}
}
