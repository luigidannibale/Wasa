package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	userID, er := strconv.Atoi(r.Header.Get("Authorization"))
	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}
	var userIds []int
	userIds = append(userIds, userID)
	i, e := rt.db.VerifyUserIds(userIds)
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		if i == 0 {
			http.Error(w, "The userID can't be found", http.StatusNotFound)
		}
		return
	}

	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Error taking the username", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
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
	e = json.NewEncoder(w).Encode(user)
	if e != nil {
		http.Error(w, "Operation successful but an error occured while returning the user ", http.StatusInternalServerError)
	}
}
