package api

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/luigidannibale/Wasa/service/utils"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//This should get the body
	var user utils.User
	json.NewDecoder(r.Body).Decode(&user)

	//This takes the userID from parameters
	userID, err := strconv.Atoi(ps.ByName("userID"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message : could not convert the userID")
		return
	}
	userToFollowID, err := strconv.Atoi(ps.ByName("userToFollowID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("message : could not convert the userToFollowID")
		return
	}

	s, err := rt.db.Follow(userID, userToFollowID)

	//Checks for DB-side errrors(404,500)
	if err != nil {
		if e := err.Error(); e == "UserNotFound" || e == "FollowedNotFound" {
			w.WriteHeader(http.StatusNotFound)
			if e == "UserNotFound" {
				json.NewEncoder(w).Encode(userID)
			} else if e == "FollowedNotFound" {
				json.NewEncoder(w).Encode(userToFollowID)
			}
			return
		} else if e == "AlreadyFollowed" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode("message : " + s)
}
