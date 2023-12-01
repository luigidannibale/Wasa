package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	username := r.URL.Query().Get("username")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error taking the username")
		return
	}

	user, s, err := rt.db.GetUserByUsername(username)

	if err != nil {
		if e := err.Error(); e == "NotFound" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		e := json.NewEncoder(w).Encode(s)
		if e != nil {
			http.Error(w, "Couldn't encode error message : "+s, http.StatusInternalServerError)
		}
		return

	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w).Encode(user)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the user ", http.StatusInternalServerError)
	}

}
