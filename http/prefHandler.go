package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/eduboard/backend"
	"encoding/json"
)

func (a *AppServer) AddPrefHandler() httprouter.Handle {
	type request struct {
		CourseID  string `json:"courseID"`
		PartnerID string `json:"partnerID"`
		Value     int    `json:"value"`
	}
	type response struct {
		ID        string
		CourseID  string `json:"courseID"`
		UserID    string `json:"userID"`
		PartnerID string `json:"partnerID"`
		Value     int    `json:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var request request
		var pref eduboard.Pref

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pref.UserID = r.Header.Get("userID")
		pref.CourseID = request.CourseID
		pref.PartnerID = request.PartnerID
		pref.Value = request.Value

		newPref, err := a.PrefService.CreatePref(&pref)
		if err != nil {
			a.Logger.Printf("error creating pref: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := response{
			ID:        newPref.ID.Hex(),
			CourseID:  newPref.CourseID,
			UserID:    newPref.UserID,
			PartnerID: newPref.PartnerID,
			Value:     newPref.Value,
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
