package api

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	userToFollowID, err := strconv.Atoi(r.URL.Query().Get("userToFollowID"))
	if err != nil {
		http.Error(w, "Could not convert the userToFollowID", http.StatusBadRequest)
		return
	}
	userID, er := strconv.Atoi(r.Header.Get("Authorization"))
	if er != nil {
		http.Error(w, "Couldn't identify userId for authentication", http.StatusUnauthorized)
		return
	}
	if userID == userToFollowID {
		http.Error(w, "The follower and followed can't have the same id", http.StatusForbidden)
		return
	}
	var userIds []int
	userIds = append(userIds, userID)
	userIds = append(userIds, userToFollowID)
	i, e := rt.db.VerifyUserIds(userIds)
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		if i == 0 {
			http.Error(w, "The userID can't be found", http.StatusNotFound)
		} else if i == 1 {
			http.Error(w, "The userToFollowID can't be found", http.StatusNotFound)
		}
		return
	}

	s, err := rt.db.CreateFollow(userID, userToFollowID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if err.Error() == "AlreadyFollowed" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	e = json.NewEncoder(w).Encode(s)
	if e != nil {
		http.Error(w, s, http.StatusInternalServerError)
	}
}
