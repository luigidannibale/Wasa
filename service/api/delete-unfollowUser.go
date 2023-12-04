package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	userIDauth, e := strconv.Atoi(r.Header.Get("Authorization"))
	if e != nil {
		http.Error(w, "Couldn't identify userId for authentication "+e.Error(), http.StatusUnauthorized)
		return
	}
	e = rt.db.VerifyUserId(userIDauth)
	if e != nil {
		http.Error(w, "The userID provided for authentication can't be found", http.StatusUnauthorized)
		return
	}
	userIDparam, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Could not convert the userID "+err.Error(), http.StatusBadRequest)
		return
	}
	if userIDauth != userIDparam {
		http.Error(w, "Authentication userID and parameter userID don't match", http.StatusForbidden)
		return
	}
	userID := userIDauth

	userToUnfollowID, err := strconv.Atoi(ps.ByName("followedID"))
	if err != nil {
		http.Error(w, "Could not convert the userToUnfollowID", http.StatusBadRequest)
		return
	}

	if userID == userToUnfollowID {
		http.Error(w, "The follower and followed can't have the same id", http.StatusForbidden)
		return
	}
	s, err := rt.db.DeleteFollow(userID, userToUnfollowID)

	if err != nil {
		if e := err.Error(); e == "UserNotFound" || e == "FollowedNotFound" {
			http.Error(w, s, http.StatusNotFound)
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
