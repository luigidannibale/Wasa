package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var errorMessage string = ""
	//userToUnfollowID, err := strconv.Atoi(r.URL.Query().Get("userToUnfollowID"))
	userToUnfollowID, err := strconv.Atoi(ps.ByName("followedID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage = "Could not convert the userToUnfollowID"
	}
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage = "Could not convert the userID"
	}
	if userID == userToUnfollowID {
		w.WriteHeader(http.StatusForbidden)
		errorMessage = "The follower and followed can't have the same id"
	}
	if errorMessage != "" {
		e := json.NewEncoder(w).Encode(errorMessage)
		if e != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
		}
		return
	}

	s, err := rt.db.DeleteFollow(userID, userToUnfollowID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if e := err.Error(); e == "UserNotFound" || e == "FollowedNotFound" {
			w.WriteHeader(http.StatusNotFound)
			e := json.NewEncoder(w).Encode(s)
			if e != nil {
				http.Error(w, s, http.StatusInternalServerError)
			}
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	e := json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, s, http.StatusInternalServerError)
	}
}
